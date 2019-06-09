import {EventEmitter, Injectable, Output} from '@angular/core';
import {User} from "../domain/user";
import {EMPTY, Observable, of, throwError} from "rxjs";
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {catchError, tap} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  @Output() userLogged = new EventEmitter<User>();
  @Output() userLoggedOut = new EventEmitter<User>();

  constructor(private http: HttpClient) {
  }

  private loggedUserKey = "loggedUser";

  public login(user: User): Observable<User> {
    return this.http.post<User>(environment.baseUrl + "/user/login", user).pipe(
      catchError((err) => {
        console.log("Login error");
        return throwError(err);
      }),
      tap((loggedUser: User) => {
        window.localStorage.setItem(this.loggedUserKey, JSON.stringify(loggedUser));
        this.userLogged.emit(loggedUser);
        return of(loggedUser);
      })
    );
  }

  public signUp(user: User): Observable<User> {
    return this.http.post<User>(environment.baseUrl + "/user/signup", user).pipe(
      catchError(
        (err) => {
          console.log("Signup error");
          return throwError(err);
        }
      ));
  }

  public logOut(): void {
    let user = <User>JSON.parse(window.localStorage.getItem(this.loggedUserKey));
    window.localStorage.removeItem(this.loggedUserKey);
    this.userLoggedOut.emit(user); //Emit who logged out
  }

  public getLoggedUser(): Observable<User> {
    let loggedUserJson = window.localStorage.getItem(this.loggedUserKey);
    if (loggedUserJson) {
      let parsedJson = <User>JSON.parse(loggedUserJson);
      return of(new User(parsedJson.username, parsedJson.password, parsedJson.isAdmin));
    } else {
      return EMPTY;
    }
  }
}
