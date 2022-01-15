package ch.rasc.consumer;

import com.google.protobuf.InvalidProtocolBufferException;

import ch.rasc.producer.Adress.AddressBook;
import ch.rasc.producer.Adress.Person;
import ch.rasc.producer.Adress.Person.PhoneNumber;
import software.amazon.awssdk.auth.credentials.ProfileCredentialsProvider;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.DeleteMessageRequest;
import software.amazon.awssdk.services.sqs.model.Message;
import software.amazon.awssdk.services.sqs.model.ReceiveMessageRequest;
import software.amazon.awssdk.services.sqs.model.ReceiveMessageResponse;

public class Consumer {

	private final static String queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4";

	public static void main(String[] args) throws InvalidProtocolBufferException {

		try (ProfileCredentialsProvider awsCredentials = ProfileCredentialsProvider
				.create("home");
				SqsClient sqsClient = SqsClient.builder()
						.credentialsProvider(awsCredentials).region(Region.US_EAST_1)
						.build()) {

			while (true) {
				System.out.println("polling for message");
				ReceiveMessageRequest request = ReceiveMessageRequest.builder()
						.queueUrl(queueURL).maxNumberOfMessages(10).waitTimeSeconds(20)
						.messageAttributeNames("Body").build();
				ReceiveMessageResponse response = sqsClient.receiveMessage(request);

				for (Message message : response.messages()) {
					AddressBook ab = AddressBook.parseFrom(message.messageAttributes()
							.get("Body").binaryValue().asByteArray());
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
					sqsClient.deleteMessage(
							DeleteMessageRequest.builder().queueUrl(queueURL)
									.receiptHandle(message.receiptHandle()).build());
				}
			}
		}

	}

}
