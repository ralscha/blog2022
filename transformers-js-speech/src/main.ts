import {bootstrapApplication} from '@angular/platform-browser';
import {PreloadAllModules, provideRouter, RouteReuseStrategy, withPreloading, withHashLocation} from '@angular/router';
import {IonicRouteStrategy, provideIonicAngular} from '@ionic/angular/standalone';

import {routes} from './app/app.routes';
import {AppComponent} from './app/app.component';

bootstrapApplication(AppComponent, {
  providers: [
    {provide: RouteReuseStrategy, useClass: IonicRouteStrategy},
    provideIonicAngular(),
    provideRouter(routes, withHashLocation(), withPreloading(PreloadAllModules)),
  ],
});
