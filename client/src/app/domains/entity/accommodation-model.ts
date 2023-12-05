import { DateAvailability } from "./date-availability.model";

export interface Accommodation {
    id: string;
    name: string;
    userId:string,
    username:string,
    location: string;
    conveniences: string;
    minNumOfVisitors: number;
    maxNumOfVisitors:number;
    AvailableAccommodationDates:DateAvailability[];
  }