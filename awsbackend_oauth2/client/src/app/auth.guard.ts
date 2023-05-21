import {Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {OAuthService} from 'angular-oauth2-oidc';

@Injectable()
export class AuthGuard {
  constructor(private readonly router: Router, private readonly oauthService: OAuthService) {
  }

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
