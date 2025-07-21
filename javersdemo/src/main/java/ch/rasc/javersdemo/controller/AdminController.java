package ch.rasc.javersdemo.controller;

import java.security.Principal;
import java.util.HashMap;
import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/admin")
@PreAuthorize("hasRole('ADMIN')")
public class AdminController {

  @GetMapping("/test")
  public ResponseEntity<Map<String, String>> adminTest(Principal principal) {
    Map<String, String> response = new HashMap<>();
    response.put("message", "Admin access granted");
    response.put("user", principal.getName());
    response.put("role", "ADMIN");
    response.put("timestamp", String.valueOf(System.currentTimeMillis()));
    return ResponseEntity.ok(response);
  }

}
