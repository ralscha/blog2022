package ch.rasc.consumer;

import java.io.IOException;

import ch.rasc.producer.Adress.AddressBook;
import ch.rasc.producer.Adress.Person;
import ch.rasc.producer.Adress.Person.PhoneNumber;
import software.amazon.awssdk.auth.credentials.ProfileCredentialsProvider;
import software.amazon.awssdk.awscore.exception.AwsServiceException;
import software.amazon.awssdk.core.ResponseInputStream;
import software.amazon.awssdk.core.exception.SdkClientException;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.DeleteObjectRequest;
import software.amazon.awssdk.services.s3.model.GetObjectRequest;
import software.amazon.awssdk.services.s3.model.GetObjectResponse;
import software.amazon.awssdk.services.s3.model.InvalidObjectStateException;
import software.amazon.awssdk.services.s3.model.NoSuchKeyException;
import software.amazon.awssdk.services.s3.model.S3Exception;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.DeleteMessageRequest;
import software.amazon.awssdk.services.sqs.model.Message;
import software.amazon.awssdk.services.sqs.model.ReceiveMessageRequest;
import software.amazon.awssdk.services.sqs.model.ReceiveMessageResponse;

public class Consumer {

	private final static String queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4";
	private final static String messageBucket = "messages-a5b0326";

	public static void main(String[] args)
			throws NoSuchKeyException, InvalidObjectStateException, S3Exception,
			AwsServiceException, SdkClientException, IOException {

		try (ProfileCredentialsProvider awsCredentials = ProfileCredentialsProvider
				.create("home");
				S3Client s3Client = S3Client.builder().credentialsProvider(awsCredentials)
						.region(Region.US_EAST_1).build();
				SqsClient sqsClient = SqsClient.builder()
						.credentialsProvider(awsCredentials).region(Region.US_EAST_1)
						.build()) {

			while (true) {
				System.out.println("polling for message");
				ReceiveMessageRequest request = ReceiveMessageRequest.builder()
						.queueUrl(queueURL).maxNumberOfMessages(10).waitTimeSeconds(20)
						.build();
				ReceiveMessageResponse response = sqsClient.receiveMessage(request);

				for (Message message : response.messages()) {
					String s3Key = message.body();
					try (ResponseInputStream<GetObjectResponse> ris = s3Client
							.getObject(GetObjectRequest.builder().bucket(messageBucket)
									.key(s3Key).build())) {
						AddressBook ab = AddressBook.parseFrom(ris);

						for (Person person : ab.getPeopleList()) {
							System.out.println(person.getId());
							System.out.println(person.getName());
							System.out.println(person.getEmail());
							for (PhoneNumber phone : person.getPhonesList()) {
								System.out.print(phone.getType());
								System.out.print(" : ");
								System.out.println(phone.getNumber());
							}
						}
					}

					s3Client.deleteObject(DeleteObjectRequest.builder()
							.bucket(messageBucket).key(s3Key).build());
					sqsClient.deleteMessage(
							DeleteMessageRequest.builder().queueUrl(queueURL)
									.receiptHandle(message.receiptHandle()).build());
				}

			}
		}

	}

}
