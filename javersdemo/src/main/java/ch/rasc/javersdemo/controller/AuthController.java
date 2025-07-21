package ch.rasc.javersdemo.controller;

import java.util.HashMap;
import java.util.Map;

import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import ch.rasc.javersdemo.service.JwtService;

@RestController
@RequestMapping("/api/auth")
public class AuthController {

  private final AuthenticationManager authenticationManager;
  private final JwtService jwtService;

  public AuthController(AuthenticationManager authenticationManager,
      JwtService jwtService) {
    this.authenticationManager = authenticationManager;
    this.jwtService = jwtService;
  }

  @PostMapping("/authenticate")
  public ResponseEntity<Map<String, String>> authenticate(
      @RequestBody AuthRequest request) {
    try {
      Authentication authentication = this.authenticationManager
          .authenticate(new UsernamePasswordAuthenticationToken(request.username(),
              request.password()));

      UserDetails userDetails = (UserDetails) authentication.getPrincipal();
      String token = this.jwtService.generateToken(userDetails);

      Map<String, String> response = new HashMap<>();
      response.put("token", token);
      response.put("username", userDetails.getUsername());
      response.put("message", "Authentication successful");

      return ResponseEntity.ok(response);
    }
    catch (Exception e) {
      Map<String, String> errorResponse = new HashMap<>();
      errorResponse.put("error", "Authentication failed");
      errorResponse.put("message", "Bad Credential");
      return ResponseEntity.status(401).body(errorResponse);
    }
  }

  public record AuthRequest(String username, String password) {
  }
}
