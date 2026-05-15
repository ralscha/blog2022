package ch.rasc.shedlockdemo.service;

import org.jooq.DSLContext;
import org.jooq.Field;
import org.jooq.Table;
import static org.jooq.impl.DSL.field;
import static org.jooq.impl.DSL.name;
import static org.jooq.impl.DSL.table;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import net.javacrumbs.shedlock.core.LockAssert;
import net.javacrumbs.shedlock.spring.annotation.SchedulerLock;

@Service
public class NotificationService {

  private static final Logger logger = LoggerFactory.getLogger(NotificationService.class);
  private static final Table<?> APP_USER = table(name("app_user"));
  private static final Field<String> APP_USER_EMAIL = field(name("email"), String.class);

  private final DSLContext dsl;

  public NotificationService(DSLContext dsl) {
    this.dsl = dsl;
  }

  @Scheduled(fixedDelay = 60_000)
  @Transactional
  public void processNotificationsWithNoLock() {
    logger.info("Starting process notifications with no lock");
    this.dsl.selectFrom(APP_USER).where(APP_USER_EMAIL.isNotNull())
        .forEach(appUser -> {
          sendNotification(appUser.get(APP_USER_EMAIL));
        });
  }

  @Scheduled(fixedDelay = 300_000)
  @SchedulerLock(name = "processNotifications", lockAtMostFor = "4m")
  @Transactional
  public void processNotifications() {
    logger.info("Starting process notifications");

    LockAssert.assertLocked();

    this.dsl.selectFrom(APP_USER).where(APP_USER_EMAIL.isNotNull())
        .forEach(appUser -> {
          sendNotification(appUser.get(APP_USER_EMAIL));
        });

  }

  private void sendNotification(String email) {
    logger.debug("Sending notification to {}", email);
  }
}
