import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';

import 'rxjs/add/operator/switchMap';

import { AppInitService } from '../app.init.service';
import { GlobalSearchService } from './global-search.service';

import { MessageService } from '../shared/message-service/message.service';

@Component({
  selector: 'global-search',
  templateUrl: 'global-search.component.html'
})
export class GlobalSearchComponent implements OnInit {

  hasSignedIn: boolean;
  globalSearch: {[key: string]: any} = {};
  
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private appInitService: AppInitService,
    private globalSearchService: GlobalSearchService,
    private messageService: MessageService
  ) {}

  ngOnInit(): void {
    this.hasSignedIn = (this.appInitService.currentUser !== null);
    this.route.queryParamMap.subscribe(params=>{
      this.search(params.get("q"));
    });
  }

  search(q: string) {
    this.globalSearchService
      .search(q)
      .then(search=>{
        this.globalSearch = search;
      })
      .catch(err=>this.messageService.dispatchError(err));
  }

  navigateTo(link) {
    this.router.navigate([link], {
      queryParams: {
        'token': this.appInitService.token
      }
    });
  }

}