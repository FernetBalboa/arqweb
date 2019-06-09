import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EditPoiComponent } from './edit-poi.component';

describe('EditPoiComponent', () => {
  let component: EditPoiComponent;
  let fixture: ComponentFixture<EditPoiComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EditPoiComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EditPoiComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
