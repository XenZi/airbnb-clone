import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RoleBasedPageComponent } from './role-based-page.component';

describe('RoleBasedPageComponent', () => {
  let component: RoleBasedPageComponent;
  let fixture: ComponentFixture<RoleBasedPageComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [RoleBasedPageComponent]
    });
    fixture = TestBed.createComponent(RoleBasedPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
