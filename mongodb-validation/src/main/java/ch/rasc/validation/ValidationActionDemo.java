package ch.rasc.validation;

import org.bson.Document;

import com.mongodb.ConnectionString;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import com.mongodb.client.model.CreateCollectionOptions;
import com.mongodb.client.model.ValidationAction;
import com.mongodb.client.model.ValidationLevel;
import com.mongodb.client.model.ValidationOptions;
import com.mongodb.client.result.InsertOneResult;

public class ValidationActionDemo {
  public static void main(String[] args) {

    ConnectionString connectionString = new ConnectionString(
        "mongodb://admin:password@localhost:27017");

    try (MongoClient mongoClient = MongoClients.create(connectionString)) {
      MongoDatabase db = mongoClient.getDatabase("validation");
      Document jsonSchema = Document.parse("""
          {
            $jsonSchema: {
              bsonType: "object",
              required: ["name", "phone"],
              properties: {
                 name: {
                   bsonType: "string",
                   description: "must be a string"
                 },
                 phone: {
                   bsonType: "string",
                   description: "must be a string"
                 }
              }
            }
          }
          """);

      ValidationOptions validationOptions = new ValidationOptions().validator(jsonSchema)
          .validationLevel(ValidationLevel.MODERATE)
          .validationAction(ValidationAction.WARN);

      db.createCollection("contacts",
          new CreateCollectionOptions().validationOptions(validationOptions));

      MongoCollection<Document> contacts = db.getCollection("contacts");
      Document contact = new Document("name", "Alice");
      InsertOneResult result = contacts.insertOne(contact);
      System.out.println("Inserted contact with ID: " + result.getInsertedId());
    }

    try (MongoClient mongoClient = MongoClients.create(connectionString)) {
      MongoDatabase db = mongoClient.getDatabase("validation2");
      Document jsonSchema = Document.parse("""
          {
            $jsonSchema: {
              bsonType: "object",
              required: ["name", "phone"],
              properties: {
                 name: {
                   bsonType: "string",
                   description: "must be a string"
                 },
                 phone: {
                   bsonType: "string",
                   description: "must be a string"
                 }
              }
            }
          }
          """);

      ValidationOptions validationOptions = new ValidationOptions().validator(jsonSchema)
          .validationLevel(ValidationLevel.MODERATE)
          .validationAction(ValidationAction.ERROR);

      db.createCollection("contacts",
          new CreateCollectionOptions().validationOptions(validationOptions));

      MongoCollection<Document> contacts = db.getCollection("contacts");
      Document contact = new Document("name", "Alice");
      InsertOneResult result = contacts.insertOne(contact);
      System.out.println("Inserted contact with ID: " + result.getInsertedId());
    }

  }
}
