import {provideRouter, RouteReuseStrategy, withHashLocation} from '@angular/router';
import {provideHttpClient} from '@angular/common/http';
import {bootstrapApplication} from '@angular/platform-browser';
import {AppComponent} from './app/app.component';
import {IonicRouteStrategy, provideIonicAngular} from '@ionic/angular/standalone';
import {routes} from "./app/app-routing.routes";


bootstrapApplication(AppComponent, {
  providers: [
    {provide: RouteReuseStrategy, useClass: IonicRouteStrategy},
    provideHttpClient(),
    provideIonicAngular(),
    provideRouter(routes, withHashLocation()),
  ]
})
  .catch(err => console.error(err));
