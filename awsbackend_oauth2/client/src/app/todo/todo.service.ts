import {inject, Injectable} from '@angular/core';
import {BehaviorSubject, Observable, tap} from 'rxjs';
import {Todo} from './todo';
import {TodoPostResponse} from './todo-post-response';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';

@Injectable({providedIn: 'root'})
export class TodoService {
  private readonly httpClient = inject(HttpClient);

  private todosMap: Map<string, Todo> = new Map();

  private readonly todosSubject = new BehaviorSubject<Todo[]>([]);
  private readonly todos$ = this.todosSubject.asObservable();

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

  private publish(): void {
    this.todosSubject.next([...this.todosMap.values()])
  }

}
