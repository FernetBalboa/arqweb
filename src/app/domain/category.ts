export class Category {
  id: string;
  name: string;
  hidden: boolean;
  iconClass: string; //FontAwesome 5 icons available


  constructor(id: string, name: string, hidden: boolean, iconClass: string) {
    this.id = id;
    this.name = name;
    this.hidden = hidden;
    this.iconClass = iconClass;
  }
}
