import { ChangeDetectionStrategy, ChangeDetectorRef, Component, inject } from '@angular/core';
import {
  IonButton,
  IonContent,
  IonHeader,
  IonIcon,
  IonProgressBar,
  IonTitle,
  IonToolbar,
} from '@ionic/angular';

import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { finalize } from 'rxjs';
import { addIcons } from 'ionicons';
import { cameraOutline } from 'ionicons/icons';

interface SketchResponse {
  description: string;
  imageBase64: string;
}

@Component({
  changeDetection: ChangeDetectionStrategy.Eager,
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrl: './home.page.scss',
  imports: [IonHeader, IonToolbar, IonTitle, IonContent, IonButton, IonIcon, IonProgressBar],
})
export class HomePage {
  image: ArrayBuffer | undefined;
  imageData: string | undefined;
  error: string | undefined;
  sketchResponse: SketchResponse | undefined;
  processing = false;
  private readonly httpClient = inject(HttpClient);
  private readonly changeDetectorRef = inject(ChangeDetectorRef);

  constructor() {
    addIcons({ cameraOutline });
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
      this.changeDetectorRef.markForCheck();
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
    this.changeDetectorRef.markForCheck();
    this.httpClient
      .post<SketchResponse>(`${environment.SERVER_URL}/sketch`, this.image)
      .pipe(
        finalize(() => {
          this.processing = false;
          this.changeDetectorRef.markForCheck();
        }),
      )
      .subscribe({
        next: (response) => {
          this.sketchResponse = response;
          this.changeDetectorRef.markForCheck();
        },
        error: (error) => {
          console.log(error);
          this.error = error.message;
          this.changeDetectorRef.markForCheck();
        },
      });
  }
}
