import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ConfirmAccountPageComponent } from './pages/confirm-account-page/confirm-account-page.component';
import { BasePageComponent } from './pages/base-page/base-page.component';
import { ResetPasswordPageComponent } from './pages/reset-password-page/reset-password-page.component';

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
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
