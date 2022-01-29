import {AuthConfig} from 'angular-oauth2-oidc';
import {environment} from '../environments/environment';

export const authCodeFlowConfig: AuthConfig = {
  issuer: environment.ISSUER,
  redirectUri: 'http://localhost:8100',
  logoutUrl: 'http://localhost:8100/sign-out',
  clientId: environment.CLIENT_ID,
  responseType: 'code',
  scope: 'openid',
  showDebugInformation: !environment.production,
  timeoutFactor: 0.01,
  strictDiscoveryDocumentValidation: false
};
