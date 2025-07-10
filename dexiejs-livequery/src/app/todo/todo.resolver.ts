import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from "@angular/router";
import {Todo} from "./todo-db";
import {Observable} from "rxjs";
import {TodoService} from "./todo.service";
import {v4 as uuidv4} from 'uuid';
import {inject, Injectable} from "@angular/core";
import {add, format} from 'date-fns'

@Injectable({providedIn: 'root'})
export class TodoResolver implements Resolve<Todo> {
  private readonly todoService = inject(TodoService);

  resolve(route: ActivatedRouteSnapshot, _state: RouterStateSnapshot): Observable<Todo> | Promise<Todo> | Todo {

    const id = route.paramMap.get('id');
    if (id) {
      return this.todoService.getTodo(id).then((result) => {
        if (result) {
          return result;
        }
        return this.returnEmpty();
      });
    }
    return this.returnEmpty();

  }

  private returnEmpty(): Todo {
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
