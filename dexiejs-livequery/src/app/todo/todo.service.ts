import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs';
import {v4 as uuidv4} from 'uuid';
import {add, format} from 'date-fns'
import {Todo, TodoDb} from './todo-db';

@Injectable()
export class TodoService implements Resolve<Todo | undefined> {

  private readonly db: TodoDb;

  constructor() {
    this.db = new TodoDb();
  }

  updateTodo(todo: Todo): Promise<string> {
    return this.db.todos.put(todo);
  }

  deleteTodo(id: string): Promise<void> {
    return this.db.todos.delete(id);
  }

  allTodos(): Promise<Todo[]> {
    return this.db.todos.toArray();
  }

  async getTodo(id: string): Promise<Todo | undefined> {
    return this.db.todos.get(id);
  }

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Todo | undefined> | Promise<Todo | undefined> | Todo | undefined {
    const id = route.paramMap.get('id');
    if (id) {
      return this.getTodo(id);
    } else {
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

}
