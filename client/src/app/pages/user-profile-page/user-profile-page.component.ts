import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { User } from 'src/app/domains/entity/user-profile.model'
import { ProfileService } from 'src/app/services/profile/profile.service';


@Component({
  selector: 'app-user-profile-page',
  templateUrl: './user-profile-page.component.html',
  styleUrls: ['./user-profile-page.component.scss']
})
export class UserProfilePageComponent {
  profileID: string | undefined
  user!: User

  constructor(
    private route: ActivatedRoute,
    private profileService: ProfileService,
  ) {}

  ngOnInit(){
    console.log("zapoceo")
    this.getUserID()
    console.log(this.profileID)
    this.getUserById()
    
    
  }
  getUserID() {
    this.route.paramMap.subscribe((params) => {
      this.profileID = String(params.get('id'));
    });
  }

  getUserById() {
    this.profileService.getUserById(this.profileID as string).subscribe((data) => {
      console.log(data)
      console.log("dva")
      this.user = data.data;
    });
  }


}
