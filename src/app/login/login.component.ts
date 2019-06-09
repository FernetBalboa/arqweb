import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {LoginService} from "../service/login.service";
import {User} from "../domain/user";
import {Router} from "@angular/router";

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.sass']
})
export class LoginComponent implements OnInit {

  loginControl: FormGroup;
  signUpControl: FormGroup;
  selectedTabIndex: number; //0 for login, 1 for signUp

  constructor(private loginService: LoginService, private router: Router) { }

  ngOnInit() {
    this.loginControl = new FormGroup(
      {
        username: new FormControl('', [Validators.required, Validators.maxLength(30)]),
        password: new FormControl('', [Validators.required, Validators.maxLength(30)])
      }
    );

    this.signUpControl = new FormGroup(
      {
        username: new FormControl('', [Validators.required, Validators.maxLength(30)]),
        password: new FormControl('', [Validators.required, Validators.maxLength(30)])
      }
    );
    this.selectedTabIndex = 0; //Set login tab

  }

  login(loginData: any): void {
    this.loginService.login(new User(loginData.username, loginData.password, false))
      .subscribe(
      //Go back to map view
      (loggedUser: User) => this.router.navigateByUrl("/map")
    );
  }

  signUp(signUpData: any): void {
    this.loginService.signUp(new User(signUpData.username, signUpData.password, false)).subscribe();
    this.selectedTabIndex = 0;
    //TODO: show user creation popup
    //Allow quick login after signUp and reset signUp form
    this.loginControl.get("username").setValue(signUpData.username);
    this.loginControl.get("password").setValue(signUpData.password);
    this.signUpControl.get("username").setValue("");
    this.signUpControl.get("password").setValue("");
  }

}
