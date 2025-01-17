import {Injectable} from '@angular/core';
import {ToastController} from '@ionic/angular';

@Injectable({
  providedIn: 'root'
})
export class MessagesService {

  constructor(private readonly toastCtrl: ToastController) {
  }

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
