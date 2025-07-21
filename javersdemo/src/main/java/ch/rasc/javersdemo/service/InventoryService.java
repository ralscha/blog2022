package ch.rasc.javersdemo.service;

import static ch.rasc.javersdemo.db.tables.InventoryItems.INVENTORY_ITEMS;

import java.time.LocalDateTime;
import java.util.List;

import org.jooq.DSLContext;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import ch.rasc.javersdemo.entity.InventoryItem;

@Service
@Transactional
public class InventoryService {

  private final DSLContext dsl;
  private final AuditService auditService;

  public InventoryService(DSLContext dsl, AuditService auditService) {
    this.dsl = dsl;
    this.auditService = auditService;
  }

  public List<InventoryItem> getAllItems() {
    return this.dsl.selectFrom(INVENTORY_ITEMS).fetch().into(InventoryItem.class);
  }

  public InventoryItem getItemById(Long id) {
    return this.dsl.selectFrom(INVENTORY_ITEMS).where(INVENTORY_ITEMS.ID.eq(id))
        .fetchOne().into(InventoryItem.class);
  }

  public InventoryItem createItem(InventoryItem item) {
    InventoryItem savedItem = save(item);
    this.auditService.commitEntity(savedItem);
    return savedItem;
  }

  public InventoryItem updateItem(InventoryItem item) {
    InventoryItem savedItem = save(item);
    this.auditService.commitEntity(savedItem);
    return savedItem;
  }

  public void deleteItem(Long id) {
    InventoryItem item = getItemById(id);
    if (item != null) {
      deleteById(id);
      this.auditService.commitEntityDeletion(item);
    }
  }

  public InventoryItem updateQuantity(Long id, Integer newQuantity) {
    InventoryItem item = getItemById(id);
    if (item == null) {
      throw new RuntimeException("Item not found");
    }
    item.setQuantity(newQuantity);
    InventoryItem savedItem = save(item);
    this.auditService.commitEntity(savedItem);
    return savedItem;
  }

  // Repository methods merged into service
  private InventoryItem save(InventoryItem item) {
    item.setUpdatedAt(LocalDateTime.now());
    if (item.getId() == null) {
      // Insert new item
      var record = this.dsl.insertInto(INVENTORY_ITEMS)
          .set(INVENTORY_ITEMS.NAME, item.getName())
          .set(INVENTORY_ITEMS.DESCRIPTION, item.getDescription())
          .set(INVENTORY_ITEMS.QUANTITY, item.getQuantity())
          .set(INVENTORY_ITEMS.PRICE, item.getPrice())
          .set(INVENTORY_ITEMS.CATEGORY, item.getCategory())
          .set(INVENTORY_ITEMS.SKU, item.getSku())
          .set(INVENTORY_ITEMS.CREATED_AT,
              item.getCreatedAt() != null ? item.getCreatedAt() : LocalDateTime.now())
          .set(INVENTORY_ITEMS.UPDATED_AT, LocalDateTime.now())
          .returning(INVENTORY_ITEMS.ID).fetchOne();

      item.setId(record.getId());
      return item;
    }
    // Update existing item
    this.dsl.update(INVENTORY_ITEMS).set(INVENTORY_ITEMS.NAME, item.getName())
        .set(INVENTORY_ITEMS.DESCRIPTION, item.getDescription())
        .set(INVENTORY_ITEMS.QUANTITY, item.getQuantity())
        .set(INVENTORY_ITEMS.PRICE, item.getPrice())
        .set(INVENTORY_ITEMS.CATEGORY, item.getCategory())
        .set(INVENTORY_ITEMS.SKU, item.getSku())
        .set(INVENTORY_ITEMS.UPDATED_AT, LocalDateTime.now())
        .where(INVENTORY_ITEMS.ID.eq(item.getId())).execute();

    return item;
  }

  private void deleteById(Long id) {
    this.dsl.deleteFrom(INVENTORY_ITEMS).where(INVENTORY_ITEMS.ID.eq(id)).execute();
  }

  public long count() {
    return this.dsl.selectCount().from(INVENTORY_ITEMS).fetchOne(0, long.class);
  }

}
