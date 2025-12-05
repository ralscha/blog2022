import { provideZoneChangeDetection } from "@angular/core";
import {provideRouter, RouteReuseStrategy, Routes, withHashLocation} from '@angular/router';
import {IonicRouteStrategy, provideIonicAngular} from '@ionic/angular/standalone';
import {bootstrapApplication} from '@angular/platform-browser';
import {provideHttpClient} from '@angular/common/http';
import {AppComponent} from './app/app.component';
import {HomePage} from "./app/home/home.page";

const routes: Routes = [
  {
    path: 'home',
    component: HomePage
  },
  {
    path: '',
    redirectTo: 'home',
    pathMatch: 'full'
  },
];

bootstrapApplication(AppComponent, {
  providers: [
    provideZoneChangeDetection(),provideIonicAngular(),
    provideHttpClient(),
    provideRouter(routes, withHashLocation()),
    {provide: RouteReuseStrategy, useClass: IonicRouteStrategy}
  ]
})
  .catch(err => console.log(err));
