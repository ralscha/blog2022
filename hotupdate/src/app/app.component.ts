import {Component} from '@angular/core';
import {IonApp, IonRouterOutlet} from '@ionic/angular/standalone';
import {CapacitorUpdater} from '@capgo/capacitor-updater'
import {HttpClient} from "@angular/common/http";
import {firstValueFrom} from "rxjs";
import {environment} from "../environments/environment";
import {Dialog} from "@capacitor/dialog";
import {App, AppState} from "@capacitor/app";

type UpdateInfo = {
  version: string;
  downloadURL: string;
}

@Component({
  selector: 'app-root',
  templateUrl: 'app.component.html',
  standalone: true,
  imports: [IonApp, IonRouterOutlet],
})
export class AppComponent {
  constructor(private readonly httpClient: HttpClient) {
    CapacitorUpdater.notifyAppReady()

    App.addListener('appStateChange', (state: AppState) => {
      if (state.isActive) {
        this.#checkForUpgrade();
      }
    });
    this.#checkForUpgrade();
  }

  #checkForUpgrade() {
    this.#readUpdateInfo().then(updateInfo => {
      if (updateInfo.version !== environment.version) {
        this.#upgrade(updateInfo);
      }
    });
  }

  #readUpdateInfo(): Promise<UpdateInfo> {
    return firstValueFrom(this.httpClient.get<UpdateInfo>('https://static.rasc.ch/update.json'));
  }

  async #upgrade(updateInfo: UpdateInfo): Promise<void> {
    const version = await CapacitorUpdater.download({
      url: updateInfo.downloadURL,
      version: updateInfo.version,
    })

    const {value: okButtonClicked} = await Dialog.confirm({
      title: 'New Version Available',
      message: `Do you want to upgrade to version ${updateInfo.version}?`,
      okButtonTitle: 'Upgrade',
      cancelButtonTitle: 'Later',
    });
    if (okButtonClicked) {
      await CapacitorUpdater.set(version);
    }

  }
}
