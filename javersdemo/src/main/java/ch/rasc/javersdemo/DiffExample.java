package ch.rasc.javersdemo;

import java.util.Arrays;
import java.util.List;

import org.javers.core.Javers;
import org.javers.core.JaversBuilder;
import org.javers.core.diff.Change;
import org.javers.core.diff.Diff;
import org.javers.core.diff.changetype.NewObject;
import org.javers.core.diff.changetype.ReferenceChange;
import org.javers.core.diff.changetype.ValueChange;
import org.javers.core.diff.changetype.container.ContainerElementChange;
import org.javers.core.diff.changetype.container.ElementValueChange;
import org.javers.core.diff.changetype.container.ListChange;
import org.javers.core.diff.changetype.container.ValueAdded;
import org.javers.core.diff.changetype.container.ValueRemoved;
import org.javers.core.metamodel.annotation.Id;

public class DiffExample {

  record User(@Id String id, String name, int age, List<String> roles, Address address,
      Todo todo) {
  }

  record Address(String street, String city) {
  }

  record Todo(@Id String id, String title, boolean completed) {
  }

  public static void main(String[] args) {
    Javers javers = JaversBuilder.javers().build();
    // Javers javers = JaversBuilder.javers().withInitialChanges(false).build();

    Address address = new Address("123 Main St", "Anytown");
    Address newAddress = new Address("1234 Main St", "Anytown");
    Diff diff = javers.compare(address, newAddress);
    for (Change change : diff.getChanges()) {
      ValueChange valueChange = (ValueChange) change;
      System.out.println("Property '" + valueChange.getPropertyName() + "' changed from '"
          + valueChange.getLeft() + "' to '" + valueChange.getRight() + "'");
    }
    System.out.println();

    List<String> h1 = List.of("admin", "editor", "viewer", "reporter");
    List<String> h2 = List.of("admin", "viewer", "reporter");
    diff = javers.compareCollections(h1, h2, String.class);
    for (Change change : diff.getChanges()) {
      ListChange listChange = (ListChange) change;

      for (ContainerElementChange c : listChange.getChanges()) {
        switch (c) {
        case ElementValueChange evc -> System.out
            .println("Value changed: " + evc.getIndex() + " from '" + evc.getLeftValue()
                + "' to '" + evc.getRightValue() + "'");
        case ValueAdded va -> System.out
            .println("Added: " + va.getValue() + " at index " + va.getIndex());
        case ValueRemoved vr -> System.out
            .println("Removed: " + vr.getValue() + " at index " + vr.getIndex());
        default -> throw new IllegalArgumentException("Unexpected value: " + c);
        }
      }
    }
    System.out.println();

    // --- Demonstrate ValueChange ---
    User user = new User("U1", "Alice", 30, List.of("admin", "editor"), address, null);
    User userAfterValueChange = new User("U1", "Alicia", 31, user.roles(), user.address(),
        null);

    diff = javers.compare(user, userAfterValueChange);
    System.out.println("Has Changes: " + diff.hasChanges());
    System.out.println("Diff (Value Changes): \n" + diff.prettyPrint());

    for (Change change : diff.getChanges()) {
      ValueChange vc = (ValueChange) change;
      System.out.println("  ValueChange: Property '" + vc.getPropertyName()
          + "' changed from '" + vc.getLeft() + "' to '" + vc.getRight() + "'");
    }
    System.out.println();

    // --- Demonstrate ReferenceChange ---
    newAddress = new Address("456 Oak Ave", "Newville");
    Todo todo1 = new Todo("T1", "Buy groceries", false);
    User userAfterAddressAndTodoChange = new User("U1", userAfterValueChange.name(),
        userAfterValueChange.age(), user.roles(), newAddress, todo1);

    diff = javers.compare(userAfterValueChange, userAfterAddressAndTodoChange);
    for (Change change : diff.getChanges()) {
      switch (change) {
      case NewObject noc -> System.out.println("  New Object: " + noc);
      case ReferenceChange rc -> System.out
          .println("  ReferenceChange: Property '" + rc.getPropertyName()
              + "' changed from " + rc.getLeft() + " to " + rc.getRight());
      case ValueChange vc -> System.out
          .println("  ValueChange: Property '" + vc.getPropertyName() + "' changed from '"
              + vc.getLeft() + "' to '" + vc.getRight() + "'");
      default -> System.out.println("  Other Change: " + change);
      }
    }
    System.out.println();

    // --- Demonstrate ListChange ---
    List<String> updatedRoles = Arrays.asList("admin", "viewer", "reporter");
    User userAfterRoleChange = new User("U1", userAfterAddressAndTodoChange.name(),
        userAfterAddressAndTodoChange.age(), updatedRoles,
        userAfterAddressAndTodoChange.address(), todo1);

    diff = javers.compare(userAfterAddressAndTodoChange, userAfterRoleChange);
    System.out.println("Diff (List Change): \n" + diff.prettyPrint());

    for (Change change : diff.getChanges()) {
      ListChange lc = (ListChange) change;
      System.out.println("  ListChange on property '" + lc.getPropertyName() + "':");
      lc.getChanges().forEach(itemChange -> {
        System.out.println("    " + itemChange);
      });
    }
    System.out.println();

    // --- Demonstrate different Id ---
    User u1 = new User("U1", "Alice", 30, List.of("admin", "editor"), address, null);
    User u2 = new User("U2", "Alice", 30, List.of("admin", "editor"), address, null);
    diff = javers.compare(u1, u2);
    System.out.println("Diff with different IDs: \n" + diff.prettyPrint());

  }
}
