import { inject, Injectable } from '@angular/core';
import { ToastController } from '@ionic/angular/standalone';

@Injectable({
  providedIn: 'root'
})
export class ToastService {
  private toastController = inject(ToastController);

  async showToast(
    message: string,
    color: 'success' | 'danger' | 'warning' | 'medium' = 'medium',
    duration: number = 3000,
    position: 'top' | 'bottom' | 'middle' = 'top'
  ) {
    const toast = await this.toastController.create({
      message,
      duration,
      color,
      position
    });
    await toast.present();
  }
}
