import { Component, inject, OnInit, signal } from '@angular/core';
import {
  email,
  FormField,
  FormRoot,
  form,
  required
} from '@angular/forms/signals';
import { Router } from '@angular/router';
import {
  AlertController,
  IonBackButton,
  IonButton,
  IonButtons,
  IonCard,
  IonCardContent,
  IonCardHeader,
  IonCardTitle,
  IonContent,
  IonHeader,
  IonInput,
  IonItem,
  IonLabel,
  IonText,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';
import { PocketbaseService } from '../services/pocketbase.service';
import { User } from '../models/user.model';
import { ToastService } from '../services/toast.service';
import { FormErrorService } from '../services/form-error.service';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.page.html',
  styleUrl: './profile.page.css',
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonCard,
    IonCardHeader,
    IonCardTitle,
    IonCardContent,
    IonLabel,
    IonInput,
    IonButton,
    IonText,
    IonButtons,
    IonBackButton,
    FormField,
    FormRoot,
    IonItem
  ]
})
export class ProfilePage implements OnInit {
  formErrorService = inject(FormErrorService);
  profileModel = signal({
    email: '',
    name: ''
  });
  profileForm = form(this.profileModel, path => {
    required(path.email);
    email(path.email);
  });
  isLoading = signal(false);
  currentUser = signal<User | null>(null);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private toastService = inject(ToastService);
  private alertController = inject(AlertController);

  ngOnInit(): void {
    this.loadProfile();
  }

  async loadProfile(): Promise<void> {
    this.isLoading.set(true);
    const user = this.pocketbaseService.currentUser();
    if (user) {
      this.currentUser.set(user);
      this.profileModel.set({
        email: user.email,
        name: user.name || ''
      });
    }
    this.isLoading.set(false);
  }

  async onSubmit(): Promise<void> {
    if (this.profileForm().valid() && !this.isLoading()) {
      const user = this.currentUser();
      if (!user) return;

      this.isLoading.set(true);

      const updatedUser = await this.pocketbaseService.updateProfile(
        user.id,
        this.profileModel()
      );
      this.currentUser.set(updatedUser);
      await this.toastService.showToast(
        'Profile updated successfully!',
        'success'
      );

      this.isLoading.set(false);
    }
  }

  async requestPasswordReset() {
    const user = this.currentUser();
    if (!user) return;

    const alert = await this.alertController.create({
      header: 'Password Reset',
      message: `Send password reset email to ${user.email}?`,
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel'
        },
        {
          text: 'Send',
          handler: async () => {
            await this.pocketbaseService.requestPasswordReset(user.email);
            await this.toastService.showToast(
              'Password reset email sent!',
              'success'
            );
          }
        }
      ]
    });

    await alert.present();
  }

  async logout() {
    const alert = await this.alertController.create({
      header: 'Logout',
      message: 'Are you sure you want to logout?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel'
        },
        {
          text: 'Logout',
          handler: async () => {
            await this.pocketbaseService.logout();
            this.router.navigate(['/login']);
          }
        }
      ]
    });

    await alert.present();
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }
}
