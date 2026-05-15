package ch.rasc.javersdemo.service;

import java.math.BigDecimal;
import java.time.LocalDateTime;
import java.util.List;

import org.jooq.DSLContext;
import org.jooq.Field;
import org.jooq.Table;
import static org.jooq.impl.DSL.field;
import static org.jooq.impl.DSL.name;
import static org.jooq.impl.DSL.table;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import ch.rasc.javersdemo.entity.InventoryItem;

@Service
@Transactional
public class InventoryService {

  private static final Table<?> INVENTORY_ITEMS = table(name("inventory_items"));
  private static final Field<Long> INVENTORY_ITEMS_ID = field(name("id"), Long.class);
  private static final Field<String> INVENTORY_ITEMS_NAME = field(name("name"), String.class);
  private static final Field<String> INVENTORY_ITEMS_DESCRIPTION =
    field(name("description"), String.class);
  private static final Field<Integer> INVENTORY_ITEMS_QUANTITY =
    field(name("quantity"), Integer.class);
  private static final Field<BigDecimal> INVENTORY_ITEMS_PRICE =
    field(name("price"), BigDecimal.class);
  private static final Field<String> INVENTORY_ITEMS_CATEGORY =
    field(name("category"), String.class);
  private static final Field<String> INVENTORY_ITEMS_SKU = field(name("sku"), String.class);
  private static final Field<LocalDateTime> INVENTORY_ITEMS_CREATED_AT =
    field(name("created_at"), LocalDateTime.class);
  private static final Field<LocalDateTime> INVENTORY_ITEMS_UPDATED_AT =
    field(name("updated_at"), LocalDateTime.class);

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
    return this.dsl.selectFrom(INVENTORY_ITEMS).where(INVENTORY_ITEMS_ID.eq(id))
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
        .set(INVENTORY_ITEMS_NAME, item.getName())
        .set(INVENTORY_ITEMS_DESCRIPTION, item.getDescription())
        .set(INVENTORY_ITEMS_QUANTITY, item.getQuantity())
        .set(INVENTORY_ITEMS_PRICE, item.getPrice())
        .set(INVENTORY_ITEMS_CATEGORY, item.getCategory())
        .set(INVENTORY_ITEMS_SKU, item.getSku())
        .set(INVENTORY_ITEMS_CREATED_AT,
              item.getCreatedAt() != null ? item.getCreatedAt() : LocalDateTime.now())
        .set(INVENTORY_ITEMS_UPDATED_AT, LocalDateTime.now())
        .returning(INVENTORY_ITEMS_ID).fetchOne();

      item.setId(record.get(INVENTORY_ITEMS_ID));
      return item;
    }
    // Update existing item
    this.dsl.update(INVENTORY_ITEMS).set(INVENTORY_ITEMS_NAME, item.getName())
      .set(INVENTORY_ITEMS_DESCRIPTION, item.getDescription())
      .set(INVENTORY_ITEMS_QUANTITY, item.getQuantity())
      .set(INVENTORY_ITEMS_PRICE, item.getPrice())
      .set(INVENTORY_ITEMS_CATEGORY, item.getCategory())
      .set(INVENTORY_ITEMS_SKU, item.getSku())
      .set(INVENTORY_ITEMS_UPDATED_AT, LocalDateTime.now())
      .where(INVENTORY_ITEMS_ID.eq(item.getId())).execute();

    return item;
  }

  private void deleteById(Long id) {
    this.dsl.deleteFrom(INVENTORY_ITEMS).where(INVENTORY_ITEMS_ID.eq(id)).execute();
  }

  public long count() {
    return this.dsl.selectCount().from(INVENTORY_ITEMS).fetchOne(0, long.class);
  }

}
