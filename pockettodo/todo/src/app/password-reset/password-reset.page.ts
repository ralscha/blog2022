import { Component, inject, signal } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
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
    ReactiveFormsModule,
    RouterLink,
    IonRouterLinkWithHref
  ]
})
export class PasswordResetPage {
  formErrorService = inject(FormErrorService);
  resetForm: FormGroup;
  isLoading = signal(false);
  emailSent = signal(false);
  private fb = inject(FormBuilder);
  private pocketbaseService = inject(PocketbaseService);
  private toastService = inject(ToastService);

  constructor() {
    this.resetForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]]
    });
  }

  get email() {
    return this.resetForm.get('email');
  }

  async onSubmit() {
    if (this.resetForm.valid && !this.isLoading()) {
      this.isLoading.set(true);

      await this.pocketbaseService.requestPasswordReset(
        this.resetForm.value.email
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
