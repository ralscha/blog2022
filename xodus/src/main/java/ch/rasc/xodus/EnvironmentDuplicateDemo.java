package ch.rasc.xodus;

import java.io.File;

import jetbrains.exodus.ByteIterable;
import jetbrains.exodus.bindings.LongBinding;
import jetbrains.exodus.bindings.StringBinding;
import jetbrains.exodus.env.Cursor;
import jetbrains.exodus.env.Environment;
import jetbrains.exodus.env.Environments;
import jetbrains.exodus.env.Store;
import jetbrains.exodus.env.StoreConfig;
import jetbrains.exodus.env.Transaction;

public class EnvironmentDuplicateDemo {

  public static void main(String[] args) {

    File dbDir = new File("./demo2");
    try (Environment env = Environments.newInstance(dbDir)) {
      Transaction rwTxn = env.beginTransaction();
      Store userStore = env.openStore("User", StoreConfig.WITHOUT_DUPLICATES, rwTxn);
      Store permissionsStore = env.openStore("Permissions", StoreConfig.WITH_DUPLICATES,
          rwTxn);

      userStore.put(rwTxn, LongBinding.longToEntry(1),
          StringBinding.stringToEntry("admin"));
      userStore.put(rwTxn, LongBinding.longToEntry(2),
          StringBinding.stringToEntry("user"));

      permissionsStore.put(rwTxn, LongBinding.longToEntry(1),
          StringBinding.stringToEntry("read"));
      permissionsStore.put(rwTxn, LongBinding.longToEntry(1),
          StringBinding.stringToEntry("write"));
      permissionsStore.put(rwTxn, LongBinding.longToEntry(2),
          StringBinding.stringToEntry("read"));
      rwTxn.commit();

      env.executeInReadonlyTransaction(txn -> {

        String adminUser = StringBinding
            .entryToString(userStore.get(txn, LongBinding.longToEntry(1)));
        System.out.println(adminUser); // admin

        String permission = StringBinding
            .entryToString(permissionsStore.get(txn, LongBinding.longToEntry(1)));
        System.out.println(permission); // read

        try (Cursor cursor = permissionsStore.openCursor(txn)) {
          ByteIterable bi = cursor.getSearchKey(LongBinding.longToEntry(1));
          if (bi != null) {
            System.out.println(StringBinding.entryToString(bi));
            while (cursor.getNextDup()) {
              bi = cursor.getValue();
              System.out.println(StringBinding.entryToString(bi)); // first read then
                                                                   // write
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