package ch.rasc.javersdemo.controller;

import java.util.List;

import org.javers.core.Javers;
import org.javers.core.metamodel.object.CdoSnapshot;
import org.javers.repository.jql.QueryBuilder;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import ch.rasc.javersdemo.entity.InventoryItem;

@RestController
@RequestMapping("/api/audit")
public class AuditController {

  private final Javers javers;

  public AuditController(Javers javers) {
    this.javers = javers;
  }


  @GetMapping("/inventory-items/{id}/snapshots")
  public String getInventoryItemSnapshots(@PathVariable Long id,
      @RequestParam(defaultValue = "10") int limit) {
    List<CdoSnapshot> changes = this.javers.findSnapshots(
        QueryBuilder.byInstanceId(id, InventoryItem.class).limit(limit).build());
    return this.javers.getJsonConverter().toJson(changes);
  }

  @GetMapping("/inventory-items/{id}/changes")
  public String getInventoryItemChanges(@PathVariable Long id) {
    var changes = this.javers
        .findChanges(QueryBuilder.byInstanceId(id, InventoryItem.class).build());
    return this.javers.getJsonConverter().toJson(changes);
  }
  
  @GetMapping("/inventory-items/{id}/shadows")
  public List<InventoryItem> getInventoryItemShadows(@PathVariable Long id) {
    return this.javers
        .findShadows(QueryBuilder.byInstanceId(id, InventoryItem.class).build()).stream()
        .map(shadow -> (InventoryItem)shadow.get()).toList();
  }

}
