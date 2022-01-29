import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, Resolve, RouterStateSnapshot} from '@angular/router';
import {BehaviorSubject, Observable, tap} from 'rxjs';
import {v4 as uuidv4} from 'uuid';
import {add, format} from 'date-fns'
import {Todo} from './todo';
import {TodoPostResponse} from './todo-post-response';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';

@Injectable()
export class TodoService implements Resolve<Todo> {
  private todosMap: Map<string, Todo> = new Map();

  private readonly todosSubject = new BehaviorSubject<Todo[]>([]);
  private readonly todos$ = this.todosSubject.asObservable();

  constructor(private readonly httpClient: HttpClient) {
  }

  loadTodos(): void {
    this.httpClient.get<Todo[]>(`${environment.API_URL}/todos`).subscribe(todos => {
      this.todosMap.clear();
      for (const todo of todos) {
        this.todosMap.set(todo.id, todo);
      }
      this.publish();
    });
  }

  getTodos(): Observable<Todo[]> {
    return this.todos$;
  }

  getTodo(id: string): Todo | undefined {
    return this.todosMap.get(id);
  }

  deleteTodo(id: string): Observable<void> {
    return this.httpClient.delete<void>(`${environment.API_URL}/todos/${id}`)
      .pipe(
        tap(() => {
          this.todosMap.delete(id);
          this.publish();
        }));
  }

  save(todo: Todo): Observable<TodoPostResponse> {
    return this.httpClient.post<TodoPostResponse>(`${environment.API_URL}/todos`, todo)
      .pipe(
        tap(() => {
          this.todosMap.set(todo.id, todo)
          this.publish();
        })
      )
  }

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<Todo> | Promise<Todo> | Todo {
    const id = route.paramMap.get('id');
    if (id) {
      // @ts-ignore
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

  private publish(): void {
    this.todosSubject.next([...this.todosMap.values()])
  }

}
