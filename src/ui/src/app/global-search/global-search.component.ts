import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AppInitService } from '../shared.service/app-init.service';
import { GlobalSearchService } from './global-search.service';
import { GlobalSearchResult } from './global-search.types';

@Component({
  templateUrl: 'global-search.component.html'
})
export class GlobalSearchComponent implements OnInit {
  token = '';
  hasSignedIn = false;
  globalSearch: GlobalSearchResult;

  constructor(private router: Router,
              private route: ActivatedRoute,
              private appInitService: AppInitService,
              private globalSearchService: GlobalSearchService) {
    this.globalSearch = new GlobalSearchResult();
  }

  ngOnInit(): void {
    this.hasSignedIn = this.appInitService.currentUser.userId > 0;
    this.route.queryParamMap.subscribe(params => this.search(params.get('q')));
  }

  search(q: string) {
    this.globalSearchService.search(q, this.appInitService.token).subscribe(
      res => this.globalSearch = res
    );
  }

  navigateTo(link) {
    this.router.navigate([link], {
      queryParams: {
        token: this.appInitService.token
      }
    });
  }

}
