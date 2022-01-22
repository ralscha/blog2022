package ch.rasc.s3selectdemo;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;

import com.amazonaws.auth.AWSStaticCredentialsProvider;
import com.amazonaws.auth.BasicSessionCredentials;
import com.amazonaws.regions.Regions;
import com.amazonaws.services.s3.AmazonS3;
import com.amazonaws.services.s3.AmazonS3Client;
import com.amazonaws.services.s3.model.CompressionType;
import com.amazonaws.services.s3.model.ExpressionType;
import com.amazonaws.services.s3.model.InputSerialization;
import com.amazonaws.services.s3.model.JSONInput;
import com.amazonaws.services.s3.model.JSONOutput;
import com.amazonaws.services.s3.model.JSONType;
import com.amazonaws.services.s3.model.OutputSerialization;
import com.amazonaws.services.s3.model.SelectObjectContentEvent;
import com.amazonaws.services.s3.model.SelectObjectContentEvent.RecordsEvent;
import com.amazonaws.services.s3.model.SelectObjectContentEventStream;
import com.amazonaws.services.s3.model.SelectObjectContentEventVisitor;
import com.amazonaws.services.s3.model.SelectObjectContentRequest;
import com.amazonaws.services.s3.model.SelectObjectContentResult;

public class Main {
	public static void main(String[] args) throws IOException {

		String accessKey = "";
		String secretKey = "";
		String sessionToken = "";

		AmazonS3 s3Client = AmazonS3Client.builder().withRegion(Regions.US_EAST_1)
				.withCredentials(new AWSStaticCredentialsProvider(
						new BasicSessionCredentials(accessKey, secretKey, sessionToken)))
				.build();

		// if the application runs on AWS
		// AmazonS3 s3Client = AmazonS3Client.builder().build();

		String query = "select p.id,p.name from S3Object[*].pokemon[*] p";
		selectObject(s3Client, query, false);

		query = "select p.id,p.name,p.type from S3Object[*].pokemon[*] p where p.type[0] = 'Fire' or p.type[1] = 'Fire'";
		selectObject(s3Client, query, false);

		query = "select p from S3Object[*].pokemon[*] p where p.name = 'Charmander'";
		selectObject(s3Client, query, false);
		
		query = "select count(*) from S3Object[*].pokemon[*] p";
		selectObject(s3Client, query, true);		
	}

	private static void selectObject(AmazonS3 s3Client, String query, boolean compressed)
			throws IOException {
		String bucket = "select-demo";
		String keyUncompressed = "pokedex.json";
		String keyCompressed = "pokedex.json.bz2";
		SelectObjectContentRequest request = new SelectObjectContentRequest();
		request.setBucketName(bucket);
		if (compressed) {
			request.setKey(keyCompressed);
		}
		else {
			request.setKey(keyUncompressed);
		}
		request.setExpression(query);
		request.setExpressionType(ExpressionType.SQL);

		InputSerialization inputSerialization = new InputSerialization();
		/*
		 * LINES means that each line in the input data contains a single JSON object.
		 * DOCUMENT means that a single JSON object can span multiple lines in the input.
		 * Using DOCUMENT might result in slower performance in some cases.
		 */
		inputSerialization.setJson(new JSONInput().withType(JSONType.DOCUMENT));
		if (compressed) {
			inputSerialization.setCompressionType(CompressionType.BZIP2);
		}
		else {
			inputSerialization.setCompressionType(CompressionType.NONE);
		}
		request.setInputSerialization(inputSerialization);

		OutputSerialization outputSerialization = new OutputSerialization();
		outputSerialization.setJson(new JSONOutput());
		request.setOutputSerialization(outputSerialization);

		SelectObjectContentEventVisitor listener = new SelectObjectContentEventVisitor() {
			@Override
			public void visit(RecordsEvent event) {
				System.out.println("record event");
			}

			@Override
			public void visit(SelectObjectContentEvent.StatsEvent event) {
				System.out.println("stats event: ");
				System.out.println("uncompressed bytes processed: "
						+ event.getDetails().getBytesProcessed());
			}

			@Override
			public void visit(SelectObjectContentEvent.EndEvent event) {
				System.out.println("end event");
			}
		};

		try (ByteArrayOutputStream baos = new ByteArrayOutputStream();
				SelectObjectContentResult result = s3Client.selectObjectContent(request);
				SelectObjectContentEventStream payload = result.getPayload();
				InputStream is = payload.getRecordsInputStream(listener)) {
			is.transferTo(baos);
			System.out.println(new String(baos.toByteArray()));
		}

		// without visitor
		try (ByteArrayOutputStream baos = new ByteArrayOutputStream();
				SelectObjectContentResult result = s3Client.selectObjectContent(request);
				SelectObjectContentEventStream payload = result.getPayload();
				InputStream is = payload.getRecordsInputStream()) {
			is.transferTo(baos);
			System.out.println(new String(baos.toByteArray()));
		}
	}
}