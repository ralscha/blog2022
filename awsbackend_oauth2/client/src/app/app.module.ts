import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {RouteReuseStrategy} from '@angular/router';
import {IonicModule, IonicRouteStrategy} from '@ionic/angular';
import {AppComponent} from './app.component';
import {AppRoutingModule} from './app-routing.module';
import {OAuthModule} from 'angular-oauth2-oidc';
import {HttpClientModule} from '@angular/common/http';
import {AuthGuard} from './auth.guard';
import {StartComponent} from './start/start.component';
import {environment} from '../environments/environment';
import {SignOutComponent} from './sign-out/sign-out.component';

@NgModule({
  declarations: [AppComponent, StartComponent, SignOutComponent],
  imports: [BrowserModule, HttpClientModule, OAuthModule.forRoot({
    resourceServer: {
      allowedUrls: [environment.API_URL],
      sendAccessToken: true
    }
  }), IonicModule.forRoot(), AppRoutingModule],
  providers: [{
    provide: RouteReuseStrategy,
    useClass: IonicRouteStrategy
  }, AuthGuard/*, { provide: OAuthStorage, useValue: localStorage }*/],
  bootstrap: [AppComponent],
})
export class AppModule {
}
