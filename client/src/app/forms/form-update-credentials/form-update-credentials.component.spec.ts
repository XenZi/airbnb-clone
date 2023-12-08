import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FormUpdateCredentialsComponent } from './form-update-credentials.component';

describe('FormUpdateCredentialsComponent', () => {
  let component: FormUpdateCredentialsComponent;
  let fixture: ComponentFixture<FormUpdateCredentialsComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FormUpdateCredentialsComponent]
    });
    fixture = TestBed.createComponent(FormUpdateCredentialsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
