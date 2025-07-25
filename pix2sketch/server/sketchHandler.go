package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/openai/openai-go"
	"net/http"
)

type SketchResponse struct {
	Description string `json:"description"`
	ImageBase64 string `json:"imageBase64"`
}

func (app *application) sketch(w http.ResponseWriter, r *http.Request) {
	imageData := new(bytes.Buffer)
	_, err := imageData.ReadFrom(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// send request to gpt-4o on azure to describe the image
	base64Image := base64.StdEncoding.EncodeToString(imageData.Bytes())
	imageUrl := "data:image/png;base64," + base64Image

	resp, err := app.azureOpenAIClient.Chat.Completions.New(r.Context(), openai.ChatCompletionNewParams{
		Model: "gpt-4o",
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfArrayOfContentParts: []openai.ChatCompletionContentPartUnionParam{
							{
								OfText: &openai.ChatCompletionContentPartTextParam{
									Text: "What's in this image? Return a detailed description:",
								},
							},
							{
								OfImageURL: &openai.ChatCompletionContentPartImageParam{
									ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
										URL: imageUrl,
									},
								},
							},
						},
					},
				},
			},
		},
		MaxTokens:   openai.Int(4096),
		N:           openai.Int(1),
		Temperature: openai.Float(0.0),
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	imageDescription := resp.Choices[0].Message.Content

	// send request to stability on aws bedrock to generate a sketch based on the description

	awsRegion := "us-west-2"

	// only when application runs outside AWS
	staticCredentials := aws.NewCredentialsCache(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(app.config.AwsBedrockUserAccessKey, app.config.AwsBedrockUserSecretKey, "")))
	sdkConfig, err := config.LoadDefaultConfig(r.Context(), config.WithRegion(awsRegion),
		config.WithCredentialsProvider(staticCredentials))

	// use this config if application runs in AWS. For example: Lambda, EC2, ECS, ...
	// sdkConfig, err := config.LoadDefaultConfig(r.Context(), config.WithRegion(awsRegion))

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	client := bedrockruntime.NewFromConfig(sdkConfig)
	modelId := "stability.sd3-large-v1:0"

	requestImageDescription := "Create a high quality drawing based on the following description:\n\n" + imageDescription

	nativeRequest := map[string]any{
		"prompt":        requestImageDescription,
		"mode":          "text-to-image",
		"aspect_ratio":  "1:1",
		"output_format": "jpeg",
	}

	body, err := json.Marshal(nativeRequest)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	response, err := client.InvokeModel(r.Context(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	})

	if response == nil {
		app.serverError(w, r, err)
		return
	}

	modelResponse := map[string]any{}
	err = json.Unmarshal(response.Body, &modelResponse)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	images := modelResponse["images"].([]any)
	if len(images) == 0 {
		app.serverError(w, r, err)
		return
	}

	imageBase64 := images[0].(string)

	sketchResponse := SketchResponse{
		Description: imageDescription,
		ImageBase64: imageBase64,
	}

	err = JSON(w, http.StatusOK, sketchResponse)
	if err != nil {
		app.serverError(w, r, err)
	}

}
