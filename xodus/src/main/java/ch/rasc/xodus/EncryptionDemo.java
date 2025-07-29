package ch.rasc.xodus;

import jetbrains.exodus.ByteIterable;
import jetbrains.exodus.bindings.IntegerBinding;
import jetbrains.exodus.bindings.StringBinding;
import jetbrains.exodus.entitystore.PersistentEntityStore;
import jetbrains.exodus.entitystore.PersistentEntityStores;
import jetbrains.exodus.env.Environment;
import jetbrains.exodus.env.EnvironmentConfig;
import jetbrains.exodus.env.Environments;
import jetbrains.exodus.env.Store;
import jetbrains.exodus.env.StoreConfig;
import jetbrains.exodus.env.Transaction;

public class EncryptionDemo {
  public static void main(String[] args) {
    EnvironmentConfig config = new EnvironmentConfig();
    config
        .setCipherId("jetbrains.exodus.crypto.streamciphers.ChaChaStreamCipherProvider");

    // for demo purposes hard coded, in a real application read from a secure location
    String cipherKey = "000102030405060708090a0b0c0d0e0f000102030405060708090a0b0c0d0e0f";
    long iv = 314159262718281828L;

    config.setCipherKey(cipherKey);
    config.setCipherBasicIV(iv);

    try (Environment env = Environments.newInstance("./encrypted_db", config)) {
      Transaction rwTx = env.beginTransaction();
      Store sensorStore = env.openStore("sensor_data", StoreConfig.WITHOUT_DUPLICATES,
          rwTx);

      ByteIterable sensor = StringBinding.stringToEntry("sensor1");
      for (int i = 0; i < 10; i++) {
        sensorStore.put(rwTx, sensor, IntegerBinding.intToCompressedEntry(i));
      }

      rwTx.commit();

    }

    try (Environment env = Environments.newInstance("./encrypted_db", config);
        PersistentEntityStore store = PersistentEntityStores.newInstance(env)) {
      // Use the store as usual
    }
  }
}
