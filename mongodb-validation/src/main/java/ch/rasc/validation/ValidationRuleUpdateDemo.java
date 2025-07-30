package ch.rasc.validation;

import java.util.List;

import org.bson.Document;
import org.bson.conversions.Bson;

import com.mongodb.ConnectionString;
import com.mongodb.MongoWriteException;
import com.mongodb.client.ListCollectionsIterable;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoCollection;
import com.mongodb.client.MongoDatabase;
import com.mongodb.client.model.Aggregates;
import com.mongodb.client.model.Filters;
import com.mongodb.client.model.Updates;
import com.mongodb.client.result.InsertOneResult;

public class ValidationRuleUpdateDemo {
  public static void main(String[] args) {
    // create collection without validation rules
    // invalid document
    // add validation rules
    // update existing document and show the difference between ValidationLevel.MODERATE
    // and ValidationLevel.STRICT

    ConnectionString connectionString = new ConnectionString(
        "mongodb://admin:password@localhost:27017");
    try (MongoClient mongoClient = MongoClients.create(connectionString)) {
      run(mongoClient);
    }
  }

  private static void run(MongoClient mongo) {
    MongoDatabase db = mongo.getDatabase("validation");
    db.drop();

    MongoCollection<Document> contacts = db.getCollection("contacts");
    Document contact = new Document("name", "Alice");
    InsertOneResult result = contacts.insertOne(contact);
    System.out.println("Inserted contact with ID: " + result.getInsertedId());

    printCollectionInfo(db, "contacts");

    Document jsonSchema = new Document("$jsonSchema",
        new Document("bsonType", "object").append("required", List.of("name", "phone"))
            .append("properties",
                new Document("name",
                    new Document("bsonType", "string").append("description",
                        "must be a string")).append("phone",
                            new Document("bsonType", "string").append("description",
                                "must be a string"))));

    Document collModCommand = new Document("collMod", "contacts")
        .append("validator", jsonSchema).append("validationLevel", "strict") // "strict"
                                                                             // or
                                                                             // "moderate"
        .append("validationAction", "error"); // "error" or "warn"

    Document collModCommandResult = db.runCommand(collModCommand);
    System.out.println("collMod command executed successfully. Result: "
        + collModCommandResult.toJson());

    printCollectionInfo(db, "contacts");

    // insert valid document
    Document validContact = new Document("name", "Bob").append("phone", "123456789");
    InsertOneResult insertResult = contacts.insertOne(validContact);
    System.out.println("Inserted valid contact with ID: " + insertResult.getInsertedId());

    // update existing document
    Bson updateFilter = Filters.eq("name", "Alice");
    Bson updateOperation = Updates.set("name", "Alice Updated");
    try {
      contacts.updateOne(updateFilter, updateOperation);
    }
    catch (MongoWriteException e) {
      System.out.println("Update failed: " + e);
    }

    // change validation level to moderate
    Document collModCommandModerate = new Document("collMod", "contacts")
        .append("validationLevel", "moderate");
    Document collModCommandModerateResult = db.runCommand(collModCommandModerate);
    System.out.println("Changed validation level to moderate. Result: "
        + collModCommandModerateResult.toJson());

    printCollectionInfo(db, "contacts");

    // update existing document with validation level moderate
    Bson updateFilterModerate = Filters.eq("name", "Alice");
    Bson updateOperationModerate = Updates.set("name", "Alice Updated Moderate");
    contacts.updateOne(updateFilterModerate, updateOperationModerate);

    // list all valid documents in the collection
    for (Document document : contacts.find(jsonSchema)) {
      System.out.println("Valid document: " + document.toJson());
    }

    // list all valid documents in the collection with aggregation
    for (Document document : contacts.aggregate(List.of(Aggregates.match(jsonSchema)))) {
      System.out.println("Aggregated valid document: " + document.toJson());
    }

    // list all invalid documents in the collection
    for (Document document : contacts.find(Filters.nor(jsonSchema))) {
      System.out.println("Invalid document: " + document.toJson());
    }

    // update all invalid documents to make them valid
    Bson updateInvalidFilter = Filters.nor(jsonSchema);
    Bson updateInvalidOperation = Updates.set("phone", "000000000");
    contacts.updateMany(updateInvalidFilter, updateInvalidOperation);

    // list all valid documents after update
    System.out.println("After updating invalid documents:");
    for (Document document : contacts.find(jsonSchema)) {
      System.out.println("Valid document: " + document.toJson());
    }

    // delete all invalid documents
    contacts.deleteMany(Filters.nor(jsonSchema));

  }

  private static void printCollectionInfo(MongoDatabase db, String collectionName) {
    Document collectionInfo = db.listCollections().first();

    System.out.println("Collection: " + collectionName);
    Document options = collectionInfo.get("options", Document.class);
    System.out.println("Validation Rules: " + options.get("validator"));
    System.out.println("Validation Level: " + options.get("validationLevel"));
    System.out.println("Validation Action: " + options.get("validationAction"));
  }
}
