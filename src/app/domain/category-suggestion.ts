import {Category} from "./category";

export class CategorySuggestion {
  id: string;
  category: Category;
  status: SuggestionStatus;


  constructor(category: Category, status: SuggestionStatus) {
    this.category = category;
    this.status = status;
  }
}

export enum SuggestionStatus {
  APPROVED,
  WAITING_FOR_APPROVAL,
  REJECTED
}
