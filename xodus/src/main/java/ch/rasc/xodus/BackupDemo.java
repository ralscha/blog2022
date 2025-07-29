package ch.rasc.xodus;

import java.io.File;

import jetbrains.exodus.ByteIterable;
import jetbrains.exodus.bindings.IntegerBinding;
import jetbrains.exodus.bindings.StringBinding;
import jetbrains.exodus.env.Environment;
import jetbrains.exodus.env.Environments;
import jetbrains.exodus.env.Store;
import jetbrains.exodus.env.StoreConfig;
import jetbrains.exodus.env.Transaction;
import jetbrains.exodus.util.CompressBackupUtil;

public class BackupDemo {
  public static void main(String[] args) throws Exception {
    try (Environment env = Environments.newInstance("./backup_db")) {
      Transaction rwTx = env.beginTransaction();
      Store sensorStore = env.openStore("sensor_data", StoreConfig.WITHOUT_DUPLICATES,
          rwTx);

      ByteIterable sensor = StringBinding.stringToEntry("sensor1");
      for (int i = 0; i < 10; i++) {
        sensorStore.put(rwTx, sensor, IntegerBinding.intToCompressedEntry(i));
      }

      rwTx.commit();

      File backupFile = CompressBackupUtil.backup(env,
          new File(env.getLocation(), "sensor_data_backup"), "sensor_data_", true);
      System.out.println("Backup file: " + backupFile.getAbsolutePath());
    }
  }
}
