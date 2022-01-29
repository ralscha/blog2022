import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {TodoService} from '../todo.service';
import {MessagesService} from '../messages.service';
import {NgForm} from '@angular/forms';
import {HttpErrorResponse} from '@angular/common/http';
import {TodoPostResponse} from '../todo-post-response';
import {Todo} from '../todo';

@Component({
  selector: 'app-edit-page',
  templateUrl: './edit.page.html',
  styleUrls: ['./edit.page.scss'],
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

  private static displayFieldErrors(form: NgForm, fieldErrors: { [key: string]: string }): void {
    for (const [key, value] of Object.entries(fieldErrors)) {
      const comp = form.form.get(key);
      if (comp) {
        comp.setErrors({[value]: true});
      }
    }
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

    this.todoService.save(this.selectedTodo).subscribe({
      next: () => this.handleSuccessResponse(),
      error: this.handleErrorResponse(todoForm)
    });
  }

  dateChanged(event: any) {
    this.dueDate = event.detail.value;
    this.showCalendar = false;
  }

  private handleSuccessResponse(): void {
    this.messagesService.showSuccessToast('Todo successfully saved', 500);
    this.router.navigate(['/todo']);
  }

  private handleErrorResponse(form: NgForm) {
    return (errorResponse: HttpErrorResponse) => {
      const response: TodoPostResponse = errorResponse.error;
      if (response) {
        if (response.fieldErrors) {
          EditPage.displayFieldErrors(form, response.fieldErrors)
        } else {
          this.messagesService.showErrorToast('Saving Todo failed');
        }
      }
    };
  }

}
