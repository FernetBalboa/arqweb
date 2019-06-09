import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {LoginService} from "../service/login.service";
import {User} from "../domain/user";
import {Router} from "@angular/router";

@Component({
  selector: 'app-navigation-toolbar',
  templateUrl: './navigation-toolbar.component.html',
  styleUrls: ['./navigation-toolbar.component.sass']
})
export class NavigationToolbarComponent implements OnInit {
  appTitle = 'Map Ui';
  loggedUser: User;
  @Output() sidenavToggled = new EventEmitter();

  constructor(private loginService: LoginService, private router: Router) {

  }

  ngOnInit() {
    //Remember logged user
    this.loginService.getLoggedUser().subscribe(
      (loggedUser: User) => {
        this.loggedUser = loggedUser;
      },
    );

    //Subscribe to new logins
    this.loginService.userLogged.subscribe(
      (loggedUser: User) => {
        this.loggedUser = loggedUser;
      }
    );

    //Subscribe to logouts
    this.loginService.userLoggedOut.subscribe(
      () => {
        this.loggedUser = null;
      }
    );

  }

  private logOut(): void {
    this.loginService.logOut();
    this.router.navigateByUrl("/");
  }

  //TODO: delete
  private logAdmin() {
    this.loginService.signUp(new User("Fernet", "pass", true)).subscribe(
      () => {
        this.loginService.login(new User("Fernet", "pass", null)).subscribe(
          () => this.router.navigateByUrl("/")
        );
      }
    );
  }

  private toggleSidenav() {
    this.sidenavToggled.emit();
  }
}
