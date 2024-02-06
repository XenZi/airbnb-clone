import { Injectable } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';
import { UnloadService } from '../unload/unload.service';

@Injectable({
  providedIn: 'root'
})
export class PreviouseRouteService {

  private previousUrl!: string;
  private currentUrl: string;

  constructor(private router: Router, private unloadService: UnloadService) {
    this.currentUrl = this.router.url;
    router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {        
        this.previousUrl = this.currentUrl;
        this.currentUrl = event.url;
        this.unloadService.setPreviousUrl(this.previousUrl);
      };
    });
  }

  public getPreviousUrl() {
    return this.previousUrl;
  }   

}
