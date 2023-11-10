import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SimpleProfileMenuComponent } from './simple-profile-menu.component';

describe('SimpleProfileMenuComponent', () => {
  let component: SimpleProfileMenuComponent;
  let fixture: ComponentFixture<SimpleProfileMenuComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [SimpleProfileMenuComponent]
    });
    fixture = TestBed.createComponent(SimpleProfileMenuComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
