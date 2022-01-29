export interface Todo {
  id: string;
  priority: 'high' | 'normal' | 'low';
  dueDate?: string;
  description: string;
}
