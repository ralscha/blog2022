import { ErrorHandler, inject, Injectable, NgZone } from '@angular/core';
import { ToastService } from './toast.service';

@Injectable()
export class GlobalErrorHandler implements ErrorHandler {
  private toastService = inject(ToastService);
  private zone = inject(NgZone);

  handleError(error: any): void {
    if (error.originalError) {
      error = error.originalError;
    }

    const message = error.message || 'An unexpected error occurred';

    console.error(error);

    this.zone.run(() => {
      this.toastService.showToast(message, 'danger');
    });
  }
}
