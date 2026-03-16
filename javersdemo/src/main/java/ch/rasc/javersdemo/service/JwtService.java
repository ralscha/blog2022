package ch.rasc.javersdemo.service;

import java.security.SecureRandom;
import java.time.Instant;
import java.util.Date;
import java.util.Map;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.stereotype.Service;

import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.auth0.jwt.interfaces.JWTVerifier;

@Service
public class JwtService {

  @Value("${jwt.expiration:86400000}") // 24 hours in milliseconds
  private int jwtExpiration;

  private final Algorithm algorithm;
  private final JWTVerifier verifier;

  public JwtService() {
    byte[] secret = new byte[64];
    new SecureRandom().nextBytes(secret);
    this.algorithm = Algorithm.HMAC512(secret);
    this.verifier = JWT.require(this.algorithm).build();
  }

  public String extractUsername(String token) {
    return verifier.verify(token).getSubject();
  }

  public String generateToken(UserDetails userDetails) {
    return generateToken(Map.of(), userDetails);
  }

  public String generateToken(Map<String, Object> extraClaims, UserDetails userDetails) {
    return buildToken(extraClaims, userDetails, this.jwtExpiration);
  }

  private String buildToken(Map<String, Object> extraClaims, UserDetails userDetails,
      long expiration) {
    return JWT.create().withPayload(extraClaims).withSubject(userDetails.getUsername())
        .withIssuedAt(Instant.now())
        .withExpiresAt(Instant.now().plusMillis(expiration))
        .sign(this.algorithm);
  }

  public boolean isTokenValid(String token, UserDetails userDetails) {
    final String username = extractUsername(token);
    return username.equals(userDetails.getUsername()) && !isTokenExpired(token);
  }

  private boolean isTokenExpired(String token) {
    return JWT.decode(token).getExpiresAt().before(new Date());
  }

}
