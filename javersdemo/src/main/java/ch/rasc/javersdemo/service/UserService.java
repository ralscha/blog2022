package ch.rasc.javersdemo.service;

import java.time.LocalDateTime;

import org.jooq.DSLContext;
import org.jooq.Field;
import org.jooq.Table;
import static org.jooq.impl.DSL.field;
import static org.jooq.impl.DSL.name;
import static org.jooq.impl.DSL.table;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import ch.rasc.javersdemo.entity.User;

@Service
@Transactional
public class UserService {

  private static final Table<?> USERS = table(name("users"));
  private static final Field<Long> USERS_ID = field(name("id"), Long.class);
  private static final Field<String> USERS_USERNAME = field(name("username"), String.class);
  private static final Field<String> USERS_PASSWORD = field(name("password"), String.class);
  private static final Field<String> USERS_ROLE = field(name("role"), String.class);
  private static final Field<Boolean> USERS_ENABLED = field(name("enabled"), Boolean.class);
  private static final Field<LocalDateTime> USERS_CREATED_AT =
      field(name("created_at"), LocalDateTime.class);
  private static final Field<LocalDateTime> USERS_UPDATED_AT =
      field(name("updated_at"), LocalDateTime.class);

  private final DSLContext dsl;

  public UserService(DSLContext dsl) {
    this.dsl = dsl;
  }

  public User findById(Long id) {
    return this.dsl.selectFrom(USERS).where(USERS_ID.eq(id)).fetchOne().into(User.class);
  }

  public User findByUsername(String username) {
    return this.dsl.selectFrom(USERS).where(USERS_USERNAME.eq(username)).fetchOne()
        .into(User.class);
  }

  public User createUser(User user) {
    return save(user);
  }

  public User updateUser(User user) {
    return save(user);
  }

  public long getUserCount() {
    return this.dsl.selectCount().from(USERS).fetchOne(0, long.class);
  }

  // Repository methods merged into service
  private User save(User user) {
    if (user.getId() == null) {
      // Insert new user
      var record = this.dsl.insertInto(USERS).set(USERS_USERNAME, user.getUsername())
          .set(USERS_PASSWORD, user.getPassword()).set(USERS_ROLE, user.getRole().name())
          .set(USERS_ENABLED, user.getEnabled())
          .set(USERS_CREATED_AT,
              user.getCreatedAt() != null ? user.getCreatedAt() : LocalDateTime.now())
          .set(USERS_UPDATED_AT, LocalDateTime.now()).returning(USERS_ID).fetchOne();

      user.setId(record.get(USERS_ID));
      return user;
    }

    // Update existing user
    this.dsl.update(USERS).set(USERS_USERNAME, user.getUsername())
        .set(USERS_PASSWORD, user.getPassword()).set(USERS_ROLE, user.getRole().name())
        .set(USERS_ENABLED, user.getEnabled()).set(USERS_UPDATED_AT, LocalDateTime.now())
        .where(USERS_ID.eq(user.getId())).execute();

    return user;
  }

}
