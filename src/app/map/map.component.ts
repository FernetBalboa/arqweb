import {ChangeDetectorRef, Component, OnInit, ViewChild} from '@angular/core';
import * as L from "leaflet";
import {icon, LatLng, latLng, Layer, LeafletMouseEvent, Map, marker, Marker, tileLayer} from "leaflet";
import {PointOfInterest} from "../domain/point-of-interest";
import {AddMarkerFormComponent} from "../add-marker-form/add-marker-form.component";
import {PoiService} from "../service/poi.service";

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrls: ['./map.component.sass']
})
export class MapComponent implements OnInit {

  wMaps = tileLayer('http://maps.wikimedia.org/osm-intl/{z}/{x}/{y}.png', {
    detectRetina: true,
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  });

  streetMaps = tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    detectRetina: true,
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
  });

  initialSetup: object;

  //Control base layers and overlays
  layersControl: object;

  layers: Layer[];

  addingNewPOI = false;

  newMarker: Marker;

  private addMarkerComponent: AddMarkerFormComponent;
  private map: Map;
  private static setMarkerIconUrl = 'leaflet/marker-icon.png';
  private static setMarkerShadowUrl = 'leaflet/marker-shadow.png';


  //Use @ViewChild to inject child component. Using setter
  //because of the *ngIf in the html template.
  @ViewChild('addMarkerComponent') set markerComponent(addMarkerComponent: AddMarkerFormComponent) {
    this.addMarkerComponent = addMarkerComponent;
    this.cdRef.detectChanges();
  }

  constructor(private cdRef:ChangeDetectorRef, private poiService: PoiService) {
  }

  ngOnInit() {
    this.initialSetup = {
      layers: [
        this.streetMaps
      ],
      zoom: 15,
      center: latLng([ -34.581199, -58.421058 ])
    };

    this.layersControl= {
      baseLayers: {
        "Street": this.streetMaps,
        "Wikimedia": this.wMaps
      },
      overlays: {}
    };

    this.layers = [];
  }

  onMapReady(map: Map) {
    this.map = map;
    this.poiService.Search(null, null, 20,
      this.map.getBounds(), false).subscribe(
      (POIs: PointOfInterest[]) => {
        let layers = [];
        POIs.forEach((POI: PointOfInterest) => {
          layers.push(MapComponent.POIToMarker(POI))
        });
        this.layers = layers;
      }
    );
  }

  showSavePOIPopup(latLong: LatLng) {
    if (this.addingNewPOI) {
      //Remove last marker
      this.layers.pop();
    }

    let newMarker = marker(latLong, {
      icon: icon({
        iconSize: [ 25, 41 ],
        iconAnchor: [ 13, 41 ],
        iconUrl: 'https://cdn.rawgit.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
        shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/0.7.7/images/marker-shadow.png',
      }),
      clickable: true,
      draggable: true,
      riseOnHover: true
    });
    this.layers.push(newMarker);
    this.newMarker = newMarker;

    //Zoom on marker when clicked
    this.newMarker.on("click", function (me: LeafletMouseEvent) {
      let latLng = [me.latlng];
      // noinspection JSPotentiallyInvalidUsageOfClassThis
      this.map.fitBounds(L.latLngBounds(latLng), {maxZoom: 15})
    },{map: this.map});

    this.addingNewPOI = true;
    //Subscribes add-marker component to the new marker to track coordinates
    //SetTimeout is used to give time for the @ViewChild to set the component
    //after the *ngIf activation
    setTimeout(() => this.addMarkerComponent.resetPOIForm(), 1);
    //Set map on top of screen to give more space to the add POI form
    setTimeout(() => document.getElementById("map").scrollIntoView(true), 100);

  }

  hideNewPOIForm(removeLastLayer = true) {
    this.addingNewPOI = false;
    if (removeLastLayer) {
      this.layers.pop();
    }
  }

  onPOIAdd(savedPOI: PointOfInterest) {
    this.hideNewPOIForm(false);
    this.newMarker.bindPopup(<string>MapComponent.markerPopupHtml(savedPOI), {maxWidth: 700, className: 'popup'});
    this.newMarker.dragging.disable();
    let newIcon = icon({
      iconUrl: MapComponent.setMarkerIconUrl,
      shadowUrl: MapComponent.setMarkerShadowUrl
    });
    this.newMarker.setIcon(newIcon);

  }

  static markerPopupHtml(POI: PointOfInterest): string {
    let html: string = ` 
<head>
    <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.0/css/all.css" integrity="sha384-lZN37f5QGtY3VHgisS14W3ExzMWZxybE1SJSEsQp9S+oqd12jhcu+A56Ebc1zFSJ" crossorigin="anonymous">
</head>
<h2>${this.toTitleCase(POI.title)}</h2>
<h3>${this.toTitleCase(POI.category.name)} <i class="${POI.category.iconClass}"></i></h3>  
<p>Description: ${POI.description}</p>
`;
    if (POI.picture.data) {
      html = html + `<img src="data:${POI.picture.contentType};base64, ${POI.picture.data}" style="width: 100px" style="height: 45px" "/>`
    }
    return html;
  }

  private static toTitleCase(str: string): string {
    if (str) {
      return str.replace(
        /\w\S*/g,
        function (txt) {
          return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
        }
      );
    } else {
      return "";
    }
  }

  /**
   * Filters the layers in the map.
   *
   * @param poiFilters object with the values from {@link POIFilterComponent#filtersControl}.
   */
  filterPOIs(poiFilters: any) {
    let category = poiFilters.category;
    let title = poiFilters.title;
    let markerLimit = 20; //Arbitrary marker limit
    let bounds = this.map.getBounds();

    this.layers = [];
    this.poiService.Search(title, category, markerLimit, bounds, false).subscribe(
      (searchResult: PointOfInterest[]) => searchResult.forEach(
        (POI: PointOfInterest) => {
          this.layers.push(MapComponent.POIToMarker(POI))
        }, this
      )
    );
  }


  public static POIToMarker(POI: PointOfInterest): Marker {
    let marker = new Marker([POI.lat, POI.long],
      {
        draggable: false,
        clickable: true,
        riseOnHover: true,
        icon: icon({
          iconSize: [ 25, 41 ],
          iconAnchor: [ 13, 41 ],
          iconUrl: this.setMarkerIconUrl,
          shadowUrl: this.setMarkerShadowUrl,
        }),
        title: POI.title
      });

    marker.bindPopup(<string>MapComponent.markerPopupHtml(POI), {maxWidth: 700, className: 'popup'});

    return marker
  }

}
