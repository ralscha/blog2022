package ch.rasc.xodus;

import java.io.File;

import jetbrains.exodus.entitystore.Entity;
import jetbrains.exodus.entitystore.EntityId;
import jetbrains.exodus.entitystore.EntityIterable;
import jetbrains.exodus.entitystore.PersistentEntityStore;
import jetbrains.exodus.entitystore.PersistentEntityStores;

public class EntityStoreDemo {

  public static void main(String[] args) {

    File dbDir = new File("my_entity_store_db");
    try (PersistentEntityStore entityStore = PersistentEntityStores.newInstance(dbDir)) {

      EntityId user1EntityId = entityStore.computeInTransaction(txn -> {
        Entity user1 = txn.newEntity("User");
        user1.setProperty("userId", 1);
        user1.setProperty("username", "john_doe");
        user1.setProperty("age", 25);
        user1.setProperty("email", "john.doe@example.com");
        user1.setProperty("active", true);
        user1.setProperty("created", System.currentTimeMillis());
        // user1.setBlob("image", new File(...));

        Entity user2 = txn.newEntity("User");
        user2.setProperty("userId", 2);
        user2.setProperty("username", "jane_doe");
        user2.setProperty("age", 42);
        user2.setProperty("email", "jane.doe@example.com");
        user2.setProperty("active", false);
        user2.setProperty("created", System.currentTimeMillis());

        Entity post1 = txn.newEntity("Post");
        post1.setProperty("postId", txn.getSequence("postsSequence").increment());
        post1.setProperty("title", "Hello World!");
        post1.setLink("author", user1);

        user1.addLink("posts", post1);

        Entity comment1 = txn.newEntity("Comment");
        comment1.setProperty("commentId",
            txn.getSequence("commentsSequence").increment());
        comment1.setProperty("text", "Great post!");
        comment1.setLink("post", post1);
        comment1.setLink("author", user2);

        user2.addLink("comments", comment1);

        return user1.getId();
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        EntityIterable allUsers = txn.getAll("User");

        // Fast check if cached
        long count = allUsers.count();
        if (count >= 0) {
          System.out.println("User count (cached): " + count);
        }
        else {
          // Slower
          long size = allUsers.size();
          System.out.println("User count: " + size);
        }

        // Check if empty
        if (allUsers.isEmpty()) {
          System.out.println("No users found");
        }
      });

      entityStore.executeInReadonlyTransaction(_ -> {
        Entity u1 = entityStore.getEntity(user1EntityId);
        System.out.println(u1.getProperty("username"));
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        System.out.println("All users:");
        EntityIterable allUsers = txn.getAll("User");
        for (Entity entity : allUsers) {
          System.out.println("User: " + entity.getProperty("username"));
        }

        System.out.println("Sort by age DESC");
        EntityIterable allUsersSorted = txn.sort("User", "age", false);
        for (Entity entity : allUsersSorted) {
          System.out.println("User: " + entity.getProperty("username"));
        }

        System.out.println("Sort by age ASC and then by username DESC");
        allUsersSorted = txn.sort("User", "age", txn.sort("User", "username", false),
            true);
        for (Entity entity : allUsersSorted) {
          System.out.println("User: " + entity.getProperty("username"));
        }
      });

      String username = entityStore.computeInReadonlyTransaction(txn -> {
        EntityIterable users = txn.find("User", "username", "john_doe");
        if (!users.isEmpty()) {
          Entity u = users.getFirst();
          return u.getProperty("username").toString();
        }
        return null;
      });
      System.out.println("Username: " + username);

      entityStore.executeInReadonlyTransaction(txn -> {
        EntityIterable users = txn.find("User", "age", 18, 25);
        System.out.println("Users between 18 and 25:");
        for (Entity entity : users) {
          System.out.println("User: " + entity.getProperty("username"));
        }
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        // intersect: AND
        EntityIterable users = txn.findStartingWith("User", "username", "john")
            .intersect(txn.findStartingWith("User", "email", "john.doe@"));
        System.out.println(
            "Users with username starting with 'john' AND email starting with 'john.doe@':");
        for (Entity entity : users) {
          System.out.println("User: " + entity.getProperty("username"));
        }

        // union: OR
        EntityIterable users2 = txn.findStartingWith("User", "username", "john")
            .union(txn.findStartingWith("User", "email", "john.doe@").distinct());
        System.out.println(
            "Users with username starting with 'john' OR email starting with 'john.doe@':");
        for (Entity entity : users2) {
          System.out.println("User: " + entity.getProperty("username"));
        }

        // minus: AND NOT
        EntityIterable users3 = txn.findStartingWith("User", "username", "john")
            .minus(txn.findStartingWith("User", "email", "john.doe@"));
        System.out.println(
            "Users with username starting with 'john' AND NOT email starting with 'john.doe@':");
        for (Entity entity : users3) {
          System.out.println("User: " + entity.getProperty("username"));
        }

      });

      entityStore.executeInReadonlyTransaction(txn -> {
        // Find users between age 20 and 30
        EntityIterable youngUsers = txn.find("User", "age", 20, 30);
        System.out.println("Users between age 20 and 30:");
        for (Entity user : youngUsers) {
          System.out.println("User: " + user.getProperty("username"));
        }

        // Find users whose username starts with "john"
        EntityIterable johnUsers = txn.findStartingWith("User", "username", "JOHN");
        System.out.println("Users whose username starts with 'john':");
        for (Entity user : johnUsers) {
          System.out.println("User: " + user.getProperty("username"));
        }

        for (Entity user : johnUsers) {
          System.out.println("User: " + user.getProperty("username"));
        }

        EntityIterable usersContainingDoe = txn.findContaining("User", "username", "DOE",
            true);
        System.out.println("Users whose username contains 'doe':");
        for (Entity user : usersContainingDoe) {
          System.out.println("User: " + user.getProperty("username"));
        }
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        Entity johnUser = txn.find("User", "username", "JOHN_DOE").getFirst();
        if (johnUser != null) {
          System.out.println("Found user: " + johnUser.getProperty("email"));
        }
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        EntityIterable u = txn.find("User", "userId", 1);
        if (!u.isEmpty()) {
          Entity user1 = u.getFirst();
          EntityIterable postsByUser1 = txn.findLinks("Post", user1, "author");
          for (Entity post : postsByUser1) {
            System.out.println("Post by user1: " + post.getProperty("title"));
          }
        }

        u = txn.find("User", "userId", 2);
        if (!u.isEmpty()) {
          Entity user2 = u.getFirst();
          EntityIterable commentsByUser2 = user2.getLinks("comments");
          for (Entity comment : commentsByUser2) {
            System.out.println("Comment by user2: " + comment.getProperty("text"));
          }
        }

        // get all posts from users under 30
        EntityIterable users = txn.find("User", "age", 0, 30);
        EntityIterable posts = users.selectManyDistinct("posts");
        for (Entity post : posts) {
          System.out.println("Post: " + post.getProperty("title"));
        }

        // Get all authors of posts (single link)
        EntityIterable postAuthors = txn.getAll("Post").selectDistinct("author");
        for (Entity author : postAuthors) {
          System.out.println("Author: " + author.getProperty("username"));
        }
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        EntityIterable users = txn.getAll("User");
        long exactCount = users.size();
        long roughCount = users.getRoughCount();
        boolean isEmpty = users.isEmpty();
        System.out.println("Exact count: " + exactCount);
        System.out.println("Rough count: " + roughCount);
        System.out.println("Is empty: " + isEmpty);
      });

      entityStore.executeInReadonlyTransaction(txn -> {
        Entity johnUser = txn.find("User", "username", "john_doe").getFirst();
        if (johnUser != null) {
          EntityIterable userPosts = johnUser.getLinks("posts");
          for (Entity post : userPosts) {
            System.out.println("Post: " + post.getProperty("title"));
          }
        }

        Entity janeUser = txn.find("User", "username", "jane_doe").getFirst();
        if (janeUser != null) {
          EntityIterable postComments = txn.findLinks("Comment", janeUser, "author");
          for (Entity comment : postComments) {
            System.out.println("Comment: " + comment.getProperty("text"));
          }
        }
      });

      // findWith
      entityStore.executeInReadonlyTransaction(txn -> {
        EntityIterable usersWithAge = txn.findWithProp("User", "age");
        System.out.println("Number of users with age property: " + usersWithAge.size());

        EntityIterable usersWithoutAge = txn.getAll("User")
            .minus(txn.findWithProp("User", "age"));
        System.out
            .println("Number of users without age property: " + usersWithoutAge.size());

        EntityIterable usersWithImage = txn.findWithBlob("User", "image");
        System.out.println("Users with image blob: " + usersWithImage.size());

        EntityIterable usersWithPosts = txn.findWithLinks("User", "posts");
        System.out.println("Users with posts links: " + usersWithPosts.size());

      });
    }

  }

}