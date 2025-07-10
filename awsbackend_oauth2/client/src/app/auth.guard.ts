import {inject, Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {OAuthService} from 'angular-oauth2-oidc';

@Injectable()
export class AuthGuard {
  private readonly router = inject(Router);
  private readonly oauthService = inject(OAuthService);


  canActivate() {
    if (
      this.oauthService.hasValidAccessToken() &&
      this.oauthService.hasValidIdToken()
    ) {
      return true;
    } else {
      this.router.navigate(['/']);
      return false;
    }
  }
}
