import { inject, ErrorHandler, Injectable } from '@angular/core';
import { ToastService } from './toast.service';

@Injectable()
export class GlobalErrorHandler implements ErrorHandler {
  private toastService = inject(ToastService);

  handleError(error: any): void {
    if (error.originalError) {
      error = error.originalError;
    }

    const message = error.message || 'An unexpected error occurred';

    console.error(error);

    void this.toastService.showToast(message, 'danger');
  }
}
