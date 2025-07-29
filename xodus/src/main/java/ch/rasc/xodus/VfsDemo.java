package ch.rasc.xodus;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;

import jetbrains.exodus.env.Environment;
import jetbrains.exodus.env.Environments;
import jetbrains.exodus.vfs.VirtualFileSystem;

public class VfsDemo {

  public static void main(String[] args) {
    File dbDir = new File("vfs_demo_db");

    try (Environment env = Environments.newInstance(dbDir)) {
      VirtualFileSystem vfs = new VirtualFileSystem(env);

      try {
        demonstrateFileOperations(env, vfs);
        demonstrateFilePositioning(env, vfs);
      }
      finally {
        vfs.shutdown();
      }
    }

  }

  private static void demonstrateFileOperations(Environment env, VirtualFileSystem vfs) {
    env.executeInTransaction(txn -> {
      try {
        // Create a new file. Throws exception if file already exists.
        jetbrains.exodus.vfs.File file1 = vfs.createFile(txn, "demo/test1.txt");
        System.out.println("Created file: " + file1.getPath() + " (descriptor: "
            + file1.getDescriptor() + ")");

        // Write content to the file
        try (OutputStream output = vfs.writeFile(txn, file1)) {
          String content = "This is a test file.\nLine 2 of content.\n";
          output.write(content.getBytes(StandardCharsets.UTF_8));
        }

        // Create file with openFile (create if not exists)
        jetbrains.exodus.vfs.File file2 = vfs.openFile(txn, "demo/test2.txt", true);
        System.out.println("Opened/created file: " + file2.getPath());

        // Write content to the second file
        try (OutputStream output = vfs.writeFile(txn, file2)) {
          String content = "This is another test file.\n";
          output.write(content.getBytes(StandardCharsets.UTF_8));
        }

        // Append content to the second file
        try (OutputStream output = vfs.appendFile(txn, file2)) {
          String appendContent = "Appending more content to test2.txt.\n";
          output.write(appendContent.getBytes(StandardCharsets.UTF_8));
        }

        // Create file with unique auto-generated path
        jetbrains.exodus.vfs.File uniqueFile = vfs.createUniqueFile(txn, "temp/unique_");
        System.out.println("Created unique file: " + uniqueFile.getPath());

        // Check if file exists by trying to open it (will return null if doesn't
        // exist)
        jetbrains.exodus.vfs.File file3 = vfs.openFile(txn, "demo/test1.txt", false);
        if (file3 != null) {
          System.out.println("File demo/test1.txt exists");
        }
        else {
          System.out.println("File demo/test1.txt does not exist");
        }

        // Rename a file
        vfs.renameFile(txn, file2, "demo/renamed_test2.txt");

        // Delete a file
        vfs.deleteFile(txn, "demo/test1.txt");

        // Return total size of all files in the VFS
        long size = vfs.diskUsage(txn);
        System.out.println("Total disk usage: " + size + " bytes");

        // List all files
        System.out.println("Files:");
        for (jetbrains.exodus.vfs.File file : vfs.getFiles(txn)) {
          System.out.println(
              " - " + file.getPath() + " (descriptor: " + file.getDescriptor() + ")");
        }

        // Read content from a file
        try (InputStream input = vfs.readFile(txn, file2)) {
          byte[] buffer = input.readAllBytes();
          String readContent = new String(buffer, StandardCharsets.UTF_8);
          System.out.println("Content of renamed_test2.txt:");
          System.out.println(readContent);
        }

      }
      catch (Exception e) {
        System.err.println("Error in file operations: " + e.getMessage());
      }
    });
  }

  private static void demonstrateFilePositioning(Environment env, VirtualFileSystem vfs) {
    System.out.println("\nFile Positioning Demo");

    env.executeInTransaction(txn -> {
      try {
        jetbrains.exodus.vfs.File file = vfs.openFile(txn, "demo/positioned.txt", true);

        String initialContent = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ";
        try (OutputStream output = vfs.writeFile(txn, file)) {
          output.write(initialContent.getBytes(StandardCharsets.UTF_8));
        }

        // Read from specific position
        long position = 10;
        try (InputStream input = vfs.readFile(txn, file, position)) {
          byte[] buffer = new byte[5];
          int bytesRead = input.read(buffer);
          if (bytesRead > 0) {
            String readFromPosition = new String(buffer, 0, bytesRead,
                StandardCharsets.UTF_8);
            System.out
                .println("Read from position " + position + ": " + readFromPosition);
          }
        }

        // Write at specific position (overwrite)
        String overwriteContent = "XXXXX";
        try (OutputStream output = vfs.writeFile(txn, file, 15)) {
          output.write(overwriteContent.getBytes(StandardCharsets.UTF_8));
        }

        // Read entire file to see the change
        try (InputStream input = vfs.readFile(txn, file)) {
          byte[] buffer = input.readAllBytes();
          String finalContent = new String(buffer, StandardCharsets.UTF_8);
          System.out.println("Final content after overwrite: " + finalContent);
        }

        // Get file length
        long length = vfs.getFileLength(txn, file);
        System.out.println("File length: " + length + " bytes");

      }
      catch (IOException e) {
        System.err.println("Error in positioning operations: " + e.getMessage());
      }
    });
  }
}
