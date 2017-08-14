import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AppInitService } from '../app.init.service';

@Component({
  selector: 'main-content',
  templateUrl: 'main-content.component.html'
})  
export class MainContentComponent {
  
  token: string;
  
  constructor(
    private appInitService: AppInitService,
    private router: Router
  ) {
    this.token = this.appInitService.token;
    this.appInitService.tokenMessage$.subscribe(token=>{
      this.token = token;
    })
  }

  navigateTo(link) {
    this.appInitService.token = this.token;
    this.router.navigate([link], {
      queryParams: {
        'token': this.token
      }
    })
  }
}