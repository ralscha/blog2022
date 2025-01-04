import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {NgForm} from '@angular/forms';
import {Todo} from '../todo-db';

@Component({
    selector: 'app-edit-page',
    templateUrl: './edit.page.html',
    styleUrls: ['./edit.page.scss'],
    standalone: false
})
export class EditPage implements OnInit {

  selectedTodo!: Todo;
  showCalendar = false;
  dueDate?: string

  constructor(private readonly route: ActivatedRoute,
              private readonly router: Router,
              private readonly messagesService: MessagesService,
              private readonly todoService: TodoService) {
  }

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
