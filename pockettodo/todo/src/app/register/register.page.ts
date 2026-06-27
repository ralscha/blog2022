import { Component, computed, inject, signal } from '@angular/core';
import {
  email,
  FormField,
  FormRoot,
  form,
  minLength,
  required
} from '@angular/forms/signals';
import { Router, RouterLink } from '@angular/router';
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
  selector: 'app-register',
  templateUrl: './register.page.html',
  styleUrl: './register.page.css',
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
export class RegisterPage {
  formErrorService = inject(FormErrorService);
  registerModel = signal({
    email: '',
    name: '',
    password: '',
    passwordConfirm: ''
  });
  registerForm = form(this.registerModel, path => {
    required(path.email);
    email(path.email);
    required(path.password);
    minLength(path.password, 6);
    required(path.passwordConfirm);
  });
  passwordMismatch = computed(() => {
    const { password, passwordConfirm } = this.registerModel();
    return passwordConfirm.length > 0 && password !== passwordConfirm;
  });
  isLoading = signal(false);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private toastService = inject(ToastService);

  async onSubmit(): Promise<void> {
    if (
      this.registerForm().valid() &&
      !this.passwordMismatch() &&
      !this.isLoading()
    ) {
      this.isLoading.set(true);

      await this.pocketbaseService.register(this.registerModel());
      await this.toastService.showToast(
        'Registration successful! You can now login.',
        'success'
      );
      this.router.navigate(['/login']);

      this.isLoading.set(false);
    }
  }
}
