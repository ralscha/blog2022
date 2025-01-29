package com.example.demo;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.List;

import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.nodes.Element;
import org.springframework.ai.tool.annotation.Tool;
import org.springframework.stereotype.Service;

import com.example.demo.WikipediaArticleFetcher.SearchResponse.SearchResult;
import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.databind.ObjectMapper;

@Service
public class WikipediaArticleFetcher {

  private final ObjectMapper objectMapper = new ObjectMapper();

  @JsonIgnoreProperties(ignoreUnknown = true)
  record SearchResponse(Query query) {
    @JsonIgnoreProperties(ignoreUnknown = true)
    record Query(List<SearchResult> search) {
    }

    @JsonIgnoreProperties(ignoreUnknown = true)
    record SearchResult(String title) {
    }
  }

  @Tool(description = "Searches for a Wikipedia article and returns the text")
  public String search(WikipediaQuery query) {
    System.out.println(
        "Calling WikipediaArticleFetcher with parameters: " + query.searchQuery());
    String searchUrl = "https://en.wikipedia.org/w/api.php?action=query&list=search&srsearch="
        + query.searchQuery().replace(" ", "%20") + "&format=json";

    try (HttpClient client = HttpClient.newHttpClient()) {
      HttpRequest request = HttpRequest.newBuilder().uri(URI.create(searchUrl)).build();
      var response = client.send(request, HttpResponse.BodyHandlers.ofString());
      var responseBody = response.body();
      SearchResponse searchResponse = this.objectMapper.readValue(responseBody,
          SearchResponse.class);

      if (searchResponse.query().search().isEmpty()) {
        return "";
      }

      List<SearchResult> searchResponses = searchResponse.query().search().subList(0,
          Math.min(searchResponse.query().search().size(), 3));
      StringBuilder context = new StringBuilder();
      for (SearchResult sr : searchResponses) {
        String url = "https://en.wikipedia.org/wiki/" + sr.title().replace(" ", "_");
        Document doc = Jsoup.connect(url).get();
        Element mainElement = doc.select("div[id=mw-content-text]").first();
        String text = mainElement.text();
        context.append(text.replaceAll("\\[.*?\\]+", ""));
        context.append("\n\n");
      }
      return context.toString();
    }
    catch (Exception e) {
      throw new RuntimeException(e);
    }
  }

}
