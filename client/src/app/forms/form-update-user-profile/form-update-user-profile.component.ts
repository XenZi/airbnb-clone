import { Component, Input } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { timeout } from 'rxjs';
import { Accommodation } from 'src/app/domains/entity/accommodation-model';
import { UserAuth } from 'src/app/domains/entity/user-auth.model';
import { User } from 'src/app/domains/entity/user-profile.model';
import { Role } from 'src/app/domains/enums/roles.enum';
import { ToastNotificationType } from 'src/app/domains/enums/toast-notification-type.enum';
import { AccommodationsService } from 'src/app/services/accommodations-service/accommodations.service';
import { ProfileService } from 'src/app/services/profile/profile.service';
import { ToastService } from 'src/app/services/toast/toast.service';
import { UserService } from 'src/app/services/user/user.service';
import { formatErrors } from 'src/app/utils/formatter.utils';

@Component({
  selector: 'app-form-update-user-profile',
  templateUrl: './form-update-user-profile.component.html',
  styleUrls: ['./form-update-user-profile.component.scss']
})
export class FormUpdateUserProfileComponent {
  @Input() userID!:string

  updateProfileForm: FormGroup;
  loggedUser!: UserAuth | null
  user!: User
  errors: string = '';
  isCaptchaValidated: boolean = false;

  constructor(
    private profileService: ProfileService,
    private formBuilder: FormBuilder,
    private toastService: ToastService,
    private userService: UserService
  ){
    this.updateProfileForm = this.formBuilder.group({
      firstName: [''],
      lastName: [''],
      email: [''],
      residence: [''],
      username: [''],
      age: [''],
    })
  }
  ngOnInit(){

    this.getUserInfo();
    setTimeout(() => {
      this.updateProfileForm = this.formBuilder.group({
        firstName: [this.user.firstName, Validators.required],
        lastName: [this.user.lastName, Validators.required],
        email: [this.user.email, Validators.required],
        residence: [this.user.residence, Validators.required],
        username: [this.user.username, Validators.required],
        age: [this.user.age, Validators.required]
      });
    }, 300);

  }

  getUserInfo(){
   this.loggedUser = this.userService.getLoggedUser()
   this.profileService.getUserById(this.user.id as string).subscribe((data) => {
    this.user = data
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
    this.loggedUser ={
      id: "id ssl mock",
      email: "mail ssl mock",
      role: Role.Guest,
      username: "username ssl mock",
      confirmed: true
    }
  } 
  }

  onSubmit(e: Event){
    e.preventDefault()

    if (!this.updateProfileForm.valid){
      console.log("yee")
      Object.keys(this.updateProfileForm.controls).forEach((key) =>{
        const controlErrors = this.updateProfileForm.get(key)?.errors;
        if (controlErrors){
          this.errors += formatErrors(key);
        }
      });
      this.toastService.showToast(
        'Error',
        this.errors,
        ToastNotificationType.Error
      );
      this.errors = '';
      return;
    }
    this.profileService.update(
      this.user.id,
      this.updateProfileForm.value.firstName,
      this.updateProfileForm.value.lastName,
      this.updateProfileForm.value.email,
      this.updateProfileForm.value.residence,
      this.user.role,
      this.updateProfileForm.value.username,
      this.updateProfileForm.value.age,

    )
    



  }




}

