import { Injectable } from '@angular/core';
import { LocalStorageService } from '../localstorage/local-storage.service';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';

@Injectable({
  providedIn: 'root',
})
export class UserService {
  constructor(private localStorage: LocalStorageService) {}

  getLoggedUser(): UserAuth | null {
    let savedLocally = this.localStorage.getItem('user');
    if (savedLocally == null) {
      return null;
    }

    return JSON.parse(savedLocally) as UserAuth;
  }
}
