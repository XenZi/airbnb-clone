import { Injectable } from '@angular/core';
import { LocalStorageService } from '../localstorage/local-storage.service';
import { HttpClient } from '@angular/common/http';
import { apiURL } from 'src/app/domains/constants';
import { ModalService } from '../modal/modal.service';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  constructor(
    private localStorageService: LocalStorageService,
    private http: HttpClient,
    private modalService: ModalService
  ) {}

  login(email: string, password: string): void {
    this.http
      .post(
        `${apiURL}/login`,
        {
          email,
          password,
        },
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )
      .subscribe({
        next: (data: any) => {
          this.localStorageService.setItem('token', data.data?.Token);
          this.localStorageService.setItem(
            'user',
            JSON.stringify(data.data?.User)
          );
          this.modalService.close();
        },
        error: (err) => {
          this.modalService.close();
        },
      });
  }
}
