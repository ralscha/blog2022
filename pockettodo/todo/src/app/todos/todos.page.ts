import { Component, inject, signal } from '@angular/core';
import { Router } from '@angular/router';
import {
  AlertController,
  IonButton,
  IonButtons,
  IonCheckbox,
  IonContent,
  IonFab,
  IonFabButton,
  IonHeader,
  IonIcon,
  IonItem,
  IonItemOption,
  IonItemOptions,
  IonItemSliding,
  IonLabel,
  IonList,
  IonRefresher,
  IonRefresherContent,
  IonText,
  IonTitle,
  IonToolbar,
  ViewWillEnter
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import {
  add,
  create,
  eye,
  eyeOff,
  logOut,
  person,
  refresh,
  trash
} from 'ionicons/icons';
import { PocketbaseService } from '../services/pocketbase.service';
import { Todo } from '../models/todo.model';
import { ToastService } from '../services/toast.service';

@Component({
  selector: 'app-todos',
  templateUrl: './todos.page.html',
  styleUrl: './todos.page.css',
  imports: [
    IonHeader,
    IonToolbar,
    IonTitle,
    IonContent,
    IonFab,
    IonFabButton,
    IonIcon,
    IonList,
    IonLabel,
    IonCheckbox,
    IonButton,
    IonButtons,
    IonItemSliding,
    IonItemOptions,
    IonItemOption,
    IonText,
    IonRefresher,
    IonRefresherContent,
    IonItem
  ]
})
export class TodosPage implements ViewWillEnter {
  todos = signal<Todo[]>([]);
  isLoading = signal(false);
  hideCompleted = signal(false);
  private pocketbaseService = inject(PocketbaseService);
  private router = inject(Router);
  private alertController = inject(AlertController);
  private toastService = inject(ToastService);

  constructor() {
    addIcons({ add, create, trash, logOut, person, refresh, eyeOff, eye });
  }

  ionViewWillEnter() {
    this.loadTodos();
  }

  async loadTodos() {
    this.isLoading.set(true);
    const todos = await this.pocketbaseService.getTodos(this.hideCompleted());
    this.todos.set(todos);
    this.isLoading.set(false);
  }

  async toggleHideCompleted() {
    this.hideCompleted.update(v => !v);
    await this.loadTodos();
  }

  async toggleTodo(todo: Todo) {
    const updatedTodo = await this.pocketbaseService.updateTodo(todo.id, {
      completed: !todo.completed
    });

    this.todos.update(todos =>
      todos.map(t => (t.id === updatedTodo.id ? updatedTodo : t))
    );

    await this.toastService.showToast(
      `Todo ${updatedTodo.completed ? 'completed' : 'reopened'}!`,
      'success',
      2000,
      'bottom'
    );
  }

  async deleteTodo(todo: Todo) {
    const alert = await this.alertController.create({
      header: 'Delete Todo',
      message: `Are you sure you want to delete "${todo.title}"?`,
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel'
        },
        {
          text: 'Delete',
          role: 'destructive',
          handler: async () => {
            await this.pocketbaseService.deleteTodo(todo.id);
            this.todos.update(todos => todos.filter(t => t.id !== todo.id));
            await this.toastService.showToast('Todo deleted!', 'success');
          }
        }
      ]
    });

    await alert.present();
  }

  editTodo(todo: Todo) {
    this.router.navigate(['/edit-todo', todo.id]);
  }

  createTodo() {
    this.router.navigate(['/edit-todo']);
  }

  async logout() {
    const alert = await this.alertController.create({
      header: 'Logout',
      message: 'Are you sure you want to logout?',
      buttons: [
        {
          text: 'Cancel',
          role: 'cancel'
        },
        {
          text: 'Logout',
          handler: async () => {
            await this.pocketbaseService.logout();
            this.router.navigate(['/login']);
          }
        }
      ]
    });

    await alert.present();
  }

  goToProfile() {
    this.router.navigate(['/profile']);
  }

  async doRefresh(event: any) {
    await this.loadTodos();
    event.target.complete();
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString();
  }

  getDueDateColor(dueDate: string): string {
    if (!dueDate) {
      return '';
    }

    const due = new Date(dueDate);
    const now = new Date();
    const diffTime = due.getTime() - now.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays < 0) {
      return 'danger'; // Overdue
    }
    if (diffDays <= 1) {
      return 'warning'; // Due today or tomorrow
    }
    return 'medium'; // Future
  }
}
