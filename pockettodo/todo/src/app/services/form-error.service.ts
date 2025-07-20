import { Injectable } from '@angular/core';
import { AbstractControl } from '@angular/forms';

@Injectable({
  providedIn: 'root'
})
export class FormErrorService {
  constructor() {}

  getErrorMessage(
    control: AbstractControl | null,
    controlName: string
  ): string {
    if (!control || !control.invalid || !control.touched) {
      return '';
    }

    if (control.errors?.['required']) {
      return `${this.capitalizeFirstLetter(controlName)} is required`;
    }

    if (control.errors?.['email']) {
      return 'Please enter a valid email';
    }

    if (control.errors?.['minlength']) {
      return `${this.capitalizeFirstLetter(controlName)} must be at least ${control.errors?.['minlength'].requiredLength} characters`;
    }

    if (control.errors?.['maxlength']) {
      return `${this.capitalizeFirstLetter(controlName)} cannot exceed ${control.errors?.['maxlength'].requiredLength} characters`;
    }

    if (control.errors?.['passwordMismatch']) {
      return 'Passwords do not match';
    }

    return '';
  }

  private capitalizeFirstLetter(string: string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }
}
