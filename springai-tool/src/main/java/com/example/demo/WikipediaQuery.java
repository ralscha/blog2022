package com.example.demo;

import org.springframework.ai.tool.annotation.ToolParam;

public record WikipediaQuery(
    @ToolParam(description = "The search query") String searchQuery) {
}
