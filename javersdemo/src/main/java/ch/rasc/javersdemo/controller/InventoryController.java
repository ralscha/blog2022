package ch.rasc.javersdemo.controller;

import java.util.List;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PatchMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ch.rasc.javersdemo.entity.InventoryItem;
import ch.rasc.javersdemo.service.InventoryService;
import jakarta.validation.Valid;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;

@RestController
@RequestMapping("/api/inventory")
public class InventoryController {

  private final InventoryService inventoryService;

  public InventoryController(InventoryService inventoryService) {
    this.inventoryService = inventoryService;
  }

  @GetMapping("/items")
  public ResponseEntity<List<InventoryItem>> getAllItems() {
    List<InventoryItem> items = this.inventoryService.getAllItems();
    return ResponseEntity.ok(items);
  }

  @GetMapping("/items/{id}")
  public ResponseEntity<InventoryItem> getItemById(@PathVariable Long id) {
    InventoryItem item = this.inventoryService.getItemById(id);
    if (item == null) {
      return ResponseEntity.notFound().build();
    }
    return ResponseEntity.ok(item);
  }

  @PostMapping("/items")
  public ResponseEntity<InventoryItem> createItem(
      @Valid @RequestBody InventoryItem item) {
    InventoryItem savedItem = this.inventoryService.createItem(item);
    return ResponseEntity.status(HttpStatus.CREATED).body(savedItem);
  }

  @PutMapping("/items/{id}")
  public ResponseEntity<InventoryItem> updateItem(@PathVariable Long id,
      @Valid @RequestBody InventoryItem item) {
    InventoryItem existingItem = this.inventoryService.getItemById(id);
    if (existingItem == null) {
      return ResponseEntity.notFound().build();
    }

    InventoryItem updatedItem = this.inventoryService.updateItem(item);
    return ResponseEntity.ok(updatedItem);
  }

  @PatchMapping("/items/{id}/quantity")
  public ResponseEntity<InventoryItem> updateQuantity(@PathVariable Long id,
      @RequestBody QuantityUpdate quantityUpdate) {
    try {
      InventoryItem updatedItem = this.inventoryService.updateQuantity(id,
          quantityUpdate.getQuantity());
      return ResponseEntity.ok(updatedItem);
    }
    catch (RuntimeException e) {
      return ResponseEntity.notFound().build();
    }
  }

  @DeleteMapping("/items/{id}")
  public ResponseEntity<Void> deleteItem(@PathVariable Long id) {
    if (this.inventoryService.getItemById(id) == null) {
      return ResponseEntity.notFound().build();
    }
    this.inventoryService.deleteItem(id);
    return ResponseEntity.noContent().build();
  }

  // Inner class for quantity update requests
  public static class QuantityUpdate {
    @NotNull(message = "Quantity is required")
    @Min(value = 0, message = "Quantity must be non-negative")
    private Integer quantity;

    public QuantityUpdate() {
    }

    public QuantityUpdate(Integer quantity) {
      this.quantity = quantity;
    }

    public Integer getQuantity() {
      return this.quantity;
    }

    public void setQuantity(Integer quantity) {
      this.quantity = quantity;
    }
  }
}
