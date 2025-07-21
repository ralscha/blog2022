package ch.rasc.javersdemo.service;

import static ch.rasc.javersdemo.db.tables.Users.USERS;

import java.time.LocalDateTime;

import org.jooq.DSLContext;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import ch.rasc.javersdemo.entity.User;

@Service
@Transactional
public class UserService {

  private final DSLContext dsl;

  public UserService(DSLContext dsl) {
    this.dsl = dsl;
  }

  public User findById(Long id) {
    return this.dsl.selectFrom(USERS).where(USERS.ID.eq(id)).fetchOne().into(User.class);
  }

  public User findByUsername(String username) {
    return this.dsl.selectFrom(USERS).where(USERS.USERNAME.eq(username)).fetchOne()
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
      var record = this.dsl.insertInto(USERS).set(USERS.USERNAME, user.getUsername())
          .set(USERS.PASSWORD, user.getPassword()).set(USERS.ROLE, user.getRole().name())
          .set(USERS.ENABLED, user.getEnabled())
          .set(USERS.CREATED_AT,
              user.getCreatedAt() != null ? user.getCreatedAt() : LocalDateTime.now())
          .set(USERS.UPDATED_AT, LocalDateTime.now()).returning(USERS.ID).fetchOne();

      user.setId(record.getId());
      return user;
    }

    // Update existing user
    this.dsl.update(USERS).set(USERS.USERNAME, user.getUsername())
        .set(USERS.PASSWORD, user.getPassword()).set(USERS.ROLE, user.getRole().name())
        .set(USERS.ENABLED, user.getEnabled()).set(USERS.UPDATED_AT, LocalDateTime.now())
        .where(USERS.ID.eq(user.getId())).execute();

    return user;
  }

}
