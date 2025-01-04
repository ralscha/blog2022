import {Component, OnInit} from '@angular/core';
import {OAuthService} from 'angular-oauth2-oidc';
import {authCodeFlowConfig} from '../auth.config';
import {Router} from '@angular/router';

@Component({
    selector: 'app-start',
    templateUrl: './start.component.html',
    standalone: false
})
export class StartComponent implements OnInit {

  constructor(private readonly oauthService: OAuthService,
              private readonly router: Router) {
  }

  ngOnInit() {
    this.oauthService.configure(authCodeFlowConfig);
    this.oauthService.loadDiscoveryDocumentAndLogin().then(success => {
      if (success) {
        this.oauthService.setupAutomaticSilentRefresh();
        this.router.navigate(['/todo']);
      }
    });
  }

}
