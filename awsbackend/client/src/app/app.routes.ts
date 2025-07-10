import {TodoResolver} from "./todo/todo.resolver";
import {EditPage} from "./todo/edit/edit.page";
import {ListPage} from "./todo/list/list.page";
import {Routes} from "@angular/router";

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'todo',
    pathMatch: 'full'
  },
  {
    path: 'todo',
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
  {path: '**', redirectTo: 'todo'}
];
