import {bootstrapApplication} from '@angular/platform-browser';
import {provideRouter, RouteReuseStrategy, withHashLocation} from '@angular/router';
import {IonicRouteStrategy, provideIonicAngular} from '@ionic/angular/standalone';

import {routes} from './app/app.routes';
import {AppComponent} from './app/app.component';
import {provideHttpClient} from "@angular/common/http";


bootstrapApplication(AppComponent, {
  providers: [
    {provide: RouteReuseStrategy, useClass: IonicRouteStrategy},
    provideHttpClient(),
    provideIonicAngular(),
    provideRouter(routes, withHashLocation()),
  ],
});
