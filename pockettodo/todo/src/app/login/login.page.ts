import { Component, inject, signal } from '@angular/core';
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
  IonInputPasswordToggle,
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
  selector: 'app-login',
  templateUrl: './login.page.html',
  styleUrl: './login.page.css',
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
    IonRouterLinkWithHref,
    IonInputPasswordToggle
  ]
})
export class LoginPage {
  formErrorService = inject(FormErrorService);
  loginModel = signal({
    email: '',
    password: ''
  });
  loginForm = form(this.loginModel, path => {
    required(path.email);
    email(path.email);
    required(path.password);
    minLength(path.password, 6);
  });
  isLoading = signal(false);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private toastService = inject(ToastService);

  async onSubmit(): Promise<void> {
    if (this.loginForm().valid() && !this.isLoading()) {
      this.isLoading.set(true);

      await this.pocketbaseService.login(this.loginModel());
      await this.toastService.showToast('Login successful!', 'success');
      this.router.navigate(['/todos']);

      this.isLoading.set(false);
    }
  }
}
