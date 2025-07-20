import { Component, inject, OnInit, signal } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
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
    ReactiveFormsModule,
    IonItem
  ]
})
export class ProfilePage implements OnInit {
  formErrorService = inject(FormErrorService);
  profileForm: FormGroup;
  isLoading = signal(false);
  currentUser = signal<User | null>(null);
  private fb = inject(FormBuilder);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private toastService = inject(ToastService);
  private alertController = inject(AlertController);

  constructor() {
    this.profileForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      name: ['']
    });
  }

  get email() {
    return this.profileForm.get('email');
  }

  get name() {
    return this.profileForm.get('name');
  }

  ngOnInit() {
    this.loadProfile();
  }

  async loadProfile() {
    this.isLoading.set(true);
    const user = this.pocketbaseService.currentUser();
    if (user) {
      this.currentUser.set(user);
      this.profileForm.patchValue({
        email: user.email,
        name: user.name || ''
      });
    }
    this.isLoading.set(false);
  }

  async onSubmit() {
    if (this.profileForm.valid && !this.isLoading()) {
      const user = this.currentUser();
      if (!user) return;

      this.isLoading.set(true);

      const updatedUser = await this.pocketbaseService.updateProfile(
        user.id,
        this.profileForm.value
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
