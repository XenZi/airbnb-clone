import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-confirm-account-page',
  templateUrl: './confirm-account-page.component.html',
  styleUrls: ['./confirm-account-page.component.scss'],
})
export class ConfirmAccountPageComponent {
  token!: string;
  constructor(private route: ActivatedRoute) {}

  ngOnInit() {
    console.log('etst');
    this.route.paramMap.subscribe((params) => {
      this.token = String(params.get('token'));
    });
    console.log(this.token);
  }
}
