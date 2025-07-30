package ch.rasc.validation;

import org.bson.Document;

import com.mongodb.ConnectionString;
import com.mongodb.MongoException;
import com.mongodb.client.MongoClient;
import com.mongodb.client.MongoClients;
import com.mongodb.client.MongoDatabase;
import com.mongodb.client.model.CreateCollectionOptions;
import com.mongodb.client.model.ValidationOptions;

public class ValidationJsonSchema {
  public static void main(String[] args) throws MongoException {
    ConnectionString connectionString = new ConnectionString(
        "mongodb://admin:password@localhost:27017");
    try (MongoClient mongoClient = MongoClients.create(connectionString)) {
      run(mongoClient);
    }
  }

  private static void run(MongoClient mongo) {
    MongoDatabase db = mongo.getDatabase("validation");
    db.drop();
    var jsonSchema = Document.parse(
        """
        {
          $jsonSchema: {
            type: "object",
            required: [ "username", "status", "address" ],
            additionalProperties: false,
            properties: {
               _id: { bsonType: "objectId" },
                  username: { type: "string", minLength: 1, description: "username" },
                  firstName: { type: ["string", "null"] },
                  lastName: { type: "string" },
                  email: { type: "string", pattern: "^.+@.+$"},
                  birthYear: { bsonType: "int", minimum: 1900, maximum: 2025 },
                  hobbies: {
                    type: "array",
                    items: [ { "type": "string", "enum": ["Reading", "Swimming", "Cycling", "Hiking", "Painting"] } ],
               minItems: 1,
               uniqueItems: true
                  },
                  status: { type: "string", enum: [ "active", "inactive" ] },
                  address: {
                    type: "object",
                    required: [ "city" ],
                    additionalProperties: false,
                    properties: {
                      city: { type: "string", minLength: 1 },
                      street: { type: "string" },
                      postalCode: { bsonType: "int" }
                    }
                  }
            }
          }
  }""");

    var validationOptions = new ValidationOptions();
    validationOptions.validator(jsonSchema);

    db.createCollection("users",
        new CreateCollectionOptions().validationOptions(validationOptions));

    var usersCollection = db.getCollection("users");

    var json = """
        {
          "username": "admin",
          "firstName": "John",
          "lastName": "Doe",
          "email": "test@test.com",
          "birthYear": 1988,
          "hobbies": ["Reading", "Swimming"],
          "status": "active",
          "address": {
             "city": "BigCity",
             "street": "MainRoad 10",
             "postalCode": 10000
          }
        }
        """;
    usersCollection.insertOne(Document.parse(json));

    json = """
        {
          username: "admin",
          firstName: null,
          status: "active",
          address: {
             city: "test",
          }
        }
        """;
    usersCollection.insertOne(Document.parse(json));

  }

}
