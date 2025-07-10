import {Component, inject} from '@angular/core';
import {MessagesService} from '../messages.service';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';

import {Client} from '@passwordlessdev/passwordless-client';
import {SigninInput} from "../api/types";
import {RouterLink} from '@angular/router';
import {
  IonButton,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonRouterLink,
  IonRow,
  IonTitle,
  IonToolbar,
  NavController
} from "@ionic/angular/standalone";

@Component({
  selector: 'app-authentication',
  templateUrl: './authentication.page.html',
  imports: [RouterLink, IonRouterLink, IonHeader, IonToolbar, IonContent, IonGrid, IonCol, IonButton, IonRow, IonTitle]
})
export class AuthenticationPage {
  readonly #navCtrl = inject(NavController);
  readonly #httpClient = inject(HttpClient);
  readonly #messagesService = inject(MessagesService);

  async login(): Promise<void> {
    const passwordlessClient = new Client({
      apiKey: environment.PASSWORDLESS_PUBLIC_KEY,
    });

    const {token, error} = await passwordlessClient.signinWithDiscoverable();

    if (error) {
      await this.#messagesService.showErrorToast('Login failed');
      return;
    }

    const signinInput: SigninInput = {
      token,
    };
    this.#httpClient.post<void>(`${environment.API_URL}/signin`, signinInput).subscribe({
      next: () => {
        this.#navCtrl.navigateRoot('/home', {replaceUrl: true});
      },
      error: () => {
        this.#messagesService.showErrorToast('Login failed');
      }
    });

  }

}
