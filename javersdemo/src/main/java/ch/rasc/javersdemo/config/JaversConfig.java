package ch.rasc.javersdemo.config;

import javax.sql.DataSource;

import org.javers.core.Javers;
import org.javers.repository.sql.DialectName;
import org.javers.repository.sql.JaversSqlRepository;
import org.javers.repository.sql.SqlRepositoryBuilder;
import org.javers.spring.jpa.TransactionalJpaJaversBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.jdbc.datasource.DataSourceUtils;
import org.springframework.transaction.PlatformTransactionManager;

import ch.rasc.javersdemo.entity.InventoryItem;

@Configuration
public class JaversConfig {

  private final DataSource dataSource;
  private final PlatformTransactionManager transactionManager;

  JaversConfig(DataSource dataSource, PlatformTransactionManager transactionManager) {
    this.dataSource = dataSource;
    this.transactionManager = transactionManager;
  }

  @Bean
  Javers javers() {
    JaversSqlRepository sqlRepository = SqlRepositoryBuilder.sqlRepository()
        .withConnectionProvider(() -> DataSourceUtils.getConnection(this.dataSource))
        .withDialect(DialectName.POSTGRES).build();

    return TransactionalJpaJaversBuilder.javers().withTxManager(this.transactionManager)
        .registerEntities(InventoryItem.class).registerJaversRepository(sqlRepository)
        .build();

  }

}
