package ch.rasc.javersdemo.entity;

import java.time.LocalDateTime;

import org.javers.core.metamodel.annotation.Id;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public class User {
  @Id
  private Long id;

  @NotBlank
  @Size(min = 3, max = 50)
  private String username;

  @NotBlank
  private String password;

  private Role role = Role.USER;

  private Boolean enabled = true;

  private LocalDateTime createdAt;

  private LocalDateTime updatedAt;

  public enum Role {
    USER, ADMIN
  }

  // Constructors
  public User() {
  }

  public User(String username, String password, Role role) {
    this.username = username;
    this.password = password;
    this.role = role;
    this.createdAt = LocalDateTime.now();
    this.updatedAt = LocalDateTime.now();
  }

  // Getters and Setters
  public Long getId() {
    return this.id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public String getUsername() {
    return this.username;
  }

  public void setUsername(String username) {
    this.username = username;
  }

  public String getPassword() {
    return this.password;
  }

  public void setPassword(String password) {
    this.password = password;
  }

  public Role getRole() {
    return this.role;
  }

  public void setRole(Role role) {
    this.role = role;
  }

  public Boolean getEnabled() {
    return this.enabled;
  }

  public void setEnabled(Boolean enabled) {
    this.enabled = enabled;
  }

  public LocalDateTime getCreatedAt() {
    return this.createdAt;
  }

  public void setCreatedAt(LocalDateTime createdAt) {
    this.createdAt = createdAt;
  }

  public LocalDateTime getUpdatedAt() {
    return this.updatedAt;
  }

  public void setUpdatedAt(LocalDateTime updatedAt) {
    this.updatedAt = updatedAt;
  }
}
