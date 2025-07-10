import {inject} from '@angular/core';
import {Routes} from '@angular/router';
import {StartComponent} from './start/start.component';
import {AuthGuard} from './auth.guard';
import {SignOutComponent} from './sign-out/sign-out.component';
import {ListPage} from "./todo/list/list.page";
import {EditPage} from "./todo/edit/edit.page";
import {TodoResolver} from "./todo/todo.resolver";

export const routes: Routes = [
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
    children: [
      {
        path: '',
        component: ListPage,
      },
      {
        path: 'edit',
        children: [
          {
            path: ':id',
            component: EditPage,
            resolve: {
              todo: TodoResolver
            }
          },
          {
            path: '',
            component: EditPage,
            resolve: {
              todo: TodoResolver
            }
          }
        ]
      }
    ]
  },
  {path: '**', redirectTo: 'start'}
];
