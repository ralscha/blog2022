package com.example.demo;

import tools.jackson.annotation.JsonClassDescription;
import tools.jackson.annotation.JsonPropertyDescription;

@JsonClassDescription("A query to search Wikipedia")
public record WikipediaQuery(
    @JsonPropertyDescription("The search query") String searchQuery) {
}