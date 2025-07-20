import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { PocketbaseService } from '../services/pocketbase.service';

export const authGuard: CanActivateFn = () => {
  const pocketbaseService = inject(PocketbaseService);
  const router = inject(Router);

  if (pocketbaseService.isLoggedIn()) {
    return true;
  } else {
    router.navigate(['/login']);
    return false;
  }
};

export const guestGuard: CanActivateFn = () => {
  const pocketbaseService = inject(PocketbaseService);
  const router = inject(Router);

  if (!pocketbaseService.isLoggedIn()) {
    return true;
  } else {
    router.navigate(['/todos']);
    return false;
  }
};
