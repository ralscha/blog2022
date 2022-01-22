import Dexie from 'dexie';

export class TodoDb extends Dexie {
  todos!: Dexie.Table<Todo, string>;

  constructor() {
    super('tododb');
    this.version(1).stores({
      todos: 'id'
    });
  }
}

export interface Todo {
  id: string;
  priority: 'high' | 'normal' | 'low';
  dueDate?: string;
  description: string;
}


