import { Component } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { AppInitService } from '../app.init.service';
import { GUIDE_STEP } from "../shared/shared.const";

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
    this.appInitService.systemInfo = this.route.snapshot.data['systeminfo'];
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

  get someAny(): boolean{
    if (this.appInitService.guideStep == GUIDE_STEP.SERVICE_LIST){
      console.log(3);
    }
    return true
  }

  guideNextStep(step: GUIDE_STEP){
    if (step == GUIDE_STEP.PROJECT_LIST){
      this.navigateTo('/projects');
      this.appInitService.guideStep = GUIDE_STEP.CREATE_PROJECT;
    }
  }
}