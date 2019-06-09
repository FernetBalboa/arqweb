import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {Observable} from 'rxjs';
import {map, startWith} from 'rxjs/operators';
import {Category} from "../domain/category";
import {PoiService} from "../service/poi.service";

/**
 * @title Filter autocomplete
 */
@Component({
  selector: 'app-poi-filter',
  templateUrl: './poi-filter.component.html',
  styleUrls: ['./poi-filter.component.sass'],
})
export class POIFilterComponent implements OnInit {
  filtersControl: FormGroup;
  availableCategories: Category[];
  filteredCategoryOptions: Observable<string[]>;
  @Output() filterChange: EventEmitter<string> = new EventEmitter();

  constructor(private poiService: PoiService) {

  }

  ngOnInit() {
    this.filtersControl = new FormGroup({
      title: new FormControl("Any", [Validators.maxLength(30)]),
      category: new FormControl("Any", [Validators.required, Validators.maxLength(30)])
    });

    this.filteredCategoryOptions = this.filtersControl.get("category").valueChanges
      .pipe(
        map(value => {return value.name}),
        startWith(''),
        // map(value => this._filter(value))
      );

    this.poiService.getCategories().subscribe(
      (categories: []) => {
        this.availableCategories = categories;
      }
    );
  }

  // private _filter(value: string): string[] {
  //   const filterValue = value.toLowerCase();
  //
  //   return this.availableCategories.filter(option => option.name.toLowerCase().includes(filterValue))
  //     .map(option => {return option.name});
  // }

  //Only emits category, but could be extended with other filters
  updateFilters(filtersControl: any) {
    this.filterChange.emit(filtersControl)
  }

}
