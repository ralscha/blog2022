import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {IonicModule} from '@ionic/angular';
import {FormsModule} from '@angular/forms';
import {ListPage} from './list/list.page';
import {RouterModule, Routes} from '@angular/router';
import {TodoService} from './todo.service';
import {EditPage} from './edit/edit.page';
import {MessagesService} from './messages.service';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';


const routes: Routes = [
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
          todo: TodoService
        }
      },
      {
        path: '',
        component: EditPage,
        resolve: {
          todo: TodoService
        }
      }
    ]
  }
];

@NgModule({ declarations: [ListPage, EditPage], imports: [CommonModule,
        FormsModule,
        IonicModule,
        RouterModule.forChild(routes)], providers: [TodoService, MessagesService, provideHttpClient(withInterceptorsFromDi())] })
export class TodoModule {
}
