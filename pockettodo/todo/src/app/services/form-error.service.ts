import { Service } from '@angular/core';
import { type FieldTree } from '@angular/forms/signals';

@Service()
export class FormErrorService {
  getErrorMessage(
    field: FieldTree<unknown> | null,
    controlName: string
  ): string {
    const state = field?.();
    if (!state || !state.invalid() || !state.touched()) {
      return '';
    }

    if (state.getError('required')) {
      return `${this.capitalizeFirstLetter(controlName)} is required`;
    }

    if (state.getError('email')) {
      return 'Please enter a valid email';
    }

    const minLength = state.getError('minLength');
    if (minLength) {
      return `${this.capitalizeFirstLetter(controlName)} must be at least ${minLength.minLength} characters`;
    }

    const maxLength = state.getError('maxLength');
    if (maxLength) {
      return `${this.capitalizeFirstLetter(controlName)} cannot exceed ${maxLength.maxLength} characters`;
    }

    if (state.getError('passwordMismatch')) {
      return 'Passwords do not match';
    }

    return '';
  }

  private capitalizeFirstLetter(string: string) {
    return string.charAt(0).toUpperCase() + string.slice(1);
  }
}
