import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { User } from 'src/app/domains/entity/user-profile.model'
import { Role } from 'src/app/domains/enums/roles.enum';
import { FormUpdateUserProfileComponent } from 'src/app/forms/form-update-user-profile/form-update-user-profile.component';
import { ModalService } from 'src/app/services/modal/modal.service';
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
    private modalService:ModalService,

    
  ) {}

  ngOnInit(){
    this.getUserID()
    console.log(this.profileID)
    this.profileService.getUserById(this.profileID as string).subscribe((data) => {
      console.log(2)
    this.user = data.data  
    })
    if (this.user === undefined){
      this.user = {
        id: "id ssl mock",
        firstName: "ime ssl mock",
        lastName: "prezime ssl mock",
        email: "mail ssl mock",
        residence: "rezidencija ssl mock",
        role: Role.Guest,
        username: "username ssl mock",
        age: 25,

      }
    } 
  }

  updateClick() {
    
    this.callUpdateProfile();
  }

  deleteClick(){
    this.callDeleteProfile();
  }

  callUpdateProfile() {
    this.modalService.open(FormUpdateUserProfileComponent, 'Update Profile', {"user": this.user});
  }

  callDeleteProfile(){
    this.profileService.delete(this.profileID as string)
  }

  getUserID() {
    this.route.paramMap.subscribe((params) => {
      this.profileID = String(params.get('id'));
    });
  }

  getUserById() {
    this.profileService.getUserById(this.profileID as string).subscribe((data) => {
      this.user = data.data;
    });
  }


}
