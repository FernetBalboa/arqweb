import { Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  CanActivate,
  CanActivateChild,
  Router,
  RouterStateSnapshot,
  UrlTree
} from "@angular/router";
import {Observable} from "rxjs";
import {LoginService} from "./login.service";
import {User} from "../domain/user";

/**
 * Protects urls from accessing without permission.
 */
@Injectable({
  providedIn: 'root'
})
export class AuthGuardService implements CanActivate, CanActivateChild {

  private adminUrls = ["/administration/edit/poi", "/administration/edit/category",
    "/administration/suggestion/category/new"];

  constructor(private loginService: LoginService, private router: Router) { }

  public canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot):
    Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {

    if (this.adminUrls.includes(state.url)) {
      let isAdmin: boolean = false;
      this.loginService.getLoggedUser().subscribe(
        (loggedUser: User) => {
          isAdmin = loggedUser.isAdmin;
        }
      );
      if (!isAdmin) {
        return this.router.parseUrl("/"); //Redirect to home
      } else {
        return true; //Allow entry to backoffice
      }
    } else {
      return true;
    }

  }

  public canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot):
    Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {

    if (this.adminUrls.includes(state.url)) {
      let isAdmin: boolean = false;
      this.loginService.getLoggedUser().subscribe(
        (loggedUser: User) => {
          isAdmin = loggedUser.isAdmin;
        }
      );
      if (!isAdmin) {
        return this.router.parseUrl("/"); //Redirect to home
      } else {
        return true; //Allow entry to backoffice
      }
    } else {
      return true;
    }

  }


}
