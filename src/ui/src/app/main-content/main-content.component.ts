import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AppInitService } from '../app.init.service';

@Component({
  selector: 'main-content',
  templateUrl: 'main-content.component.html'
})  
export class MainContentComponent {
  
  token: string;
  
  isSignIn: boolean = true;
  hasSignedIn: boolean = false;
  searchContent: string = '';

  constructor(
    private appInitService: AppInitService,
    private router: Router,
    private route: ActivatedRoute
  ) {
    if(this.appInitService.currentUser) {
      this.isSignIn = false;
      this.hasSignedIn = true;
    }
    this.token = this.appInitService.token;
    this.appInitService.tokenMessage$.subscribe(token=>this.token = token);
    this.route.queryParamMap.subscribe(params=>this.searchContent = params.get("q"));
  }

  get isSystemAdmin(): boolean {
    if(this.appInitService.currentUser) {
      return this.appInitService.currentUser["user_system_admin"] == 1;
    }
    return false;
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