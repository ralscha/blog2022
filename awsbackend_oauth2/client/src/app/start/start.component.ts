import { Component, inject, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { OAuthService } from 'angular-oauth2-oidc';
import { authCodeFlowConfig } from '../auth.config';
import { Router } from '@angular/router';

@Component({
  selector: 'app-start',
  changeDetection: ChangeDetectionStrategy.Eager,
  templateUrl: './start.component.html',
})
export class StartComponent implements OnInit {
  private readonly oauthService = inject(OAuthService);
  private readonly router = inject(Router);

  ngOnInit() {
    this.oauthService.configure(authCodeFlowConfig);
    this.oauthService.loadDiscoveryDocumentAndLogin().then((success) => {
      if (success) {
        this.oauthService.setupAutomaticSilentRefresh();
        this.router.navigate(['/todo']);
      }
    });
  }
}
