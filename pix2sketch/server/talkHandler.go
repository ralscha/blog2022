package main

import (
	"bytes"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/h2non/filetype"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"net/http"
	"os"
)

func (app *application) speechToText(w http.ResponseWriter, r *http.Request) {
	audioData := new(bytes.Buffer)
	_, err := audioData.ReadFrom(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	kind, _ := filetype.Match(audioData.Bytes())
	supportedByWhisper := false
	supportedAudioTypes := []string{"audio/mpeg", "audio/mp4", "audio/x-wav", "video/webm"}

	for _, supported := range supportedAudioTypes {
		if supported == kind.MIME.Value {
			supportedByWhisper = true
			break
		}
	}

	audioDataInput := audioData.Bytes()

	if !supportedByWhisper {
		tmpInput, err := os.CreateTemp("", "audioData")
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer os.Remove(tmpInput.Name())

		_, err = tmpInput.Write(audioData.Bytes())
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		tmpOutput, err := os.CreateTemp("", "audioDataOutput")
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		defer os.Remove(tmpOutput.Name())
		err = ffmpeg.Input(tmpInput.Name()).Output(tmpOutput.Name(), ffmpeg.KwArgs{"f": "mp3", "ab": "192000", "ar": "44100"}).
			OverWriteOutput().ErrorToStdOut().Run()

		if err != nil {
			app.serverError(w, r, err)
			return
		}

		audioDataInput, err = os.ReadFile(tmpOutput.Name())
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	// send to whisper to convert speech to text
	resp, err := app.azureOpenAIClient.GetAudioTranscription(r.Context(),
		azopenai.AudioTranscriptionOptions{
			File:           audioDataInput,
			ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatText),
			DeploymentName: to.Ptr("whisper"),
			Temperature:    to.Ptr(float32(0.0)),
		}, nil)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(*resp.Text))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) chatWithGPT4(w http.ResponseWriter, r *http.Request) {
	textData := new(bytes.Buffer)
	_, err := textData.ReadFrom(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	prompt := textData.String()
	fmt.Println(prompt)

	// send to gpt-4-turbo
	resp, err := app.azureOpenAIClient.GetChatCompletions(r.Context(), azopenai.ChatCompletionsOptions{
		Messages: []azopenai.ChatRequestMessageClassification{
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(prompt)},
		},
		MaxTokens:      to.Ptr(int32(4096)),
		DeploymentName: to.Ptr("gpt-turbo-2024-04-09"),
		N:              to.Ptr(int32(1)),
		Temperature:    to.Ptr(float32(0.0)),
	}, nil)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	gptResponse := *resp.Choices[0].Message.Content

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(gptResponse))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *application) textToSpeech(w http.ResponseWriter, r *http.Request) {
	textData := new(bytes.Buffer)
	_, err := textData.ReadFrom(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	input := textData.String()
	fmt.Println(input)

	// send to tts-hd to convert text to speech
	resp, err := app.azureOpenAIClient.GenerateSpeechFromText(r.Context(),
		azopenai.SpeechGenerationOptions{
			Input:          to.Ptr(input),
			Voice:          to.Ptr(azopenai.SpeechVoiceNova),
			DeploymentName: to.Ptr("tts-hd"),
			ResponseFormat: to.Ptr(azopenai.SpeechGenerationResponseFormatMp3),
		}, nil)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "audio/mp3")
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

}
