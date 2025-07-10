import {Component, inject, NgZone} from '@angular/core';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonIcon,
  IonProgressBar,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';

import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {addIcons} from "ionicons";
import {micCircleOutline, playCircleOutline, stopCircleOutline} from "ionicons/icons";


@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrl: './home.page.scss',
  imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonButton, IonIcon, IonProgressBar]
})
export class HomePage {
  recording = false;
  processing = false;
  speechToTextResponse: string | null = null;
  chatGPT4Response: string | null = null;
  textToSpeechResponse: Blob | null = null;
  error: string | null = null;
  mediaRecorder: MediaRecorder | null = null;
  chunks: Blob[] = [];
  private readonly httpClient = inject(HttpClient);
  private readonly zone = inject(NgZone);

  constructor() {
    addIcons({micCircleOutline, stopCircleOutline, playCircleOutline});
  }

  speechToTextRequest(blob: Blob) {
    this.processing = true;
    this.httpClient.post(`${environment.SERVER_URL}/talk/speechToText`, blob, {
      responseType: 'text'
    })
      .subscribe({
        next: response => {
          this.speechToTextResponse = response;
          this.chatWithGPT4(response);
        },
        error: error => {
          this.error = error.message;
          this.processing = false;
        }
      });
  }

  chatWithGPT4(prompt: string) {
    this.httpClient.post(`${environment.SERVER_URL}/talk/chatWithGPT4`, prompt, {
      responseType: 'text'
    })
      .subscribe({
        next: response => {
          this.chatGPT4Response = response;
          this.textToSpeechRequest(response);
        },
        error: error => {
          this.error = error.message;
          this.processing = false;
        }
      });
  }

  textToSpeechRequest(text: string) {
    this.httpClient.post(`${environment.SERVER_URL}/talk/textToSpeech`, text, {
      responseType: "blob"
    })
      .subscribe({
        next: response => {
          this.textToSpeechResponse = response;
          this.processing = false;
        },
        error: error => {
          this.error = error.message;
          this.processing = false;
        }
      });
  }

  speak() {
    this.playAudio();
  }

  playAudio() {
    if (this.textToSpeechResponse) {
      const audio = new Audio(URL.createObjectURL(this.textToSpeechResponse));
      audio.play();
    }
  }

  startRecording() {
    this.speechToTextResponse = null;
    this.chatGPT4Response = null;
    this.textToSpeechResponse = null;
    this.mediaRecorder = null;
    this.error = null;
    this.chunks = [];
    if (navigator.mediaDevices?.getUserMedia) {
      navigator.mediaDevices.getUserMedia({audio: true})
        .then(stream => {
          this.recording = true;
          this.mediaRecorder = new MediaRecorder(stream);
          this.mediaRecorder.start();
          this.mediaRecorder.ondataavailable = (e) => {
            this.chunks.push(e.data);
          };
        })
        .catch(err => this.error = err);
    } else {
      this.error = "getUserMedia not supported on your browser!";
    }
  }

  stopRecording() {
    this.recording = false;
    if (this.mediaRecorder !== null) {
      this.mediaRecorder.onstop = () => {
        const blob = new Blob(this.chunks);
        this.chunks = [];
        this.zone.run(() => this.speechToTextRequest(blob));
      }
      this.mediaRecorder.stop();
    }
  }


}
