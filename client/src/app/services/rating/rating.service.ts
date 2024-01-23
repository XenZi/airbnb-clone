import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { UserService } from '../user/user.service';

@Injectable({
  providedIn: 'root',
})
export class RatingService {
  constructor(private http: HttpClient, private userService: UserService) {}
}
