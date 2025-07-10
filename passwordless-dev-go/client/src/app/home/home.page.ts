import {Component, inject, OnInit} from '@angular/core';
import {AuthService} from '../auth.service';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {SecretOutput} from "../api/types";
import {IonButton, IonContent, IonHeader, IonTitle, IonToolbar, NavController} from "@ionic/angular/standalone";

@Component({
  selector: 'app-home',
  templateUrl: './home.page.html',
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonButton
  ]
})
export class HomePage implements OnInit {
  secretMessage: string | null = null;

  readonly #navCtrl = inject(NavController);
  readonly #httpClient = inject(HttpClient);
  readonly #authService = inject(AuthService);

  async logout(): Promise<void> {
    this.#authService.logout().subscribe(() => this.#navCtrl.navigateRoot('/login'));
  }

  ngOnInit(): void {
    this.#httpClient.get<SecretOutput>(`${environment.API_URL}/secret`)
      .subscribe(response => this.secretMessage = response.message);
  }

}
