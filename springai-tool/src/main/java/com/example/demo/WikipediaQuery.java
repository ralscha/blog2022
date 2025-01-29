package com.example.demo;

import com.fasterxml.jackson.annotation.JsonClassDescription;
import com.fasterxml.jackson.annotation.JsonPropertyDescription;

@JsonClassDescription("A query to search Wikipedia")
public record WikipediaQuery(
    @JsonPropertyDescription("The search query") String searchQuery) {
}