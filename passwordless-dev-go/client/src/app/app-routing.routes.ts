import {inject} from '@angular/core';
import {Router, Routes} from '@angular/router';
import {AuthenticationPage} from './authentication/authentication.page';
import {AuthService} from "./auth.service";
import {map} from "rxjs/operators";

export const authGuard = (authService = inject(AuthService), router = inject(Router)) => {
  if (authService.isLoggedIn()) {
    return true;
  }

  return authService.isAuthenticated().pipe(
    map(success => {
      if (success) {
        return true;
      }
      return router.createUrlTree(['/login']);
    })
  );
}

export const routes: Routes = [
  {path: '', redirectTo: 'home', pathMatch: 'full'},
  {
    path: 'home',
    canActivate: [() => authGuard()],
    loadComponent: () => import('./home/home.page').then(m => m.HomePage)
  },
  {path: 'login', component: AuthenticationPage},
  {
    path: 'registration',
    loadComponent: () => import('./registration/registration.page').then(m => m.RegistrationPage)
  }
];
