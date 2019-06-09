import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CategorySuggestionsComponent } from './category-suggestions.component';

describe('CategorySuggestionsComponent', () => {
  let component: CategorySuggestionsComponent;
  let fixture: ComponentFixture<CategorySuggestionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CategorySuggestionsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CategorySuggestionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
