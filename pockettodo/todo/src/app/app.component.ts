import {
  Component,
  inject,
  OnInit,
  ChangeDetectionStrategy
} from '@angular/core';
import { Router } from '@angular/router';
import { IonApp, IonRouterOutlet } from '@ionic/angular/standalone';
import { PocketbaseService } from './services/pocketbase.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [IonApp, IonRouterOutlet]
})
export class AppComponent implements OnInit {
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);

  ngOnInit() {
    this.checkAuthStatus();
  }

  private async checkAuthStatus() {
    await this.pocketbaseService.refreshAuth();

    if (this.pocketbaseService.isLoggedIn()) {
      this.router.navigate(['/todos']);
    } else {
      this.router.navigate(['/login']);
    }
  }
}
