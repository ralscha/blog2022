import {Component, inject, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {FormsModule, NgForm} from '@angular/forms';
import {Todo} from '../todo-db';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonContent,
  IonDatetime,
  IonHeader,
  IonInput,
  IonItem,
  IonList,
  IonModal,
  IonSelect,
  IonSelectOption,
  IonTitle,
  IonToolbar
} from "@ionic/angular/standalone";

@Component({
  selector: 'app-edit-page',
  templateUrl: './edit.page.html',
  styleUrl: './edit.page.scss',
  imports: [FormsModule, IonHeader, IonToolbar, IonButtons, IonBackButton, IonTitle, IonContent, IonList, IonItem, IonInput, IonSelect, IonSelectOption, IonModal, IonDatetime, IonButton]
})
export class EditPage implements OnInit {
  selectedTodo!: Todo;
  showCalendar = false;
  dueDate?: string
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly messagesService = inject(MessagesService);
  private readonly todoService = inject(TodoService);

  async ngOnInit(): Promise<void> {
    this.route.data.subscribe(data => {
      this.selectedTodo = data['todo'];
      this.dueDate = this.selectedTodo.dueDate;
    });
  }

  async save(todoForm: NgForm) {
    this.selectedTodo.dueDate = this.dueDate;
    this.selectedTodo.description = todoForm.value.description;
    this.selectedTodo.priority = todoForm.value.priority;

    await this.todoService.updateTodo(this.selectedTodo);
    this.messagesService.showSuccessToast('Todo successfully saved', 500);
    this.router.navigate(['/todo']);
  }

  dateChanged(event: any) {
    this.dueDate = event.detail.value;
    this.showCalendar = false;
  }
}
