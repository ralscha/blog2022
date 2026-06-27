import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TodoService } from '../todo.service';
import { MessagesService } from '../messages.service';
import { FormField, FormRoot, form } from '@angular/forms/signals';
import { HttpErrorResponse } from '@angular/common/http';
import { TodoPostResponse } from '../todo-post-response';
import { Todo } from '../todo';
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
  IonToolbar,
} from '@ionic/angular/standalone';

type TodoForm = {
  description: string;
  priority: Todo['priority'];
  dueDate: string;
};

@Component({
  selector: 'app-edit-page',
  templateUrl: './edit.page.html',
  styleUrls: ['./edit.page.scss'],
  imports: [
    FormField,
    FormRoot,
    IonHeader,
    IonToolbar,
    IonButtons,
    IonBackButton,
    IonTitle,
    IonContent,
    IonList,
    IonItem,
    IonInput,
    IonSelect,
    IonSelectOption,
    IonModal,
    IonDatetime,
    IonButton,
  ],
})
export class EditPage implements OnInit {
  selectedTodo!: Todo;
  showCalendar = signal(false);
  fieldErrors = signal<{ [key: string]: string }>({});
  todoModel = signal<TodoForm>({
    description: '',
    priority: 'normal',
    dueDate: '',
  });
  todoForm = form(this.todoModel);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly messagesService = inject(MessagesService);
  private readonly todoService = inject(TodoService);

  async ngOnInit(): Promise<void> {
    this.route.data.subscribe((data) => {
      this.selectedTodo = data['todo'];
      this.todoModel.set({
        description: this.selectedTodo.description,
        priority: this.selectedTodo.priority,
        dueDate: this.selectedTodo.dueDate ?? '',
      });
    });
  }

  async save(): Promise<void> {
    const formValue = this.todoModel();
    this.selectedTodo.dueDate = formValue.dueDate || undefined;
    this.selectedTodo.description = formValue.description;
    this.selectedTodo.priority = formValue.priority;
    this.fieldErrors.set({});

    this.todoService.save(this.selectedTodo).subscribe({
      next: () => this.handleSuccessResponse(),
      error: this.handleErrorResponse(),
    });
  }

  dateChanged(event: CustomEvent<{ value?: string | string[] | null }>): void {
    const dueDate = event.detail.value;
    this.todoModel.update((value) => ({
      ...value,
      dueDate: typeof dueDate === 'string' ? dueDate : '',
    }));
    this.showCalendar.set(false);
  }

  private handleSuccessResponse(): void {
    this.messagesService.showSuccessToast('Todo successfully saved', 500);
    this.router.navigate(['/todo']);
  }

  private handleErrorResponse() {
    return (errorResponse: HttpErrorResponse) => {
      const response: TodoPostResponse = errorResponse.error;
      if (response) {
        if (response.fieldErrors) {
          this.fieldErrors.set(response.fieldErrors);
        } else {
          this.messagesService.showErrorToast('Saving Todo failed');
        }
      }
    };
  }
}
