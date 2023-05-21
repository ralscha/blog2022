import {inject, NgModule} from '@angular/core';
import {PreloadAllModules, RouterModule, Routes} from '@angular/router';
import {StartComponent} from './start/start.component';
import {AuthGuard} from './auth.guard';
import {SignOutComponent} from './sign-out/sign-out.component';

const routes: Routes = [
  {
    path: '',
    component: StartComponent
  },
  {
    path: 'sign-out',
    component: SignOutComponent
  },
  {
    path: 'todo',
    canActivate: [() => inject(AuthGuard).canActivate()],
    loadChildren: () => import('./todo/todo.module').then(m => m.TodoModule)
  },
  {path: '**', redirectTo: 'start'}
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, {preloadingStrategy: PreloadAllModules})
  ],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
