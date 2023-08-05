package com.webpushdemo;

import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;

import org.springframework.stereotype.Service;
import org.springframework.util.StreamUtils;

import com.zerodeplibs.webpush.VAPIDKeyPair;
import com.zerodeplibs.webpush.VAPIDKeyPairs;
import com.zerodeplibs.webpush.key.PrivateKeySources;
import com.zerodeplibs.webpush.key.PublicKeySources;

@Service
public class WebPushService {

  private final VAPIDKeyPair vapidKeyPair;

  public WebPushService() {
    String privateKey = null;
    String publicKey = null;

    try (InputStream privateIs = getClass().getResourceAsStream("/vapidPrivateKey.pem")) {
      privateKey = StreamUtils.copyToString(privateIs, StandardCharsets.UTF_8);
    }
    catch (IOException e) {
      Application.logger.error("can't load vapid private key", e);
    }
    try (InputStream publicIs = getClass().getResourceAsStream("/vapidPublicKey.pem")) {
      publicKey = StreamUtils.copyToString(publicIs, StandardCharsets.UTF_8);
    }
    catch (IOException e) {
      Application.logger.error("can't load vapid public key", e);
    }

    this.vapidKeyPair = VAPIDKeyPairs.of(PrivateKeySources.ofPEMText(privateKey),
        PublicKeySources.ofPEMText(publicKey));
  }

  public String getPublicKey() {
    return this.vapidKeyPair.extractPublicKeyInUncompressedFormAsString();
  }

  public VAPIDKeyPair getKeyPair() {
    return this.vapidKeyPair;
  }

}