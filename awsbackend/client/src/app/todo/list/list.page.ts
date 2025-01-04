import {Component, OnInit} from '@angular/core';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {Observable} from 'rxjs';
import {Todo} from '../todo';

@Component({
    selector: 'app-list',
    templateUrl: './list.page.html',
    styleUrls: ['./list.page.scss'],
    standalone: false
})
export class ListPage implements OnInit {
  todos$!: Observable<Todo[]>;

  constructor(private readonly todoService: TodoService,
              private readonly messagesService: MessagesService) {
  }

  ngOnInit(): void {
    this.todos$ = this.todoService.getTodos();
    this.loadData();
  }

  async loadData(): Promise<void> {
    const loading = await this.messagesService.showLoading("Please wait...");
    try {
      await this.todoService.loadTodos();
    } finally {
      await loading.dismiss();
    }
  }

  todoTrackBy(index: number, todo: Todo): string {
    return todo.id;
  }

  async deleteTodo(id: string): Promise<void> {
    this.todoService.deleteTodo(id).subscribe(() => this.messagesService.showSuccessToast('Todo successfully deleted', 500));
  }


}
