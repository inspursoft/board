import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { AppInitService } from '../shared.service/app-init.service';
import { AppTokenService } from '../shared.service/app-token.service';
import { SharedService } from "../shared.service/shared.service";

@Component({
  selector: 'global-search',
  templateUrl: 'global-search.component.html'
})
export class GlobalSearchComponent implements OnInit {

  token: string;

  hasSignedIn: boolean;
  globalSearch: { [key: string]: any } = {};

  constructor(private router: Router,
              private route: ActivatedRoute,
              private appInitService: AppInitService,
              private appTokenService: AppTokenService,
              private sharedService: SharedService,
  ) {}

  ngOnInit(): void {
    if (this.appInitService.currentUser.user_id > 0) {
      this.hasSignedIn = true;
    }
    this.route.queryParamMap.subscribe(params => this.search(params.get('q')));
    this.route.queryParamMap.subscribe(params => params["token"] = this.token);
  }

  search(q: string) {
    this.sharedService.search(q).subscribe(search => this.globalSearch = search);
  }

  navigateTo(link) {
    this.router.navigate([link], {
      queryParams: {
        token: this.appTokenService.token
      }
    });
  }

}
