import { Routes } from '@angular/router';
import { authGuard, guestGuard } from './guards/auth.guard';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () => import('./login/login.page').then(m => m.LoginPage),
    canActivate: [guestGuard]
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./register/register.page').then(m => m.RegisterPage),
    canActivate: [guestGuard]
  },
  {
    path: 'password-reset',
    loadComponent: () =>
      import('./password-reset/password-reset.page').then(
        m => m.PasswordResetPage
      ),
    canActivate: [guestGuard]
  },
  {
    path: 'todos',
    loadComponent: () => import('./todos/todos.page').then(m => m.TodosPage),
    canActivate: [authGuard]
  },
  {
    path: 'edit-todo',
    loadComponent: () =>
      import('./edit-todo/edit-todo.page').then(m => m.EditTodoPage),
    canActivate: [authGuard]
  },
  {
    path: 'edit-todo/:id',
    loadComponent: () =>
      import('./edit-todo/edit-todo.page').then(m => m.EditTodoPage),
    canActivate: [authGuard]
  },
  {
    path: 'profile',
    loadComponent: () =>
      import('./profile/profile.page').then(m => m.ProfilePage),
    canActivate: [authGuard]
  },
  {
    path: 'home',
    loadComponent: () => import('./home/home.page').then(m => m.HomePage)
  },
  {
    path: '',
    redirectTo: 'login',
    pathMatch: 'full'
  }
];
