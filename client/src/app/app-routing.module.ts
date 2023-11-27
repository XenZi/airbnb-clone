import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ConfirmAccountPageComponent } from './pages/confirm-account-page/confirm-account-page.component';
import { BasePageComponent } from './pages/base-page/base-page.component';
import { ResetPasswordPageComponent } from './pages/reset-password-page/reset-password-page.component';
import { RoleBasedPageComponent } from './pages/role-based-page/role-based-page.component';
import { RoleGuardService } from './guards/role-guard.service';
import { ReservationFormComponent } from './forms/form-create-reservation/form-create-reservation.component';

const routes: Routes = [
  {
    path: '',
    component: BasePageComponent,
  },
  {
    path: 'confirm-account/:token',
    component: ConfirmAccountPageComponent,
  },
  {
    path: 'reset-password/:token',
    component: ResetPasswordPageComponent,
  },
  {
    path: 'role-based-page',
    component: RoleBasedPageComponent,
    canActivate: [RoleGuardService],
    data: {
      allowedRoles: ['HOST', 'Host'],
    },
  },
  { path: 'create-reservation', component: ReservationFormComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
