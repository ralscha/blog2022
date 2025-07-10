import {Component, inject, OnInit} from '@angular/core';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {Observable} from 'rxjs';
import {Todo} from '../todo';
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
import {add, ellipse, exitOutline, pencilOutline, trashOutline} from "ionicons/icons";

@Component({
  selector: 'app-list',
  templateUrl: './list.page.html',
  styleUrls: ['./list.page.scss'],
  imports: [RouterLink, IonRouterLink, AsyncPipe, DatePipe, IonHeader, IonToolbar, IonTitle, IonIcon, IonContent, IonCard, IonCardHeader, IonCardContent, IonLabel, IonRow, IonItem, IonFab, IonFabButton]
})
export class ListPage implements OnInit {
  todos$!: Observable<Todo[]>;
  private readonly todoService = inject(TodoService);
  private readonly messagesService = inject(MessagesService);

  constructor() {
    addIcons({exitOutline, ellipse, pencilOutline, trashOutline, add});
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
