import { Component, inject, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { IonApp, IonRouterOutlet } from '@ionic/angular/standalone';
import { PocketbaseService } from './services/pocketbase.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
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
