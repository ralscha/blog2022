import {inject, Injectable} from '@angular/core';
import {ToastController} from '@ionic/angular/standalone';

@Injectable({
  providedIn: 'root'
})
export class MessagesService {
  private readonly toastCtrl = inject(ToastController);


  async showErrorToast(message = 'Unexpected error occurred'): Promise<void> {
    const toast = await this.toastCtrl.create({
      message,
      duration: 4000,
      position: 'bottom',
      color: 'danger'
    });
    await toast.present();
  }

}
