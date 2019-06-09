import {Feature} from "geojson";
import {Category} from "./category";

export class PointOfInterest {

  id: string;
  title: string;
  category: Category;
  hidden: boolean;
  description: string;
  lat: number;
  long: number;
  picture: Image; //Base64 image

}

export class Image {
  data: string; //Base64
  name: string; //Filename
  contentType: string;


  constructor(data: string, name: string, contentType: string) {
    this.data = data;
    this.name = name;
    this.contentType = contentType;
  }
}
