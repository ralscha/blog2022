export interface Todo {
  id: string;
  title: string;
  description?: string;
  completed: boolean;
  user: string;
  due_date?: string;
  created: string;
  updated: string;
}

export interface CreateTodoRequest {
  title: string;
  description?: string;
  completed?: boolean;
  user: string;
  due_date?: string;
}

export interface UpdateTodoRequest {
  title?: string;
  description?: string;
  completed?: boolean;
  due_date?: string;
}
