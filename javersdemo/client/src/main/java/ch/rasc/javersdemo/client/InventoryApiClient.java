package ch.rasc.javersdemo.client;

import java.math.BigDecimal;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.HashMap;
import java.util.Map;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

public class InventoryApiClient {

  private final HttpClient httpClient;
  private final String baseUrl;
  private final ObjectMapper objectMapper;
  private String authToken;

  public InventoryApiClient(String baseUrl) {
    this.baseUrl = baseUrl.endsWith("/") ? baseUrl.substring(0, baseUrl.length() - 1)
        : baseUrl;
    this.httpClient = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(10))
        .build();
    this.objectMapper = new ObjectMapper();
  }

  public boolean authenticate(String username, String password) {
    try {
      Map<String, String> authRequest = new HashMap<>();
      authRequest.put("username", username);
      authRequest.put("password", password);

      String jsonBody = this.objectMapper.writeValueAsString(authRequest);

      HttpRequest request = HttpRequest.newBuilder()
          .uri(new URI(this.baseUrl + "/api/auth/authenticate"))
          .header("Content-Type", "application/json")
          .POST(HttpRequest.BodyPublishers.ofString(jsonBody)).build();

      HttpResponse<String> response = this.httpClient.send(request,
          HttpResponse.BodyHandlers.ofString());

      if (response.statusCode() == 200) {
        Map<String, Object> responseBody = this.objectMapper.readValue(response.body(),
            Map.class);
        this.authToken = (String) responseBody.get("token");
        return true;
      }
      System.err.println("Authentication failed: " + response.body());
      return false;

    }
    catch (Exception e) {
      System.err.println("Authentication failed: " + e.getMessage());
      return false;
    }
  }

  public String getAllItems() {
    return sendGetRequest("/api/inventory/items");
  }

  public String getItemById(Long id) {
    return sendGetRequest("/api/inventory/items/" + id);
  }

  public String createItem(Map<String, Object> item) {
    return sendPostRequest("/api/inventory/items", item);
  }

  public String updateItem(Long id, Map<String, Object> item) {
    return sendPutRequest("/api/inventory/items/" + id, item);
  }

  public String updateQuantity(Long id, Integer quantity) {
    Map<String, Object> quantityUpdate = new HashMap<>();
    quantityUpdate.put("quantity", quantity);
    return sendPatchRequest("/api/inventory/items/" + id + "/quantity", quantityUpdate);
  }

  public String deleteItem(Long id) {
    return sendDeleteRequest("/api/inventory/items/" + id);
  }

  public String logout() {
    this.authToken = null;
    return "Status: 200\nResponse: {\"message\":\"Logged out successfully\"}";
  }

  public String getInventoryItemSnapshots(Long id) {
    return getInventoryItemSnapshots(id, 10);
  }

  public String getInventoryItemSnapshots(Long id, int limit) {
    return sendGetRequest(
        "/api/audit/inventory-items/" + id + "/snapshots?limit=" + limit);
  }

  public String getInventoryItemChanges(Long id) {
    return sendGetRequest("/api/audit/inventory-items/" + id + "/changes");
  }

  public String getInventoryItemShadows(Long id) {
    return sendGetRequest("/api/audit/inventory-items/" + id + "/shadows");
  }

  private String sendGetRequest(String endpoint) {
    try {
      HttpRequest.Builder requestBuilder = HttpRequest.newBuilder()
          .uri(new URI(this.baseUrl + endpoint)).GET();

      if (this.authToken != null) {
        requestBuilder.header("Authorization", "Bearer " + this.authToken);
      }

      HttpRequest request = requestBuilder.build();
      HttpResponse<String> response = this.httpClient.send(request,
          HttpResponse.BodyHandlers.ofString());

      return formatResponse(response);

    }
    catch (Exception e) {
      return "Error: " + e.getMessage();
    }
  }

  private String sendPostRequest(String endpoint, Map<String, Object> body) {
    return sendJsonRequest("POST", endpoint, body);
  }

  private String sendPutRequest(String endpoint, Map<String, Object> body) {
    return sendJsonRequest("PUT", endpoint, body);
  }

  private String sendPatchRequest(String endpoint, Map<String, Object> body) {
    return sendJsonRequest("PATCH", endpoint, body);
  }

  private String sendJsonRequest(String method, String endpoint,
      Map<String, Object> body) {
    try {
      String jsonBody = this.objectMapper.writeValueAsString(body);

      HttpRequest.Builder requestBuilder = HttpRequest.newBuilder()
          .uri(new URI(this.baseUrl + endpoint))
          .header("Content-Type", "application/json")
          .method(method, HttpRequest.BodyPublishers.ofString(jsonBody));

      if (this.authToken != null) {
        requestBuilder.header("Authorization", "Bearer " + this.authToken);
      }

      HttpRequest request = requestBuilder.build();
      HttpResponse<String> response = this.httpClient.send(request,
          HttpResponse.BodyHandlers.ofString());

      return formatResponse(response);

    }
    catch (Exception e) {
      return "Error: " + e.getMessage();
    }
  }

  private String sendDeleteRequest(String endpoint) {
    try {
      HttpRequest.Builder requestBuilder = HttpRequest.newBuilder()
          .uri(new URI(this.baseUrl + endpoint)).DELETE();

      if (this.authToken != null) {
        requestBuilder.header("Authorization", "Bearer " + this.authToken);
      }

      HttpRequest request = requestBuilder.build();
      HttpResponse<String> response = this.httpClient.send(request,
          HttpResponse.BodyHandlers.ofString());

      return formatResponse(response);

    }
    catch (Exception e) {
      return "Error: " + e.getMessage();
    }
  }

  private static String formatResponse(HttpResponse<String> response) {
    StringBuilder result = new StringBuilder();
    result.append("Status: ").append(response.statusCode()).append("\n");
    result.append("Response: ").append(response.body());
    return result.toString();
  }

  public static void main(String[] args)
      throws JsonMappingException, JsonProcessingException {
    InventoryApiClient client = new InventoryApiClient("http://localhost:8080");
    boolean authSuccess = client.authenticate("user1", "password");
    if (!authSuccess) {
      System.out.println("Authentication failed for user1");
      return;
    }

    Map<String, Object> newItem = new HashMap<>();
    newItem.put("name", "Laptop Computer");
    newItem.put("description", "High-performance laptop for business use");
    newItem.put("quantity", 50);
    newItem.put("price", new BigDecimal("899.99"));
    newItem.put("category", "Electronics");
    newItem.put("sku", "LAP-001");
    String createResponse = client.createItem(newItem);
    createResponse = createResponse.substring(createResponse.indexOf("{"));
    Map<String, Object> result = client.objectMapper.readValue(createResponse, Map.class);
    Long insertedId = ((Integer) result.get("id")).longValue();

    client.logout();
    authSuccess = client.authenticate("user2", "password");
    if (!authSuccess) {
      System.out.println("Authentication failed for user2");
      return;
    }
    String response = client.updateQuantity(insertedId, 65);
    System.out.println(response);

    client.logout();
    authSuccess = client.authenticate("user1", "password");
    if (!authSuccess) {
      System.out.println("Authentication failed for user1");
      return;
    }
    response = client.deleteItem(insertedId);
    System.out.println(response);

    boolean adminAuthSuccess = client.authenticate("admin", "password");
    if (adminAuthSuccess) {

      System.out.println("CHANGES");
      System.out.println(client.getInventoryItemChanges(insertedId));
      System.out.println();

      System.out.println("SNAPSHOTS");
      System.out.println(client.getInventoryItemSnapshots(insertedId));
      System.out.println();

      System.out.println("SHADOWS");
      System.out.println(client.getInventoryItemShadows(insertedId));
      System.out.println();
    }

    client.logout();
  }
}
