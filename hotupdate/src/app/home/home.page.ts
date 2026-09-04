import { ChangeDetectionStrategy, Component } from '@angular/core';
import { IonContent, IonHeader, IonTitle, IonToolbar } from '@ionic/angular';
import { environment } from '../../environments/environment';

@Component({
  changeDetection: ChangeDetectionStrategy.Eager,
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss'],
  imports: [IonHeader, IonToolbar, IonTitle, IonContent],
})
export class HomePage {
  version = '';

  constructor() {
    this.version = environment.version;
  }
}
