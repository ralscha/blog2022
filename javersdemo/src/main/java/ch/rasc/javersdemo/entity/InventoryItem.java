package ch.rasc.javersdemo.entity;

import java.math.BigDecimal;
import java.time.LocalDateTime;

import org.javers.core.metamodel.annotation.DiffIgnore;
import org.javers.core.metamodel.annotation.Id;

import jakarta.validation.constraints.DecimalMin;
import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;

public class InventoryItem {
  @Id
  private Long id;

  @NotBlank
  private String name;

  private String description;

  @NotNull
  @Min(0)
  private Integer quantity = 0;

  @NotNull
  @DecimalMin("0.00")
  private BigDecimal price = BigDecimal.ZERO;

  private String category;

  private String sku;

  @DiffIgnore
  private LocalDateTime createdAt;

  @DiffIgnore
  private LocalDateTime updatedAt;

  // Constructors
  public InventoryItem() {
    this.createdAt = LocalDateTime.now();
    this.updatedAt = LocalDateTime.now();
  }

  public InventoryItem(String name, String description, Integer quantity,
      BigDecimal price, String category, String sku) {
    this.name = name;
    this.description = description;
    this.quantity = quantity;
    this.price = price;
    this.category = category;
    this.sku = sku;
    this.createdAt = LocalDateTime.now();
    this.updatedAt = LocalDateTime.now();
  }

  // Getters and Setters
  public Long getId() {
    return this.id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public String getName() {
    return this.name;
  }

  public void setName(String name) {
    this.name = name;
  }

  public String getDescription() {
    return this.description;
  }

  public void setDescription(String description) {
    this.description = description;
  }

  public Integer getQuantity() {
    return this.quantity;
  }

  public void setQuantity(Integer quantity) {
    this.quantity = quantity;
  }

  public BigDecimal getPrice() {
    return this.price;
  }

  public void setPrice(BigDecimal price) {
    this.price = price;
  }

  public String getCategory() {
    return this.category;
  }

  public void setCategory(String category) {
    this.category = category;
  }

  public String getSku() {
    return this.sku;
  }

  public void setSku(String sku) {
    this.sku = sku;
  }

  public LocalDateTime getCreatedAt() {
    return this.createdAt;
  }

  public void setCreatedAt(LocalDateTime createdAt) {
    this.createdAt = createdAt;
  }

  public LocalDateTime getUpdatedAt() {
    return this.updatedAt;
  }

  public void setUpdatedAt(LocalDateTime updatedAt) {
    this.updatedAt = updatedAt;
  }
}
