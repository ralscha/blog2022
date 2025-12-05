import {PreloadAllModules, provideRouter, RouteReuseStrategy, withHashLocation, withPreloading} from '@angular/router';
import {IonicRouteStrategy, provideIonicAngular} from '@ionic/angular/standalone';
import {AuthGuard} from './app/auth.guard';
import {provideHttpClient, withInterceptorsFromDi} from '@angular/common/http';
import {bootstrapApplication} from '@angular/platform-browser';
import {OAuthModule} from 'angular-oauth2-oidc';
import {environment} from './environments/environment';
import {routes} from './app/app.routes';
import {AppComponent} from './app/app.component';
import {importProvidersFrom, provideZoneChangeDetection} from '@angular/core';


bootstrapApplication(AppComponent, {
  providers: [
    provideZoneChangeDetection(),importProvidersFrom(OAuthModule.forRoot({
      resourceServer: {
        allowedUrls: [environment.API_URL],
        sendAccessToken: true
      }
    })),
    provideIonicAngular(),
    provideRouter(routes, withHashLocation(), withPreloading(PreloadAllModules)),
    {
      provide: RouteReuseStrategy,
      useClass: IonicRouteStrategy
    }, AuthGuard /*, { provide: OAuthStorage, useValue: localStorage }*/, provideHttpClient(withInterceptorsFromDi())
  ]
})
  .catch(err => console.error(err));
