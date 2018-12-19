import { Component, ElementRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AppInitService, AppTokenService } from '../app.init.service';
import { GUIDE_STEP, MAIN_MENU_DATA, RouteAudit, RouteNodes, RouteUserCenters } from "../shared/shared.const";
import { ICsMenuItemData } from "../shared/shared.types";
import { SharedService } from "../shared/shared.service";

@Component({
  templateUrl: './main-content.component.html',
  styleUrls:['./main-content.component.css']
})  
export class MainContentComponent {
  @ViewChild("frameDashboard") frame:ElementRef;
  navSource: Array<ICsMenuItemData>;
  isSignIn: boolean = true;
  hasSignedIn: boolean = false;
  searchContent: string = '';

  constructor(
    private appInitService: AppInitService,
    private appTokenService: AppTokenService,
    private router: Router,
    private route: ActivatedRoute,
    private sharedService: SharedService) {
    if(this.appInitService.currentUser) {
      this.isSignIn = false;
      this.hasSignedIn = true;
    }
    this.navSource = MAIN_MENU_DATA;
    this.getMenuItemByRoute(RouteNodes).visible = this.appInitService.isSystemAdmin;
    this.getMenuItemByRoute(RouteUserCenters).visible = this.appInitService.isSystemAdmin;
    this.getMenuItemByRoute(RouteAudit).visible = this.appInitService.isSystemAdmin;
    this.route.queryParamMap.subscribe(params=>{
      this.searchContent = params.get("q");
    });
    this.appInitService.systemInfo = this.route.snapshot.data['systeminfo'];
  }

  getMenuItemByRoute(route: string): ICsMenuItemData {
    return this.navSource.find((value => value.url.includes(route)));
  }

  navigateTo(link) {
    this.router.navigate([link], {queryParams: {'token': this.appInitService.token}}).then()
  }

  get isFirstLogin(): boolean{
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP{
    return this.appInitService.guideStep;
  }

  get showMaxGrafanaWindow(): boolean{
    return this.sharedService.showMaxGrafanaWindow;
  }

  get hideMaxGrafanaWindow(): boolean{
    return !this.sharedService.showMaxGrafanaWindow;
  }

  setGuideNoneStep(){
    this.appInitService.guideStep = GUIDE_STEP.NONE_STEP;
  }

  guideNextStep(step: GUIDE_STEP){
    if (step == GUIDE_STEP.PROJECT_LIST){
      this.navigateTo('/projects');
      this.appInitService.guideStep = GUIDE_STEP.CREATE_PROJECT;
    }
    if (step == GUIDE_STEP.SERVICE_LIST){
      this.navigateTo('/services');
      this.appInitService.guideStep = GUIDE_STEP.CREATE_SERVICE;
    }
  }

}