import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {PoiService} from "../service/poi.service";
import {CategorySuggestion, SuggestionStatus} from "../domain/category-suggestion";
import {Category} from "../domain/category";
import {Location} from "@angular/common";

@Component({
  selector: 'app-new-category-suggestion-form',
  templateUrl: './new-category-suggestion-form.component.html',
  styleUrls: ['./new-category-suggestion-form.component.sass']
})
export class NewCategorySuggestionFormComponent implements OnInit {

  private categorySuggestionFrom: FormGroup;

  constructor(private poiService: PoiService, private location: Location) { }

  ngOnInit() {
    this.categorySuggestionFrom = new FormGroup(
      {
        name: new FormControl('', [Validators.required, Validators.maxLength(30)]),
        hidden: new FormControl(false, [Validators.required, Validators.maxLength(30)]),
        iconClass: new FormControl('', [Validators.required, Validators.maxLength(30)])
      }
    )
  }

  onIconPickerSelect(newIcon: string): void {
    this.categorySuggestionFrom.get("iconClass").setValue(newIcon);
  }

  suggestNewCategory(categorySuggestionFrom: any): void {
    let newCategorySuggestion = new CategorySuggestion(
      new Category(null, categorySuggestionFrom.name, categorySuggestionFrom.hidden,
        categorySuggestionFrom.iconClass),
      SuggestionStatus.WAITING_FOR_APPROVAL
    );

    this.poiService.sendCategorySuggestion(newCategorySuggestion).subscribe(()=>{});
  }

  cancelNewCategorySuggestion(): void {
    this.location.back();
  }
}
