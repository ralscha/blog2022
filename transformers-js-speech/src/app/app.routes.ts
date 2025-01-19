import {Routes} from '@angular/router';

export const routes: Routes = [
  {
    path: 'home',
    loadComponent: () => import('./speech/speech.page').then((m) => m.SpeechPage),
  },
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full',
  },
];
