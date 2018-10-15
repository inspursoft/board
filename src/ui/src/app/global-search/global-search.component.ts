import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

import { AppInitService, AppTokenService } from '../app.init.service';
import { GlobalSearchService } from './global-search.service';

@Component({
  selector: 'global-search',
  templateUrl: 'global-search.component.html'
})
export class GlobalSearchComponent implements OnInit {

  token: string;
  
  hasSignedIn: boolean;
  globalSearch: {[key: string]: any} = {};
  
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private appInitService: AppInitService,
    private appTokenService: AppTokenService,
    private globalSearchService: GlobalSearchService,
  ) {}

  ngOnInit(): void {
    if(this.appInitService.currentUser) {
      this.hasSignedIn = true;
    }
    this.appTokenService.tokenMessage$.subscribe(token=>this.token = token);
    this.route.queryParamMap.subscribe(params=>this.search(params.get("q")));
    console.log(this.appInitService.currentUser);
  }

  search(q: string) {
    this.globalSearchService.search(q).subscribe(search=>{
        this.globalSearch = search;
        this.route.queryParamMap.subscribe(params=>{
          params["token"] = this.token;
        })
      });
  }

  navigateTo(link) {
    this.appTokenService.token = this.token;
    this.router.navigate([link], {
      queryParams: {
        'token': this.appTokenService.token
      }
    });
  }

}