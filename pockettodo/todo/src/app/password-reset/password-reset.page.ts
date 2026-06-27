import { Component, inject, signal } from '@angular/core';
import {
  email,
  FormField,
  FormRoot,
  form,
  required
} from '@angular/forms/signals';
import { RouterLink } from '@angular/router';
import {
  IonButton,
  IonCard,
  IonCardContent,
  IonCol,
  IonContent,
  IonGrid,
  IonHeader,
  IonInput,
  IonRouterLinkWithHref,
  IonRow,
  IonText,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';
import { PocketbaseService } from '../services/pocketbase.service';
import { ToastService } from '../services/toast.service';
import { FormErrorService } from '../services/form-error.service';

@Component({
  selector: 'app-password-reset',
  templateUrl: './password-reset.page.html',
  styleUrl: './password-reset.page.css',
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonCard,
    IonCardContent,
    IonInput,
    IonButton,
    IonText,
    IonGrid,
    IonRow,
    IonCol,
    FormField,
    FormRoot,
    RouterLink,
    IonRouterLinkWithHref
  ]
})
export class PasswordResetPage {
  formErrorService = inject(FormErrorService);
  resetModel = signal({
    email: ''
  });
  resetForm = form(this.resetModel, path => {
    required(path.email);
    email(path.email);
  });
  isLoading = signal(false);
  emailSent = signal(false);
  private pocketbaseService = inject(PocketbaseService);
  private toastService = inject(ToastService);

  async onSubmit(): Promise<void> {
    if (this.resetForm().valid() && !this.isLoading()) {
      this.isLoading.set(true);

      await this.pocketbaseService.requestPasswordReset(
        this.resetModel().email
      );
      this.emailSent.set(true);
      await this.toastService.showToast(
        'Password reset email sent! Check your inbox.',
        'success'
      );

      this.isLoading.set(false);
    }
  }
}
