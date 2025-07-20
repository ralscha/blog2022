import { Injectable, signal } from '@angular/core';
import PocketBase from 'pocketbase';
import { environment } from '../../environments/environment';
import {
  AuthData,
  LoginRequest,
  RegisterRequest,
  UpdateProfileRequest,
  User
} from '../models/user.model';
import {
  CreateTodoRequest,
  Todo,
  UpdateTodoRequest
} from '../models/todo.model';

@Injectable({
  providedIn: 'root'
})
export class PocketbaseService {
  isLoggedIn = signal<boolean>(false);
  currentUser = signal<User | null>(null);
  private pb = new PocketBase(environment.pocketbaseUrl);

  constructor() {
    this.checkAuth();

    this.pb.authStore.onChange(() => {
      this.isLoggedIn.set(this.pb.authStore.isValid);
      this.currentUser.set(
        (this.pb.authStore.record as unknown as User) || null
      );
    });
  }

  async login(credentials: LoginRequest): Promise<AuthData> {
    try {
      const authData = await this.pb
        .collection('users')
        .authWithPassword(credentials.email, credentials.password);
      return authData as unknown as AuthData;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async register(userData: RegisterRequest): Promise<User> {
    try {
      const user = await this.pb.collection('users').create(userData);
      return user as unknown as User;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async logout(): Promise<void> {
    this.pb.authStore.clear();
  }

  async requestPasswordReset(email: string): Promise<void> {
    try {
      await this.pb.collection('users').requestPasswordReset(email);
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async refreshAuth(): Promise<AuthData | null> {
    try {
      if (this.pb.authStore.isValid) {
        const authData = await this.pb.collection('users').authRefresh();
        return authData as unknown as AuthData;
      }
      return null;
    } catch (error) {
      this.logout();
      return null;
    }
  }

  async updateProfile(
    userId: string,
    data: UpdateProfileRequest
  ): Promise<User> {
    try {
      const user = await this.pb.collection('users').update(userId, data);
      return user as unknown as User;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async getTodos(hideCompleted: boolean): Promise<Todo[]> {
    try {
      let filter = undefined;
      if (hideCompleted) {
        filter = `completed = false`;
      }
      const result = await this.pb.collection('todos').getFullList({
        sort: '-created',
        filter: filter
      });
      return result as unknown as Todo[];
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async createTodo(todoData: Omit<CreateTodoRequest, 'user'>): Promise<Todo> {
    try {
      const data = {
        ...todoData,
        user: this.currentUser()?.id
      };
      const todo = await this.pb.collection('todos').create(data);
      return todo as unknown as Todo;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async updateTodo(todoId: string, data: UpdateTodoRequest): Promise<Todo> {
    try {
      const todo = await this.pb.collection('todos').update(todoId, data);
      return todo as unknown as Todo;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async deleteTodo(todoId: string): Promise<void> {
    try {
      await this.pb.collection('todos').delete(todoId);
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async getTodo(todoId: string): Promise<Todo> {
    try {
      const todo = await this.pb.collection('todos').getOne(todoId);
      return todo as unknown as Todo;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  private checkAuth(): void {
    this.isLoggedIn.set(this.pb.authStore.isValid);
    this.currentUser.set((this.pb.authStore.record as unknown as User) || null);
  }

  private handleError(error: any): Error {
    if (error?.response?.data) {
      const errorData = error.response.data;
      if (errorData.message) {
        return new Error(errorData.message);
      }
      if (errorData.data) {
        const firstError = Object.values(errorData.data)[0] as any;
        if (firstError?.message) {
          return new Error(firstError.message);
        }
      }
    }
    return new Error(error?.message || 'An unexpected error occurred');
  }
}
