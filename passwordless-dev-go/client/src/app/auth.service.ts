import {inject, Injectable} from '@angular/core';
import {Observable, of} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {catchError, map, tap} from 'rxjs/operators';
import {environment} from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly httpClient = inject(HttpClient);

  private loggedIn = false;

  isAuthenticated(): Observable<boolean> {
    return this.httpClient.post<void>(`${environment.API_URL}/authenticate`, null).pipe(
      tap(() => this.loggedIn = true),
      map(() => true),
      catchError(() => {
        this.loggedIn = false;
        return of(false);
      })
    );
  }

  logout(): Observable<void> {
    return this.httpClient.post<void>(`${environment.API_URL}/logout`, null).pipe(tap(() => this.loggedIn = false));
  }

  isLoggedIn(): boolean {
    return this.loggedIn;
  }

}



