import { Component, ElementRef, ViewChild } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AppInitService, AppTokenService } from '../app.init.service';
import { GUIDE_STEP } from "../shared/shared.const";

@Component({
  selector: 'main-content',
  templateUrl: 'main-content.component.html'
})  
export class MainContentComponent {
  @ViewChild("frameDashboard") frame:ElementRef;
  token: string;
  isOnlyShowGrafanaView: boolean = false;
  isSignIn: boolean = true;
  hasSignedIn: boolean = false;
  searchContent: string = '';

  constructor(
    private appInitService: AppInitService,
    private appTokenService: AppTokenService,
    private router: Router,
    private route: ActivatedRoute
  ) {
    if(this.appInitService.currentUser) {
      this.isSignIn = false;
      this.hasSignedIn = true;
    }
    this.token = this.appTokenService.token;
    this.appTokenService.tokenMessage$.subscribe(token=>this.token = token);
    this.route.queryParamMap.subscribe(params=>{
      this.isOnlyShowGrafanaView = params.get("isOnlyShowGrafanaView") == "true";
      this.searchContent = params.get("q");
    });
    let systemInfo = this.route.snapshot.data['systeminfo'];
    this.appInitService.systemInfo = systemInfo;
    this.appInitService.grafanaViewUrl = `http://${systemInfo['board_host']}:3000/dashboard/db/kubernetes-cluster?refresh=30s&orgId=1`;
  }

  get grafanaViewUrl():string{
    return this.appInitService.grafanaViewUrl;
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

  get isFirstLogin(): boolean{
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP{
    return this.appInitService.guideStep;
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