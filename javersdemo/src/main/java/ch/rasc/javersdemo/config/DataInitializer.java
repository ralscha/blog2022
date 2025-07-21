package ch.rasc.javersdemo.config;

import org.springframework.boot.CommandLineRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.crypto.password.PasswordEncoder;

import ch.rasc.javersdemo.entity.User;
import ch.rasc.javersdemo.service.UserService;

@Configuration
public class DataInitializer {

  @Bean
  CommandLineRunner init(UserService userService, PasswordEncoder passwordEncoder) {
    return _ -> {
      if (userService.getUserCount() == 0) {
        String encodedPassword = passwordEncoder.encode("password");
        userService.createUser(new User("user1", encodedPassword, User.Role.USER));
        userService.createUser(new User("user2", encodedPassword, User.Role.USER));
        userService.createUser(new User("admin", encodedPassword, User.Role.ADMIN));
        System.out.println("Sample users created with password: password");
      }
    };
  }
}
