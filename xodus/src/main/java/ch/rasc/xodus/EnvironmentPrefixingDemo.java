package ch.rasc.xodus;

import java.io.File;

import jetbrains.exodus.ByteIterable;
import jetbrains.exodus.bindings.IntegerBinding;
import jetbrains.exodus.bindings.StringBinding;
import jetbrains.exodus.env.Cursor;
import jetbrains.exodus.env.Environment;
import jetbrains.exodus.env.Environments;
import jetbrains.exodus.env.Store;
import jetbrains.exodus.env.StoreConfig;
import jetbrains.exodus.env.Transaction;

public class EnvironmentPrefixingDemo {

  public static void main(String[] args) {

    String sensor1Prefix = "sensor1:";
    String sensor2Prefix = "sensor2:";

    File dbDir = new File("./demo3");
    try (Environment env = Environments.newInstance(dbDir)) {
      Transaction rwTxn = env.beginTransaction();
      try {
        Store sensorStore = env.openStore("Sensor",
            StoreConfig.WITHOUT_DUPLICATES_WITH_PREFIXING, rwTxn);

        for (int i = 0; i < 10; i++) {
          sensorStore.put(rwTxn, StringBinding.stringToEntry(sensor1Prefix + i),
              IntegerBinding.intToEntry((int) (Math.random() * 100)));
          sensorStore.put(rwTxn, StringBinding.stringToEntry(sensor2Prefix + i),
              IntegerBinding.intToEntry((int) (Math.random() * 100)));
        }
      }
      finally {
        rwTxn.commit();
      }

      env.executeInReadonlyTransaction(txn -> {
        Store sensorStore = env.openStore("Sensor", StoreConfig.USE_EXISTING, txn);

        try (Cursor cursor = sensorStore.openCursor(txn)) {
          ByteIterable value = cursor
              .getSearchKeyRange(StringBinding.stringToEntry(sensor1Prefix));
          if (value != null) {
            System.out.println(StringBinding.entryToString(cursor.getKey()));
            System.out.println(IntegerBinding.entryToInt(value));
            while (cursor.getNext()) {
              String key = StringBinding.entryToString(cursor.getKey());
              if (!key.startsWith(sensor1Prefix)) {
                break; // stop if we reach a different prefix
              }
              value = cursor.getValue();
              System.out.println(key);
              System.out.println(IntegerBinding.entryToInt(value));
            }
          }
        }

      });

    }
    catch (Exception e) {
      e.printStackTrace();
    }

  }

}