import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { NewCategorySuggestionFormComponent } from './new-category-suggestion-form.component';

describe('NewCategorySuggestionFormComponent', () => {
  let component: NewCategorySuggestionFormComponent;
  let fixture: ComponentFixture<NewCategorySuggestionFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ NewCategorySuggestionFormComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(NewCategorySuggestionFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
