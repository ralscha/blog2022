import {
  Component,
  inject,
  OnInit,
  ChangeDetectionStrategy
} from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import {
  IonButton,
  IonCard,
  IonCardContent,
  IonCardHeader,
  IonCardTitle,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonRouterLink,
  IonRow,
  IonText,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';
import { PocketbaseService } from '../services/pocketbase.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.page.html',
  styleUrl: './home.page.css',
  changeDetection: ChangeDetectionStrategy.Eager,
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonButton,
    IonCard,
    IonCardHeader,
    IonCardTitle,
    IonCardContent,
    IonText,
    IonGrid,
    IonRow,
    IonCol,
    RouterLink,
    IonRouterLink
  ]
})
export class HomePage implements OnInit {
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);

  ngOnInit() {
    if (this.pocketbaseService.isLoggedIn()) {
      this.router.navigate(['/todos']);
    }
  }
}
