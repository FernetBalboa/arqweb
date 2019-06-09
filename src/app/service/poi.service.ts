import {Injectable} from '@angular/core';
import {Image, PointOfInterest} from "../domain/point-of-interest";
import {HttpClient, HttpParams} from "@angular/common/http";
import {Feature} from "geojson";
import {environment} from "../../environments/environment";
import {LatLngBounds} from "leaflet";
import {Observable, of} from "rxjs";
import {Category} from "../domain/category";
import {CategorySuggestion, SuggestionStatus} from "../domain/category-suggestion";
import {map, take} from "rxjs/operators";

@Injectable({
  providedIn: 'root'
})
export class PoiService {
  private http: HttpClient;

  constructor(http: HttpClient) {
    this.http = http;
  }

  //Saves the POI in the backend
  public savePOI(POI: PointOfInterest): Observable<PointOfInterest> {
    let POIGeoJSON = PoiService.POIToGeoJSON(POI);
    return this.http.post<PointOfInterest>(environment.baseUrl + "/poi", POIGeoJSON);
  }

  //Converts a POI to a GeoJSON
  private static POIToGeoJSON(POI: PointOfInterest): Feature {
    return {
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [POI.lat, POI.long]
      },
      properties: {
        "title": POI.title,
        "category": POI.category.name,
        "description": POI.description,
        "picture": {
          "data": POI.picture.data,
          "name": POI.picture.name,
          "contentType": POI.picture.contentType
        },
        "hidden": POI.hidden
      }
    };

  }

  public Search(title?: string, categoryName?: string, markerLimit?: number, bounds?: LatLngBounds,
         showHiddenCategories?: boolean): Observable<PointOfInterest[]> {

    let queryParams = new HttpParams();

    if (title && !(title === "Any")) {
      queryParams = queryParams.append("title", title);
    }

    if (markerLimit) {
      queryParams = queryParams.append("limit", String(markerLimit));
    }

    if (categoryName && !(categoryName === "Any")) {
      queryParams = queryParams.append("category", categoryName.toLowerCase())
    }

    if (bounds) {
      let westLong = bounds.getWest();
      let eastLong = bounds.getEast();
      let northLat = bounds.getNorth();
      let southLat = bounds.getSouth();
      queryParams = queryParams.append("west_long", String(westLong));
      queryParams = queryParams.append("east_long", String(eastLong));
      queryParams = queryParams.append("north_lat", String(northLat));
      queryParams = queryParams.append("south_lat", String(southLat));
      queryParams = queryParams.append("bound", String(true));
    }

    if (showHiddenCategories) {
      queryParams = queryParams.append("hidden", String(true));
    }

    return this.http.get<PointOfInterest[]>(environment.baseUrl + "/poi/search",
      {params: queryParams});

  }

  public updatePOI(POI: PointOfInterest): Observable<PointOfInterest> {
    let poiBody = {
      id: POI.id,
      title: POI.title,
      category: POI.category.name, //Send name instead of object
      description: POI.description,
      picture: {
        data: POI.picture.data,
        name: POI.picture.name,
        contentType: POI.picture.contentType
      },
      hidden: POI.hidden,
      lat: POI.lat,
      long: POI.long
    };
    return this.http.put<PointOfInterest>(environment.baseUrl +  `/poi/${POI.id}`, poiBody);
  }

  public getCategories(): Observable<Category[]> {
    return this.http.get<Category[]>(environment.baseUrl + "/categories")
  }

  public getCategory(name: string): Observable<Category[]> {
    let queryParams = new HttpParams();
    queryParams = queryParams.append("name", name);
    return this.http.get<Category[]>(
      environment.baseUrl + `/categories/search`, {params: queryParams});
  }

  public updateCategory(category: Category): Observable<Category> {
    return this.http.put<Category>(environment.baseUrl + `/category/${category.id}`, category)
  }

  deleteCategory(categoryName: string): Observable<Category> {
    return this.http.delete<Category>(environment.baseUrl + `/category/${categoryName}`)
  }

  sendCategorySuggestion(categorySuggestion: CategorySuggestion): Observable<CategorySuggestion> {
    let suggestion = {
      name: categorySuggestion.category.name,
      hidden: categorySuggestion.category.hidden,
      iconClass: categorySuggestion.category.iconClass,
      status: SuggestionStatus[categorySuggestion.status]
    };
    return this.http.post<CategorySuggestion>(environment.baseUrl + "/suggestions/categories/new", suggestion);
  }

  getPendingCategorySuggestions(): Observable<CategorySuggestion[]> {
    return this.http.get(environment.baseUrl + "/suggestions/categories/new").pipe(
      (take (10)),
      map((suggestionDtoList: any[]) =>{
        let suggestions: CategorySuggestion[] = [];
        suggestionDtoList.forEach((suggestionDto) => {
          let category: Category = new Category(null, suggestionDto.name, suggestionDto.hidden, suggestionDto.iconClass);
          let newSuggestion: CategorySuggestion = new CategorySuggestion(category, suggestionDto.status);
          newSuggestion.id = suggestionDto.id;
          suggestions.push(newSuggestion);
        });
        return suggestions;
      })
    );
  }

  approveSuggestion(suggestion: CategorySuggestion): Observable<any> {
    return this.http.put(environment.baseUrl + `/suggestions/categories/new/${suggestion.id}/approve`, null);
  }

  rejectSuggestion(suggestion: CategorySuggestion): Observable<any> {
    return this.http.put(environment.baseUrl + `/suggestions/categories/new/${suggestion.id}/reject`, null);
  }


}
