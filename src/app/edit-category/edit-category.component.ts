import { Component, OnInit } from '@angular/core';
import {PoiService} from "../service/poi.service";
import {Category} from "../domain/category";

@Component({
  selector: 'app-edit-category',
  templateUrl: './edit-category.component.html',
  styleUrls: ['./edit-category.component.sass']
})
export class EditCategoryComponent implements OnInit {

  searchCategoryName: string;
  category: Category;
  fallbackIcon: string ="fas fa-question-sign";
  icon: string;

  constructor(private poiService: PoiService) { }

  ngOnInit() {
    this.searchCategoryName = '';
  }

  private getCategory(): void {
    this.poiService.getCategory(this.searchCategoryName).subscribe(
      (category: Category[]) => {
        this.category = category[0];
      }
    );
    this.icon = this.category.iconClass;
  }

  private updateCategory(): void {
    this.poiService.updateCategory(this.category).subscribe(
      (newCategory: Category) => {
        this.category = newCategory;
      }
    );
  }

  private cancelUpdateCategory(): void {
    this.category = null;
  }

  private deleteCategory(): void {
    this.poiService.deleteCategory(this.category.name).subscribe(
      () => {
        this.category = null;
      }
    );
  }

  onIconPickerSelect(newIcon: string) {
    this.category.iconClass = newIcon
  }
}
