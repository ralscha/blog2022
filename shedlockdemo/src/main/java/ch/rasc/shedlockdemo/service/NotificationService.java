package ch.rasc.shedlockdemo.service;

import org.jooq.DSLContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import ch.rasc.shedlockdemo.db.tables.AppUser;
import net.javacrumbs.shedlock.core.LockAssert;
import net.javacrumbs.shedlock.spring.annotation.SchedulerLock;

@Service
public class NotificationService {

  private static final Logger logger = LoggerFactory.getLogger(NotificationService.class);

  private final DSLContext dsl;

  public NotificationService(DSLContext dsl) {
    this.dsl = dsl;
  }

  @Scheduled(fixedDelay = 60_000)
  @Transactional
  public void processNotificationsWithNoLock() {
    logger.info("Starting process notifications with no lock");
    this.dsl.selectFrom(AppUser.APP_USER).where(AppUser.APP_USER.EMAIL.isNotNull())
        .forEach(appUser -> {
          sendNotification(appUser.getEmail());
        });
  }

  @Scheduled(fixedDelay = 300_000)
  @SchedulerLock(name = "processNotifications", lockAtMostFor = "4m")
  @Transactional
  public void processNotifications() {
    logger.info("Starting process notifications");

    LockAssert.assertLocked();

    this.dsl.selectFrom(AppUser.APP_USER).where(AppUser.APP_USER.EMAIL.isNotNull())
        .forEach(appUser -> {
          sendNotification(appUser.getEmail());
        });

  }

  private void sendNotification(String email) {
    logger.debug("Sending notification to {}", email);
  }
}
