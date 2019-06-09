import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddMarkerFormComponent } from './add-marker-form.component';

describe('AddMarkerFormComponent', () => {
  let component: AddMarkerFormComponent;
  let fixture: ComponentFixture<AddMarkerFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AddMarkerFormComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddMarkerFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
