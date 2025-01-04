import {Component} from '@angular/core';
import {Todo} from '../todo-db';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {liveQuery, Observable} from 'dexie';

@Component({
    selector: 'app-list',
    templateUrl: './list.page.html',
    styleUrls: ['./list.page.scss'],
    standalone: false
})
export class ListPage {
  public readonly todos$: Observable<Todo[]>;

  constructor(private readonly todoService: TodoService,
              private readonly messagesService: MessagesService) {
    this.todos$ = liveQuery(() => this.todoService.allTodos());
  }

  todoTrackBy(index: number, todo: Todo): string {
    return todo.id;
  }

  async deleteTodo(id: string): Promise<void> {
    await this.todoService.deleteTodo(id);
    await this.messagesService.showSuccessToast('Todo successfully deleted', 500);
  }


}
