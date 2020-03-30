import { Component, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AppInitService } from '../shared.service/app-init.service';
import {
  GUIDE_STEP,
  MAIN_MENU_DATA,
  RouteAudit,
  RouteHelm,
  RouteKibana,
  RouteNodes,
  RouteUserCenters
} from '../shared/shared.const';
import { ICsMenuItemData } from '../shared/shared.types';
import { SharedService } from '../shared.service/shared.service';

@Component({
  templateUrl: './main-content.component.html',
  styleUrls: ['./main-content.component.css']
})
export class MainContentComponent {
  @ViewChild('frameDashboard') frame: ElementRef;
  navSource: Array<ICsMenuItemData>;
  isSignIn = true;
  hasSignedIn = false;
  searchContent = '';

  constructor(private appInitService: AppInitService,
              private router: Router,
              private route: ActivatedRoute,
              private sharedService: SharedService) {
    if (this.appInitService.currentUser.user_id > 0) {
      this.isSignIn = false;
      this.hasSignedIn = true;
    }
    this.navSource = MAIN_MENU_DATA;
    this.getMenuItemByRoute(this.navSource, RouteNodes).visible = this.appInitService.isSystemAdmin;
    this.getMenuItemByRoute(this.navSource, RouteUserCenters).visible = this.appInitService.isSystemAdmin;
    this.getMenuItemByRoute(this.navSource, RouteAudit).visible = this.appInitService.isSystemAdmin;
    this.getMenuItemByRoute(this.navSource, RouteKibana).visible = !this.appInitService.isMipsSystem &&
      !this.appInitService.isArmSystem;
    this.getMenuItemByRoute(this.navSource, RouteHelm).visible = !this.appInitService.isMipsSystem &&
      !this.appInitService.isArmSystem;
    this.route.queryParamMap.subscribe(params => {
      this.searchContent = params.get('q');
    });
    this.appInitService.systemInfo = this.route.snapshot.data.systeminfo;
  }

  getMenuItemByRoute(source: Array<ICsMenuItemData>, route: string): ICsMenuItemData {
    let result: ICsMenuItemData;
    source.forEach(value => {
      if (value.url.includes(route)) {
        if (!result) {
          result = value;
        }
      } else if (!result && value.children) {
        result = this.getMenuItemByRoute(value.children, route);
      }
    });
    return result;
  }

  navigateTo(link) {
    this.router.navigate([link], {queryParams: {token: this.appInitService.token}}).then();
  }

  get isFirstLogin(): boolean {
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP {
    return this.appInitService.guideStep;
  }

  get showMaxGrafanaWindow(): boolean {
    return this.sharedService.showMaxGrafanaWindow;
  }

  get hideMaxGrafanaWindow(): boolean {
    return !this.sharedService.showMaxGrafanaWindow;
  }

  setGuideNoneStep() {
    this.appInitService.guideStep = GUIDE_STEP.NONE_STEP;
  }

  guideNextStep(step: GUIDE_STEP) {
    if (step === GUIDE_STEP.PROJECT_LIST) {
      this.navigateTo('/projects');
      this.appInitService.guideStep = GUIDE_STEP.CREATE_PROJECT;
    }
    if (step === GUIDE_STEP.SERVICE_LIST) {
      this.navigateTo('/services');
      this.appInitService.guideStep = GUIDE_STEP.CREATE_SERVICE;
    }
  }

}
