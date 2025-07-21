package ch.rasc.javersdemo.service;

import java.util.List;

import org.javers.core.Javers;
import org.javers.core.commit.Commit;
import org.javers.core.diff.Change;
import org.javers.core.metamodel.object.CdoSnapshot;
import org.javers.repository.jql.QueryBuilder;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Service;

import ch.rasc.javersdemo.entity.InventoryItem;

@Service
public class AuditService {

  private final Javers javers;

  public AuditService(Javers javers) {
    this.javers = javers;
  }

  private static String getCurrentUser() {
    Authentication authentication = SecurityContextHolder.getContext()
        .getAuthentication();
    if (authentication != null && authentication.isAuthenticated()
        && !"anonymousUser".equals(authentication.getName())) {
      return authentication.getName();
    }
    return "system";
  }

  public Commit commitEntity(Object entity) {
    return this.javers.commit(getCurrentUser(), entity);
  }

  public Commit commitEntityDeletion(Object entity) {
    return this.javers.commitShallowDelete(getCurrentUser(), entity);
  }

}
