import {Component, inject} from '@angular/core';
import {Todo} from '../todo-db';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {liveQuery, Observable} from 'dexie';
import {RouterLink} from '@angular/router';
import {AsyncPipe, DatePipe} from '@angular/common';
import {
  IonCard,
  IonCardContent,
  IonCardHeader,
  IonContent,
  IonFab,
  IonFabButton,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonRouterLink,
  IonRow,
  IonTitle,
  IonToolbar
} from "@ionic/angular/standalone";
import {addIcons} from "ionicons";
import {add, ellipse, pencilOutline, trashOutline} from "ionicons/icons";

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  styleUrl: './list.page.scss',
  imports: [RouterLink, IonRouterLink, AsyncPipe, DatePipe, IonHeader, IonToolbar, IonTitle, IonContent, IonCard, IonCardHeader, IonIcon, IonLabel, IonCardContent, IonRow, IonItem, IonFab, IonFabButton]
})
export class ListPage {
  public readonly todos$: Observable<Todo[]>;
  private readonly todoService = inject(TodoService);
  private readonly messagesService = inject(MessagesService);

  constructor() {
    addIcons({ellipse, pencilOutline, trashOutline, add});
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
