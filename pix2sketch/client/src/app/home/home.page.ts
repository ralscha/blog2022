import {Component} from '@angular/core';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonIcon,
  IonProgressBar,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';

import { HttpClient } from "@angular/common/http";
import {environment} from "../../environments/environment";
import {finalize} from "rxjs";

interface SketchResponse {
  description: string;
  imageBase64: string;
}

@Component({
    selector: 'app-home',
    templateUrl: 'home.page.html',
    styleUrls: ['home.page.scss'],
    imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonButton, IonIcon, IonProgressBar]
})
export class HomePage {
  image: ArrayBuffer | undefined;
  imageData: string | undefined;
  error: string | undefined;
  sketchResponse: SketchResponse | undefined;
  processing = false;

  constructor(private readonly httpClient: HttpClient) {
  }

  selectImage() {
    this.error = undefined;
    this.image = undefined;
    this.sketchResponse = undefined;

    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/*';
    input.onchange = () => {
      const file = input.files![0];
      this.imageData = URL.createObjectURL(file);
      const reader = new FileReader();
      reader.onload = () => {
        if (reader.result instanceof ArrayBuffer) {
          this.image = reader.result;
          this.postRequest();
        }
      };
      reader.readAsArrayBuffer(file);
    };
    input.click();
  }

  postRequest() {
    this.processing = true;
    this.httpClient.post<SketchResponse>(`${environment.SERVER_URL}/sketch`, this.image)
      .pipe(finalize(() => this.processing = false))
      .subscribe({
        next: response => this.sketchResponse = response,
        error: error => {
          console.log(error);
          this.error = error.message;
        }
      });
  }
}
