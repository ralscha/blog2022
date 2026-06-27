import { Component, inject, input, OnInit, signal } from '@angular/core';
import {
  FormField,
  FormRoot,
  form,
  maxLength,
  required
} from '@angular/forms/signals';
import { Router } from '@angular/router';
import {
  IonBackButton,
  IonButton,
  IonButtons,
  IonCard,
  IonCardContent,
  IonCheckbox,
  IonContent,
  IonDatetime,
  IonHeader,
  IonInput,
  IonLabel,
  IonText,
  IonTextarea,
  IonTitle,
  IonToolbar
} from '@ionic/angular/standalone';
import { PocketbaseService } from '../services/pocketbase.service';
import { Todo } from '../models/todo.model';
import { ToastService } from '../services/toast.service';
import { FormErrorService } from '../services/form-error.service';

type TodoForm = {
  title: string;
  description: string;
  completed: boolean;
  due_date: string;
};

@Component({
  selector: 'app-edit-todo',
  templateUrl: './edit-todo.page.html',
  styleUrl: './edit-todo.page.css',
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonCard,
    IonCardContent,
    IonLabel,
    IonInput,
    IonTextarea,
    IonButton,
    IonText,
    IonDatetime,
    IonCheckbox,
    IonButtons,
    IonBackButton,
    FormField,
    FormRoot
  ]
})
export class EditTodoPage implements OnInit {
  formErrorService = inject(FormErrorService);
  todoModel = signal<TodoForm>({
    title: '',
    description: '',
    completed: false,
    due_date: ''
  });
  todoForm = form(this.todoModel, path => {
    required(path.title);
    maxLength(path.title, 255);
    maxLength(path.description, 1000);
  });
  isLoading = signal(false);
  isEditing = signal(false);
  currentTodo = signal<Todo | null>(null);
  id = input<string | null>(null);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private toastService = inject(ToastService);

  ngOnInit(): void {
    if (this.id()) {
      this.isEditing.set(true);
      this.loadTodo();
    }
  }

  async loadTodo(): Promise<void> {
    if (!this.id()) {
      return;
    }

    this.isLoading.set(true);
    const todo = await this.pocketbaseService.getTodo(this.id()!);
    this.currentTodo.set(todo);
    this.todoModel.set({
      title: todo.title,
      description: todo.description || '',
      completed: todo.completed,
      due_date: todo.due_date || ''
    });
    this.isLoading.set(false);
  }

  async onSubmit(): Promise<void> {
    if (this.todoForm().valid() && !this.isLoading()) {
      this.isLoading.set(true);

      const formValue = this.todoModel();
      const formData = {
        ...formValue,
        due_date: formValue.due_date
          ? new Date(formValue.due_date).toISOString()
          : undefined
      };

      if (this.isEditing() && this.id()) {
        await this.pocketbaseService.updateTodo(this.id()!, formData);
        await this.toastService.showToast(
          'Todo updated successfully!',
          'success'
        );
      } else {
        await this.pocketbaseService.createTodo(formData);
        await this.toastService.showToast(
          'Todo created successfully!',
          'success'
        );
      }

      this.router.navigate(['/todos']);
      this.isLoading.set(false);
    }
  }

  getCurrentDate(): string {
    return new Date().toISOString().split('T')[0];
  }
}
