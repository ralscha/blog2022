package ch.rasc.producer;

import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

import com.github.javafaker.Faker;

import ch.rasc.producer.Adress.AddressBook;
import ch.rasc.producer.Adress.Person;
import ch.rasc.producer.Adress.Person.PhoneNumber;
import ch.rasc.producer.Adress.Person.PhoneType;
import software.amazon.awssdk.auth.credentials.ProfileCredentialsProvider;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.SendMessageRequest;
import software.amazon.awssdk.services.sqs.model.SendMessageResponse;

public class Producer {

	private final static String queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4";
	private final static String messageBucket = "messages-a5b0326";

	public static void main(String[] args) {

		Faker faker = new Faker();
		List<Person> persons = new ArrayList<>();

		for (int i = 1; i < 10_000; i++) {
			PhoneNumber pn = PhoneNumber.newBuilder().setType(PhoneType.MOBILE)
					.setNumber(faker.phoneNumber().cellPhone()).build();
			Person person = Person.newBuilder().setId(i)
					.setEmail(faker.internet().emailAddress())
					.setName(faker.name().name()).addPhones(pn).build();
			persons.add(person);
		}

		AddressBook book = AddressBook.newBuilder().addAllPeople(persons).build();

		try (ProfileCredentialsProvider awsCredentials = ProfileCredentialsProvider
				.create("home");
				S3Client s3Client = S3Client.builder().credentialsProvider(awsCredentials)
						.region(Region.US_EAST_1).build();
				SqsClient sqsClient = SqsClient.builder()
						.credentialsProvider(awsCredentials).region(Region.US_EAST_1)
						.build()) {

			String s3Key = UUID.randomUUID().toString();
			s3Client.putObject(
					PutObjectRequest.builder().bucket(messageBucket).key(s3Key).build(),
					RequestBody.fromBytes(book.toByteArray()));

			SendMessageRequest request = SendMessageRequest.builder().queueUrl(queueURL)
					.messageBody(s3Key).build();

			SendMessageResponse response = sqsClient.sendMessage(request);
			System.out.println("Message sent: " + response.messageId());
		}

	}

}
