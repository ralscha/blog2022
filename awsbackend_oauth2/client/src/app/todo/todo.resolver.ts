import {inject, Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {v4 as uuidv4} from 'uuid';
import {add, format} from 'date-fns'
import {Todo} from './todo';
import {TodoService} from "./todo.service";

@Injectable({providedIn: 'root'})
export class TodoResolver implements Resolve<Todo> {
  private readonly todoService = inject(TodoService);

  resolve(route: ActivatedRouteSnapshot, _state: RouterStateSnapshot): Observable<Todo> | Promise<Todo> | Todo {
    const id = route.paramMap.get('id');
    let todo: Todo | undefined = undefined;
    if (id) {
      todo = this.todoService.getTodo(id);
    }

    if (todo) {
      return todo;
    }

    return {
      id: uuidv4(),
      description: "",
      priority: "normal",
      dueDate: format(add(new Date(), {
        days: 3,
      }), "yyyy-MM-dd")
    }

  }


}
