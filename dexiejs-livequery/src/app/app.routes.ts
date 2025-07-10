import {Routes} from '@angular/router';
import {EditPage} from "./todo/edit/edit.page";
import {ListPage} from "./todo/list/list.page";
import {TodoResolver} from "./todo/todo.resolver";

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'todo',
    pathMatch: 'full'
  },
  {
    path: 'todo',
    component: ListPage,
  },
  {
    path: 'todo/edit',
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
  },
  {path: '**', redirectTo: 'todo'}
];
