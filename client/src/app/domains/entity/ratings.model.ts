export interface AccommodationRate {
  accommodationID?: string;
  rate?: number;
  guest?: Guest;
  createdAt?: string;
  avgRating?: number;
  hostEmail?: string;
  hostID?: string;
}

export interface HostRate {
  host?: Host;
  rate?: number;
  guest?: Guest;
  createdAt?: string;
  avgRating?: number;
  hostEmail?: string;
  hostID?: string;
}

export interface Guest {
  id: string;
  email: string;
  username: string;
}

export interface Host {
  id: string;
  email: string;
  username: string;
}
