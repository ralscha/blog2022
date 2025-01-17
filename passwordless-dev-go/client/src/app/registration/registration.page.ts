import {Component, inject} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {MessagesService} from '../messages.service';
import {environment} from '../../environments/environment';
import {Router} from "@angular/router";
// @ts-expect-error - This package does not have typescript definitions
import {Client} from '@passwordlessdev/passwordless-client';
import {CreateTokenInput, CreateTokenOutput} from "../api/types";
import {FormsModule} from '@angular/forms';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonInput,
  IonItem,
  IonRow,
  IonTitle,
  IonToolbar
} from "@ionic/angular/standalone";

@Component({
  selector: 'app-registration',
  templateUrl: './registration.page.html',
  imports: [FormsModule, IonHeader, IonToolbar, IonButtons, IonBackButton, IonTitle, IonContent, IonGrid, IonRow, IonCol, IonItem, IonInput, IonButton]
})
export class RegistrationPage {

  readonly #router = inject(Router);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);

  async register(username: string | null): Promise<void> {
    if (!username) {
      return;
    }

    const passwordlessClient = new Client({
      apiKey: environment.PASSWORDLESS_PUBLIC_KEY,
    });

    const createTokenInput: CreateTokenInput = {username};

    this.#httpClient.post<CreateTokenOutput>(`${environment.API_URL}/create-token`, createTokenInput).subscribe({
      next: async (response) => {
        const {error} = await passwordlessClient.register(response.token);
        if (error) {
          await this.#messagesService.showErrorToast('Registration failed');
          return;
        }
        await this.#router.navigate(['/login']);
      }
    });

  }


}

