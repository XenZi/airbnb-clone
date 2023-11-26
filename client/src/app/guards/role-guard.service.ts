import { Injectable } from '@angular/core';
import { UserService } from '../services/user/user.service';
import {
  ActivatedRouteSnapshot,
  Router,
  RouterStateSnapshot,
  UrlTree,
} from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class RoleGuardService {
  constructor(private userService: UserService, private router: Router) {}

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): boolean | UrlTree {
    const allowedRoles = route.data['allowedRoles'] as Array<string>;
    const userRole = this.userService.getLoggedUser()?.role as string;
    if (allowedRoles.includes(userRole)) {
      return true;
    } else {
      return this.router.parseUrl('/');
    }
  }
}
