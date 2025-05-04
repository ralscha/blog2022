package com.example.demo;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.model.Generation;
import org.springframework.ai.tool.function.FunctionToolCallback;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import com.openmeteo.sdk.Variable;
import com.openmeteo.sdk.VariableWithValues;
import com.openmeteo.sdk.VariablesSearch;
import com.openmeteo.sdk.VariablesWithTime;
import com.openmeteo.sdk.WeatherApiResponse;

@SpringBootApplication
public class DemoApplication implements CommandLineRunner {

  public static void main(String[] args) {
    SpringApplication.run(DemoApplication.class, args);
  }

  private final ChatClient chatClient;

  public DemoApplication(ChatClient.Builder chatClientBuilder) {
    this.chatClient = chatClientBuilder.build();
  }

  @Override
  public void run(String... args) throws Exception {
    newsDemoWithoutFunction();
    newsDemo();
    temperatureDemo();
  }

  private void newsDemoWithoutFunction() {
    String prompt = "Who won the 2025 Men's ATP Tennis tournament Barcelona Open?";
    var response = this.chatClient.prompt().user(prompt).call().chatResponse();
    Generation generation = response.getResult();
    if (generation != null) {
      System.out.println(generation.getOutput().getText());
    }
    else {
      System.out.println("No generation");
    }
  }

  @Autowired
  private WikipediaArticleFetcher wikipediaArticleFetcher;

  private void newsDemo() {
    String prompt = "Who won the 2025 Men's ATP Tennis tournament Barcelona Open?";
    var response = this.chatClient.prompt().user(prompt)
        .tools(this.wikipediaArticleFetcher).call().chatResponse();
    Generation generation = response.getResult();
    if (generation != null) {
      System.out.println(generation.getOutput().getText());
    }
    else {
      System.out.println("No generation");
    }
  }

  record Location(float latitude, float longitude) {
  }

  private float fetchTemperature(Location location) {
    System.out.println("Calling fetchTemperature with parameters: " + location.latitude
        + ", " + location.longitude);
    try (var client = HttpClient.newHttpClient()) {

      var request = HttpRequest.newBuilder()
          .uri(URI.create("https://api.open-meteo.com/v1/forecast?latitude="
              + location.latitude + "&longitude=" + location.longitude
              + "&current=temperature_2m&format=flatbuffers"))
          .build();
      var response = client.send(request,
          java.net.http.HttpResponse.BodyHandlers.ofByteArray());

      ByteBuffer buffer = ByteBuffer.wrap(response.body()).order(ByteOrder.LITTLE_ENDIAN);
      WeatherApiResponse mApiResponse = WeatherApiResponse
          .getRootAsWeatherApiResponse(buffer.position(4));
      VariablesWithTime current = mApiResponse.current();

      VariableWithValues temperature2m = new VariablesSearch(current)
          .variable(Variable.temperature).altitude(2).first();
      if (temperature2m == null) {
        return Float.NaN;
      }
      return temperature2m.value();
    }
    catch (IOException | InterruptedException e) {
      throw new RuntimeException(e);
    }

  }

  private void temperatureDemo() {
    String prompt = "What are the current temperatures in Lisbon, Portugal and Reykjavik, Iceland?";

    FunctionToolCallback<Location, Float> callback = FunctionToolCallback
        .builder("fetchTemperature", (Location location) -> fetchTemperature(location))
        .description("Get the current temperature of a location")
        .inputType(Location.class).build();

    var response = this.chatClient.prompt().user(prompt).toolCallbacks(callback).call()
        .chatResponse();
    Generation generation = response.getResult();
    if (generation != null) {
      System.out.println(generation.getOutput().getText());
    }
    else {
      System.out.println("No generation");
    }
  }

}
