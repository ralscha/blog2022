package com.example.demo;

import java.util.Arrays;
import java.util.List;
import java.util.Map;

import org.springframework.ai.chat.client.ChatClient;
import org.springframework.ai.chat.client.advisor.api.Advisor;
import org.springframework.ai.chat.client.advisor.vectorstore.QuestionAnswerAdvisor;
import org.springframework.ai.document.Document;
import org.springframework.ai.embedding.EmbeddingModel;
import org.springframework.ai.rag.advisor.RetrievalAugmentationAdvisor;
import org.springframework.ai.rag.preretrieval.query.expansion.MultiQueryExpander;
import org.springframework.ai.rag.preretrieval.query.transformation.RewriteQueryTransformer;
import org.springframework.ai.rag.retrieval.search.VectorStoreDocumentRetriever;
import org.springframework.ai.reader.ExtractedTextFormatter;
import org.springframework.ai.reader.pdf.PagePdfDocumentReader;
import org.springframework.ai.reader.pdf.config.PdfDocumentReaderConfig;
import org.springframework.ai.transformer.splitter.TokenTextSplitter;
import org.springframework.ai.vectorstore.SearchRequest;
import org.springframework.ai.vectorstore.VectorStore;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.core.io.Resource;
import org.springframework.web.client.RestTemplate;

@SpringBootApplication
public class DemoApplication implements CommandLineRunner {

  public static void main(String[] args) {
    SpringApplication.run(DemoApplication.class, args);
  }

  private final RestTemplate restTemplate = new RestTemplate();
  private final TokenTextSplitter splitter = new TokenTextSplitter();

  private final ChatClient chatClient;
  private final ChatClient.Builder chatClientBuilder;

  public DemoApplication(ChatClient.Builder chatClientBuilder) {
    this.chatClientBuilder = chatClientBuilder;
    this.chatClient = chatClientBuilder.build();
  }

  @Autowired
  private EmbeddingModel embeddingModel;

  private void embeddingDemo() {
    List<float[]> embeddings = this.embeddingModel.embed(List.of("Hello", "World"));
    for (float[] embedding : embeddings) {
      System.out.println(Arrays.toString(embedding));
      System.out.println("Embedding size: " + embedding.length);
    }
  }

  @Autowired
  VectorStore vectorStore;

  @Override
  public void run(String... args) throws Exception {
    // simple();
    // embeddingDemo();
    ingest();
    rag1();
    rag2();
    ragAdvanced1();
    ragAdvanced2();
  }

  private void simple() {
    List<Document> documents = List.of(new Document(
        "The autumn leaves dance silently as they paint the ground in shades of amber and gold.",
        Map.of("id", "1")),
        new Document(
            "Quantum computers manipulate reality at the edge of what's possible in our universe.",
            Map.of("id", "2")),
        new Document(
            "Through the telescope, distant galaxies tell stories of time's infinite passage.",
            Map.of("id", "3")),
        new Document(
            "Children's laughter echoes across the playground like wind chimes in a gentle breeze.",
            Map.of("id", "4")),
        new Document(
            "Deep in the ocean's trenches, bioluminescent creatures create their own constellations.",
            Map.of("id", "5")));

    this.vectorStore.add(documents);

    List<Document> foundDocuments = this.vectorStore.similaritySearch(
        SearchRequest.builder().query("Tell me more about quantum computers").topK(3)
            .similarityThreshold(0.6).build());

    for (Document document : foundDocuments) {
      System.out.println(document.getMetadata().get("id"));
    }
  }

  private void rag1() {
    String response = this.chatClient.prompt()
        .advisors(QuestionAnswerAdvisor.builder(this.vectorStore).build())
        .user("What are the main causes of the climate change?").call().content();
    System.out.println(response);
    System.out.println("===========================");
  }

