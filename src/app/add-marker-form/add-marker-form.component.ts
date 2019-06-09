import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";
import {Image, PointOfInterest} from "../domain/point-of-interest";
import {Marker} from "leaflet";
import {PoiService} from "../service/poi.service";
import {Category} from "../domain/category";
import {__await} from "tslib";
import {promise} from "selenium-webdriver";

@Component({
  selector: 'app-add-marker-form',
  templateUrl: './add-marker-form.component.html',
  styleUrls: ['./add-marker-form.component.sass']
})
export class AddMarkerFormComponent implements OnInit {

  @Input() marker: Marker;
  addPOIForm: FormGroup;
  newPOI: PointOfInterest;
  @Output() addedPOI: EventEmitter<any> = new EventEmitter();
  @Output() cancelPOI: EventEmitter<any> = new EventEmitter();
  private availableCategories: Category[];

  private poiService: PoiService;

  constructor(poiService: PoiService) {
    this.poiService = poiService;
  }

  ngOnInit() {
    this.addPOIForm = AddMarkerFormComponent.defaultForm(0, 0); //Placeholder
    this.poiService.getCategories().subscribe(
      (categories: Category[]) => {
        this.availableCategories = categories;
      }
    );
  }

  //Resets the form to track the newest marker
  public resetPOIForm() {
    //Subscribe to marker events to update form accordingly
    this.marker.on('add move drag dragEnd', function () {
      let newLatLng = this.formComponent.marker.getLatLng();
      this.formComponent.addPOIForm.setControl("lat", new FormControl(newLatLng.lat));
      this.formComponent.addPOIForm.setControl("long", new FormControl(newLatLng.lng));
    }, {formComponent: this});

    let latLng = this.marker.getLatLng();
    this.addPOIForm = AddMarkerFormComponent.defaultForm(latLng.lat, latLng.lng);

  }

  static defaultForm(lat: number, long: number): FormGroup {
    return new FormGroup({
      title: new FormControl('', [Validators.required, Validators.maxLength(60)]),
      category: new FormControl('', [Validators.required, Validators.maxLength(60)]),
      description: new FormControl(''),
      lat: new FormControl(lat),
      long: new FormControl(long),
      picture: new FormGroup({
        name: new FormControl(''),
        contentType: new FormControl(''),
        data: new FormControl('')
        }
      )
    })
  }

  onFileChange(event) {
    let reader = new FileReader();
    if(event.target.files && event.target.files.length > 0) {
      let file = event.target.files[0];
      reader.readAsDataURL(file);
      reader.onload = () => {
        this.addPOIForm.get('picture').setValue({
          name: file.name,
          contentType: file.type,
          data: (<string>reader.result).split(',')[1]
        })
      };
    }
  }

  //Argument has same fields as control group
  addPOI(newPoiForm: any) {
    let latLng = this.marker.getLatLng();
    let formPicture = newPoiForm.picture;
    let picture = new Image(formPicture.data, formPicture.name, formPicture.contentType);
    this.createPOI(newPoiForm.title, newPoiForm.category,
      newPoiForm.description, latLng.lat, latLng.lng, picture);
  }

  createPOI(title: string, category: string, description: string,
            lat: number, long: number, picture: Image) {
    this.poiService.getCategory(category).subscribe(
      (category: Category[]) => {
        let newPOI = new PointOfInterest();
        newPOI.title = title;
        newPOI.category = category[0];
        newPOI.description = description;
        newPOI.lat = lat;
        newPOI.long = long;
        newPOI.picture = picture;
        newPOI.hidden = false;
        this.newPOI = newPOI;
        this.saveNewPOI();

        //Inform map of save event
        this.addedPOI.emit(this.newPOI);
      }
    );
  }

  cancelNewPOI() {
    this.cancelPOI.emit(null);
  }

  saveNewPOI(){
    this.poiService.savePOI(this.newPOI).subscribe();
  }

}
