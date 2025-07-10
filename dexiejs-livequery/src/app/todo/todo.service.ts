import {Injectable} from '@angular/core';
import {Todo, TodoDb} from './todo-db';

@Injectable({providedIn: 'root'})
export class TodoService {

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


}
