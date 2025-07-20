import { Component, inject, input, OnInit, signal } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
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
    ReactiveFormsModule
  ]
})
export class EditTodoPage implements OnInit {
  formErrorService = inject(FormErrorService);
  todoForm: FormGroup;
  isLoading = signal(false);
  isEditing = signal(false);
  currentTodo = signal<Todo | null>(null);
  id = input<string | null>(null);
  private fb = inject(FormBuilder);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private toastService = inject(ToastService);

  constructor() {
    this.todoForm = this.fb.group({
      title: ['', [Validators.required, Validators.maxLength(255)]],
      description: ['', [Validators.maxLength(1000)]],
      completed: [false],
      due_date: ['']
    });
  }

  get title() {
    return this.todoForm.get('title');
  }

  get description() {
    return this.todoForm.get('description');
  }

  get completed() {
    return this.todoForm.get('completed');
  }

  get due_date() {
    return this.todoForm.get('due_date');
  }

  ngOnInit() {
    if (this.id()) {
      this.isEditing.set(true);
      this.loadTodo();
    }
  }

  async loadTodo() {
    if (!this.id()) {
      return;
    }

    this.isLoading.set(true);
    const todo = await this.pocketbaseService.getTodo(this.id()!);
    this.currentTodo.set(todo);
    this.todoForm.patchValue({
      title: todo.title,
      description: todo.description || '',
      completed: todo.completed,
      due_date: todo.due_date || ''
    });
    this.isLoading.set(false);
  }

  async onSubmit() {
    if (this.todoForm.valid && !this.isLoading()) {
      this.isLoading.set(true);

      const formData = this.todoForm.value;

      if (formData.due_date) {
        formData.due_date = new Date(formData.due_date).toISOString();
      } else {
        formData.due_date = null;
      }

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
