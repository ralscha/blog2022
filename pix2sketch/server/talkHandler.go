package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/h2non/filetype"
	"github.com/openai/openai-go"
	ffmpeg "github.com/u2takey/ffmpeg-go"
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
	resp, err := app.azureOpenAIClient.Audio.Transcriptions.New(r.Context(), openai.AudioTranscriptionNewParams{
		Model:          "whisper",
		File:           bytes.NewReader(audioDataInput),
		ResponseFormat: openai.AudioResponseFormatText,
		Temperature:    openai.Float(0.0),
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(resp.Text))
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
	resp, err := app.azureOpenAIClient.Chat.Completions.New(r.Context(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(prompt),
					},
				},
			},
		},
		MaxTokens:   openai.Int(4096),
		Model:       "gpt-turbo-2024-04-09",
		N:           openai.Int(1),
		Temperature: openai.Float(0.0),
	}, nil)

	if err != nil {
		app.serverError(w, r, err)
		return
	}
	gptResponse := resp.Choices[0].Message.Content

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
	resp, err := app.azureOpenAIClient.Audio.Speech.New(r.Context(), openai.AudioSpeechNewParams{
		Model:          "tts-hd",
		Input:          input,
		Voice:          openai.AudioSpeechNewParamsVoiceAsh,
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatMP3,
	})

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
