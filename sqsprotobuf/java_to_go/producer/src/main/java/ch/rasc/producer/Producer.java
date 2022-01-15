package ch.rasc.producer;

import java.util.Map;

import ch.rasc.producer.Adress.AddressBook;
import ch.rasc.producer.Adress.Person;
import ch.rasc.producer.Adress.Person.PhoneNumber;
import ch.rasc.producer.Adress.Person.PhoneType;
import software.amazon.awssdk.auth.credentials.ProfileCredentialsProvider;
import software.amazon.awssdk.core.SdkBytes;
import software.amazon.awssdk.regions.Region;
import software.amazon.awssdk.services.sqs.SqsClient;
import software.amazon.awssdk.services.sqs.model.MessageAttributeValue;
import software.amazon.awssdk.services.sqs.model.SendMessageRequest;
import software.amazon.awssdk.services.sqs.model.SendMessageResponse;

public class Producer {

	private final static String queueURL = "https://sqs.us-east-1.amazonaws.com/660461151343/queue-d8494a4";

	public static void main(String[] args) {

		PhoneNumber pn = PhoneNumber.newBuilder().setType(PhoneType.MOBILE)
				.setNumber("1111").build();
		Person person = Person.newBuilder().setId(1).setEmail("john@test.com")
				.setName("John").addPhones(pn).build();
		AddressBook book = AddressBook.newBuilder().addPeople(person).build();

		try (ProfileCredentialsProvider awsCredentials = ProfileCredentialsProvider
				.create("home");
				SqsClient client = SqsClient.builder().credentialsProvider(awsCredentials)
						.region(Region.US_EAST_1).build()) {

			SendMessageRequest request = SendMessageRequest.builder().queueUrl(queueURL)
					.messageBody(" ")
					.messageAttributes(Map.of("Body",
							MessageAttributeValue.builder().dataType("Binary")
									.binaryValue(
											SdkBytes.fromByteArray(book.toByteArray()))
									.build()))
					.build();

			SendMessageResponse response = client.sendMessage(request);
			System.out.println("Message sent: " + response.messageId());
		}

	}

}