  private void rag2() {
    var qaAdvisor = QuestionAnswerAdvisor.builder(this.vectorStore)
        .searchRequest(SearchRequest.builder().similarityThreshold(0.6).topK(20).build())
        .build();
    String response = this.chatClient.prompt().advisors(qaAdvisor)
        .user("What are the main causes of the climate change?").call().content();
    System.out.println(response);
    System.out.println("===========================");
  }

  private void ragAdvanced1() {
    Advisor retrievalAugmentationAdvisor = RetrievalAugmentationAdvisor.builder()
        .queryTransformers(RewriteQueryTransformer.builder()
            .chatClientBuilder(this.chatClientBuilder.build().mutate()).build())
        .documentRetriever(VectorStoreDocumentRetriever.builder().similarityThreshold(0.5)
            .topK(20).vectorStore(this.vectorStore).build())
        .build();

    String response = this.chatClient.prompt().advisors(retrievalAugmentationAdvisor)
        .user("What are the main causes of the climate change?").call().content();
    System.out.println(response);
    System.out.println("===========================");
  }

  private void ragAdvanced2() {
    Advisor retrievalAugmentationAdvisor = RetrievalAugmentationAdvisor.builder()
        .queryExpander(MultiQueryExpander.builder()
            .chatClientBuilder(this.chatClientBuilder.build().mutate())
            .includeOriginal(true).numberOfQueries(3).build())
        .documentRetriever(VectorStoreDocumentRetriever.builder().similarityThreshold(0.5)
            .topK(20).vectorStore(this.vectorStore).build())
        .build();

    String response = this.chatClient.prompt().advisors(retrievalAugmentationAdvisor)
        .user("What are the main causes of the climate change?").call().content();
    System.out.println(response);
    System.out.println("===========================");
  }

  private final static List<String> urls = List.of(
      "https://www.ipcc.ch/site/assets/uploads/2025/01/2407_CDR_CCUS_Report.pdf",
      "https://www.ipcc.ch/site/assets/uploads/2023/07/IPCC_2023_Workshop_Report_Scenarios.pdf",
      "https://www.ipcc.ch/site/assets/uploads/2018/05/EMR_TGICA_Future.pdf",
      "https://www.ipcc.ch/site/assets/uploads/2018/02/IPCC_2017_EMR_Scenarios.pdf",
      "https://www.ipcc.ch/report/ar6/wg3/downloads/report/IPCC_AR6_WGIII_FullReport.pdf",
      "https://www.ipcc.ch/report/ar6/wg2/downloads/report/IPCC_AR6_WGII_FullReport.pdf",
      "https://www.ipcc.ch/report/ar6/syr/downloads/report/IPCC_AR6_SYR_FullVolume.pdf",
      "https://www.ipcc.ch/report/ar6/wg1/downloads/report/IPCC_AR6_WGI_FullReport.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SPM_version_report_LR.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SR15_Chapter_1_HR.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SR15_Chapter_2_LR.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SR15_Chapter_3_LR.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SR15_Chapter_4_LR.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SR15_Chapter_5_LR.pdf",
      "https://www.ipcc.ch/site/assets/uploads/sites/2/2022/06/SR15_AnnexI.pdf");

  private void ingest() {
    for (String url : urls) {
      byte[] pdfBytes = this.restTemplate.getForObject(url, byte[].class);
      List<Document> documents = loadPdfs(new ByteArrayResource(pdfBytes) {
        @Override
        public String getFilename() {
          return url;
        }
      });
      this.vectorStore.add(documents);
    }
  }

  List<Document> loadPdfs(Resource resourcePdf) {

    PagePdfDocumentReader pdfReader = new PagePdfDocumentReader(resourcePdf,
        PdfDocumentReaderConfig.builder().withPageTopMargin(0).withPageBottomMargin(0)
            .withPageExtractedTextFormatter(
                ExtractedTextFormatter.builder().withNumberOfTopTextLinesToDelete(0)
                    .withNumberOfBottomTextLinesToDelete(0).build())
            .withPagesPerDocument(1).build());

    return this.splitter.apply(pdfReader.read());
  }

}
