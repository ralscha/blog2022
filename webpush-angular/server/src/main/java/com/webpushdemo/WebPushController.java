package com.webpushdemo;

import java.io.IOException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;

import org.springframework.http.HttpStatus;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.support.WebClientAdapter;
import org.springframework.web.service.invoker.HttpServiceProxyFactory;

import com.webpushdemo.ChuckNorrisJokeService.Joke;
import com.zerodeplibs.webpush.PushSubscription;
import com.zerodeplibs.webpush.httpclient.StandardHttpClientRequestPreparer;

@RestController
public class WebPushController {

  private final Map<String, PushSubscription> pushSubscriptions = new ConcurrentHashMap<>();

  private final WebPushService webPushService;

  private final ChuckNorrisJokeService chuckNorrisJokeService;

  private final HttpClient httpClient;

  public WebPushController(WebPushService webPushService) {
    this.webPushService = webPushService;

    WebClient client = WebClient.builder().baseUrl("https://api.chucknorris.io").build();
    HttpServiceProxyFactory factory = HttpServiceProxyFactory
        .builder(WebClientAdapter.forClient(client)).build();
    this.chuckNorrisJokeService = factory.createClient(ChuckNorrisJokeService.class);

    this.httpClient = HttpClient.newHttpClient();
  }

  @GetMapping(path = "/publicKey")
  public String publicKey() {
    return this.webPushService.getPublicKey();
  }

  @PostMapping("/subscribe")
  @ResponseStatus(HttpStatus.CREATED)
  public void subscribe(@RequestBody PushSubscription subscription) {
    Application.logger.info("subscribe: " + subscription);
    this.pushSubscriptions.put(subscription.getEndpoint(), subscription);
  }

  @PostMapping("/unsubscribe")
  @ResponseStatus(HttpStatus.NO_CONTENT)
  public void unsubscribe(@RequestBody PushSubscription subscription) {
    Application.logger.info("unsubscribe: " + subscription);
    this.pushSubscriptions.remove(subscription.getEndpoint());
  }

  @Scheduled(fixedDelayString = "PT1M")
  public void sendJokes() {
    if (this.pushSubscriptions.isEmpty()) {
      return;
    }

    Joke joke = this.chuckNorrisJokeService.getRandomJoke();

    Application.logger.info("sending joke to subscribers: {}", joke.id());

    String msg = """
        {
          "notification": {
             "title": "{title}",
             "body": "{body}",
             "icon": "assets/icons/icon-72x72.png",
             "data": {
               "onActionClick": {
                 "default": {"operation": "navigateLastFocusedOrOpen", "url": "/"},               }
             }
          }
        }
        """
        .replace("{title}", "Chuck Norris Joke").replace("{body}", joke.value());

    for (PushSubscription subscription : this.pushSubscriptions.values()) {
      HttpRequest request = StandardHttpClientRequestPreparer.getBuilder()
          .pushSubscription(subscription).vapidJWTExpiresAfter(3, TimeUnit.HOURS)
          .vapidJWTSubject("mailto:example@example.com").pushMessage(msg)
          .ttl(1, TimeUnit.HOURS).urgencyNormal().topic("Joke")
          .build(this.webPushService.getKeyPair()).toRequest();

      try {
        HttpResponse<String> httpResponse = this.httpClient.send(request,
            HttpResponse.BodyHandlers.ofString());

        switch (httpResponse.statusCode()) {
        case 201 -> {
          Application.logger.info("Push message successfully sent: {}",
              httpResponse.body());
        }
        case 404, 410 -> {
          Application.logger.warn("Subscription not found or gone: {}",
              subscription.getEndpoint());
          // remove subscription
          this.pushSubscriptions.remove(subscription.getEndpoint());
        }
        case 429 -> {
          Application.logger.error("Too many requests: {}", request);
          // TODO: retry
        }
        case 400 -> {
          Application.logger.error("Invalid request: {}", request);
          // TODO: something is wrong with the request
        }
        case 413 -> {
          Application.logger.error("Payload size too large: {}", request);
          // TODO: decrease payload
        }
        default -> {
          Application.logger.error("Unhandled status code: {} / {}",
              httpResponse.statusCode(), request);
          // TODO: might be a temporary problem with the push service. retry
        }
        }

      }
      catch (IOException | InterruptedException e) {
        Application.logger.error("sending to push notification failed", e);
      }

    }

  }

}
