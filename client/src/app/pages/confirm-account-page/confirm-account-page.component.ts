import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AuthService } from 'src/app/services/auth-service/auth.service';

@Component({
  selector: 'app-confirm-account-page',
  templateUrl: './confirm-account-page.component.html',
  styleUrls: ['./confirm-account-page.component.scss'],
})
export class ConfirmAccountPageComponent {
  token!: string;
  constructor(
    private route: ActivatedRoute,
    private authService: AuthService
  ) {}

  ngOnInit() {
    this.route.paramMap.subscribe((params) => {
      this.token = String(params.get('token'));
      this.authService.confirmAccount(this.token + '=');
    });
  }
}
