import {Component, inject, OnInit} from '@angular/core';

import {HttpClient} from "@angular/common/http";
import {SwPush} from "@angular/service-worker";
import {firstValueFrom, lastValueFrom} from "rxjs";
import {environment} from "../environments/environment";

@Component({
    selector: 'app-root',
    templateUrl: 'app.component.html',
    styleUrls: ['app.component.css']
})
export class AppComponent implements OnInit {
  subscribed = false;
  webPushSupported = false;
  permissionDenied = false;
  enabled = false;
  readonly #httpClient = inject(HttpClient);
  readonly #swPush = inject(SwPush);
  #serverPublicKey: string | null = null;
  #currentSubscription: PushSubscription | null = null;

  async ngOnInit() {
    // standalone: boolean indicating whether the browser is running in standalone mode.
    // Available on Apple's iOS Safari only
    const isIOS = 'standalone' in window.navigator;
    const isIOSStandalone = 'standalone' in window.navigator && window.navigator.standalone === true;
    this.enabled = this.#swPush.isEnabled;
    if (this.#swPush.isEnabled && (!isIOS || isIOSStandalone)) {
      this.webPushSupported = true;

      // fetch the current subscription
      this.#currentSubscription = await firstValueFrom(this.#swPush.subscription)
      if (this.#currentSubscription) {
        this.subscribed = true;
        await lastValueFrom(this.#httpClient.post(`${environment.SERVER_URL}/subscribe`, this.#currentSubscription));
      }
    } else {
      this.webPushSupported = false;
    }
  }

  async subscribe() {
    if (!this.#currentSubscription) {
      this.#serverPublicKey = await lastValueFrom(this.#httpClient.get(`${environment.SERVER_URL}/publicKey`, {responseType: 'text'}))

      try {
        this.#currentSubscription = await this.#swPush.requestSubscription({
          serverPublicKey: this.#serverPublicKey!
        });
      } catch (e) {
        console.error(e);
        this.permissionDenied = true;
        return;
      }
    }

    if (this.#currentSubscription) {
      await lastValueFrom(this.#httpClient.post(`${environment.SERVER_URL}/subscribe`, this.#currentSubscription));
      this.subscribed = true;
    } else {
      this.subscribed = false;
    }
  }

  async unsubscribe() {
    if (this.#currentSubscription) {
      await this.#currentSubscription.unsubscribe();
      await lastValueFrom(this.#httpClient.post(`${environment.SERVER_URL}/unsubscribe`, this.#currentSubscription));
      this.#currentSubscription = null;
      this.subscribed = false;
    }
  }

}
