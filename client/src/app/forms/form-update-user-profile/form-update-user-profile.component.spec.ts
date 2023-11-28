import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormUpdateUserProfileComponent } from './form-update-user-profile.component';

describe('FormUpdateUserProfileComponent', () => {
  let component: FormUpdateUserProfileComponent;
  let fixture: ComponentFixture<FormUpdateUserProfileComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormUpdateUserProfileComponent]
    });
    fixture = TestBed.createComponent(FormUpdateUserProfileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
