import { Component, OnInit } from '@angular/core';
import {CategorySuggestion} from "../domain/category-suggestion";
import {PoiService} from "../service/poi.service";
import {MatSelectionListChange} from "@angular/material";

@Component({
  selector: 'app-category-suggestions',
  templateUrl: './category-suggestions.component.html',
  styleUrls: ['./category-suggestions.component.sass']
})
export class CategorySuggestionsComponent implements OnInit {

  suggestions: CategorySuggestion[];
  selectedSuggestions: CategorySuggestion[];

  constructor(private poiService: PoiService) { }

  ngOnInit() {
    this.getPendingSuggestions();

    this.selectedSuggestions = [];

  }

  private getPendingSuggestions() {
    this.poiService.getPendingCategorySuggestions().subscribe(
      (pendingSuggestions: CategorySuggestion[]) => {
        this.suggestions = pendingSuggestions;
      }
    );
  }

  approve() {
    this.selectedSuggestions.forEach(
      (suggestion: CategorySuggestion) => {
        this.poiService.approveSuggestion(suggestion).subscribe(
          () => this.getPendingSuggestions()
        );

      }
    );

  }

  reject() {
    this.selectedSuggestions.forEach(
      (suggestion: CategorySuggestion) => {
        this.poiService.rejectSuggestion(suggestion).subscribe(
          () => this.getPendingSuggestions()

        );
      }
    );
    // this.getPendingSuggestions();

  }


  getOptionClass(suggestion: CategorySuggestion): string {
    if (suggestion.category.hidden) {
      return "hidden-category"
    }
    return "";
  }
}
