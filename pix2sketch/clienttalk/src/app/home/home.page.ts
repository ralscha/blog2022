import { Component, inject, signal } from '@angular/core';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonIcon,
  IonProgressBar,
  IonTitle,
  IonToolbar,
} from '@ionic/angular/standalone';

import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { addIcons } from 'ionicons';
import { micCircleOutline, playCircleOutline, stopCircleOutline } from 'ionicons/icons';

@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrl: './home.page.scss',
  imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonButton, IonIcon, IonProgressBar],
})
export class HomePage {
  recording = signal(false);
  processing = signal(false);
  speechToTextResponse = signal<string | null>(null);
  chatGPT4Response = signal<string | null>(null);
  textToSpeechResponse = signal<Blob | null>(null);
  error = signal<string | null>(null);
  mediaRecorder: MediaRecorder | null = null;
  chunks: Blob[] = [];
  private readonly httpClient = inject(HttpClient);

  constructor() {
    addIcons({ micCircleOutline, stopCircleOutline, playCircleOutline });
  }

  speechToTextRequest(blob: Blob) {
    this.processing.set(true);
    this.httpClient
      .post(`${environment.SERVER_URL}/talk/speechToText`, blob, {
        responseType: 'text',
      })
      .subscribe({
        next: (response) => {
          this.speechToTextResponse.set(response);
          this.chatWithGPT4(response);
        },
        error: (error) => {
          this.error.set(error.message);
          this.processing.set(false);
        },
      });
  }

  chatWithGPT4(prompt: string) {
    this.httpClient
      .post(`${environment.SERVER_URL}/talk/chatWithGPT4`, prompt, {
        responseType: 'text',
      })
      .subscribe({
        next: (response) => {
          this.chatGPT4Response.set(response);
          this.textToSpeechRequest(response);
        },
        error: (error) => {
          this.error.set(error.message);
          this.processing.set(false);
        },
      });
  }

  textToSpeechRequest(text: string) {
    this.httpClient
      .post(`${environment.SERVER_URL}/talk/textToSpeech`, text, {
        responseType: 'blob',
      })
      .subscribe({
        next: (response) => {
          this.textToSpeechResponse.set(response);
          this.processing.set(false);
        },
        error: (error) => {
          this.error.set(error.message);
          this.processing.set(false);
        },
      });
  }

  speak() {
    this.playAudio();
  }

  playAudio() {
    const textToSpeechResponse = this.textToSpeechResponse();
    if (textToSpeechResponse) {
      const audio = new Audio(URL.createObjectURL(textToSpeechResponse));
      audio.play();
    }
  }

  startRecording() {
    this.speechToTextResponse.set(null);
    this.chatGPT4Response.set(null);
    this.textToSpeechResponse.set(null);
    this.mediaRecorder = null;
    this.error.set(null);
    this.chunks = [];
    if (navigator.mediaDevices?.getUserMedia) {
      navigator.mediaDevices
        .getUserMedia({ audio: true })
        .then((stream) => {
          this.recording.set(true);
          this.mediaRecorder = new MediaRecorder(stream);
          this.mediaRecorder.start();
          this.mediaRecorder.ondataavailable = (e) => {
            this.chunks.push(e.data);
          };
        })
        .catch((err) => {
          this.error.set(err);
        });
    } else {
      this.error.set('getUserMedia not supported on your browser!');
    }
  }

  stopRecording() {
    this.recording.set(false);
    if (this.mediaRecorder !== null) {
      this.mediaRecorder.onstop = () => {
        const blob = new Blob(this.chunks);
        this.chunks = [];
        this.speechToTextRequest(blob);
      };
      this.mediaRecorder.stop();
    }
  }
}
