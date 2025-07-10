import {PreloadAllModules, provideRouter, RouteReuseStrategy, withHashLocation, withPreloading} from '@angular/router';
import {IonicRouteStrategy, provideIonicAngular} from '@ionic/angular/standalone';
import {bootstrapApplication} from '@angular/platform-browser';
import {AppComponent} from './app/app.component';
import {routes} from './app/app.routes';
import {provideHttpClient} from "@angular/common/http";

bootstrapApplication(AppComponent, {
  providers: [
    provideIonicAngular(),
    provideHttpClient(),
    provideRouter(routes, withHashLocation(), withPreloading(PreloadAllModules)),
    {provide: RouteReuseStrategy, useClass: IonicRouteStrategy}
  ]
})
  .catch(err => console.error(err));
