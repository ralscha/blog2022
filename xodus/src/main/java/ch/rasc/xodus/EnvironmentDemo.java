package ch.rasc.xodus;

import java.io.File;
import java.util.List;

import jetbrains.exodus.ByteIterable;
import jetbrains.exodus.bindings.LongBinding;
import jetbrains.exodus.bindings.StringBinding;
import jetbrains.exodus.env.ContextualEnvironment;
import jetbrains.exodus.env.ContextualStore;
import jetbrains.exodus.env.Cursor;
import jetbrains.exodus.env.Environment;
import jetbrains.exodus.env.EnvironmentConfig;
import jetbrains.exodus.env.Environments;
import jetbrains.exodus.env.Store;
import jetbrains.exodus.env.StoreConfig;
import jetbrains.exodus.env.Transaction;

public class EnvironmentDemo {

  public static void main(String[] args) {

    try (Environment env = Environments.newInstance("./data")) {
      // insert
      env.executeInTransaction(txn -> {
        Store store = env.openStore("Users", StoreConfig.WITHOUT_DUPLICATES, txn);
        store.put(txn, StringBinding.stringToEntry("user1"),
            StringBinding.stringToEntry("Alice"));
      });

      // update
      env.executeInTransaction(txn -> {
        Store store = env.openStore("Users", StoreConfig.WITHOUT_DUPLICATES, txn);
        store.put(txn, StringBinding.stringToEntry("user1"),
            StringBinding.stringToEntry("Alice Smith"));
      });

      // read
      env.executeInReadonlyTransaction(txn -> {
        Store store = env.openStore("Users", StoreConfig.WITHOUT_DUPLICATES, txn);
        ByteIterable value = store.get(txn, StringBinding.stringToEntry("user1"));
        System.out.println("User1: " + StringBinding.entryToString(value));
      });

      // delete
      env.executeInTransaction(txn -> {
        Store store = env.openStore("Users", StoreConfig.WITHOUT_DUPLICATES, txn);
        store.delete(txn, StringBinding.stringToEntry("user1"));
      });

      // exclusive transaction
      env.executeInExclusiveTransaction(txn -> {
        Store store = env.openStore("Users", StoreConfig.WITHOUT_DUPLICATES, txn);
        store.put(txn, StringBinding.stringToEntry("user2"),
            StringBinding.stringToEntry("Bob"));
      });
    }

    try (Environment env = Environments.newInstance("./data")) {
      Transaction txn = env.beginReadonlyTransaction();
      try {
        // Check if store exists
        boolean storeExists = env.storeExists("Users", txn);
        System.out.println("Store 'Users' exists: " + storeExists);

        // List all stores
        List<String> storeNames = env.getAllStoreNames(txn);
        System.out.println("Stores in the environment:");
        for (String name : storeNames) {
          System.out.println("- " + name);
        }
      }
      finally {
        txn.abort(); // Readonly transaction does not need to be committed
      }
    }

    try (ContextualEnvironment env = Environments.newContextualInstance("./data")) {
      try {
        env.beginTransaction();
        ContextualStore store = env.openStore("Users", StoreConfig.WITHOUT_DUPLICATES);
        store.put(StringBinding.stringToEntry("user1"),
            StringBinding.stringToEntry("Alice"));
      }
      finally {
        env.getCurrentTransaction().commit();
      }
    }

    File dbDir = new File("./demo");
    EnvironmentConfig config = new EnvironmentConfig();

    try (Environment env = Environments.newInstance(dbDir, config)) {
      Transaction rwTx = env.beginTransaction();
      Store myStore = env.openStore("MyStore", StoreConfig.WITHOUT_DUPLICATES, rwTx);

      myStore.put(rwTx, LongBinding.longToEntry(1),
          StringBinding.stringToEntry("value1"));
      myStore.put(rwTx, LongBinding.longToEntry(2),
          StringBinding.stringToEntry("value2"));
      myStore.putRight(rwTx, LongBinding.longToEntry(3),
          StringBinding.stringToEntry("value3"));

      rwTx.commit();

      env.executeInReadonlyTransaction(txn -> {
        ByteIterable value1 = myStore.get(txn, LongBinding.longToEntry(1));
        System.out
            .println("Value for key 'key1': " + StringBinding.entryToString(value1));

        ByteIterable value4 = myStore.get(txn, LongBinding.longToEntry(4));
        System.out.println(value4 == null ? "Key 'key4' not found"
            : "Value for key 'key4': " + StringBinding.entryToString(value4));

        long count = myStore.count(txn);
        System.out.println("Store count: " + count);

        boolean exists = myStore.exists(txn, LongBinding.longToEntry(2),
            StringBinding.stringToEntry("value2"));
        System.out.println("Key 'key2' with value 'value2' exists: " + exists);

        try (Cursor cursor = myStore.openCursor(txn)) {
          while (cursor.getNext()) {
            Long key = LongBinding.entryToLong(cursor.getKey());
            String value = StringBinding.entryToString(cursor.getValue());
            System.out.println("Cursor: Key=" + key + ", Value=" + value);
          }
        }
      });

      env.executeInTransaction(txn -> {
        myStore.delete(txn, StringBinding.stringToEntry("key1"));
      });

    }

  }

}