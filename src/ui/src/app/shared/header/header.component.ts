import { Component, Input, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { TranslateService } from '@ngx-translate/core';

import { AppInitService } from '../../app.init.service';
import { AccountService } from '../../account/account.service';
import { MessageService } from '../message-service/message.service';
import { CookieService } from "ngx-cookie";

@Component({
  selector: 'header-content',
  templateUrl: 'header.component.html',
  styleUrls: [ 'header.component.css' ]
})
export class HeaderComponent implements OnInit {
  currentLang: string;
  @Input() isSignIn: boolean;
  @Input() hasSignedIn: boolean;
  @Input() searchContent: string;

  currentUser: {[key: string]: any};
  showChangePassword:boolean = false;
  showAccountSetting:boolean = false;
  authMode: string = '';
  redirectionURL: string = '';

  get brandLogoUrl(): string {
    return this.isSignIn ? '../../images/board-blue.jpg': '../../../images/board.png';
  }

  constructor(private router: Router,
              private translateService: TranslateService,
              private cookieService: CookieService,
              private appInitService: AppInitService,
              private accountService: AccountService,
              private messageService: MessageService) {
    this._assertLanguage(this.appInitService.currentLang);
  }

  ngOnInit(): void {
    if (this.hasSignedIn){
      this.currentUser = this.appInitService.currentUser || {};
      this.authMode = this.appInitService.systemInfo.auth_mode;
      this.redirectionURL = this.appInitService.systemInfo.redirection_url;
    }
  }

  _assertLanguage(lang: string) {
    lang = lang.toLowerCase();
    switch(lang) {
    case 'en':
    case 'en-us':
      lang = 'en-us';
      this.currentLang = 'HEAD_NAV.LANG_EN_US';
      break;

    case 'zh':
    case 'zh-cn': 
      lang = 'zh-cn';
      this.currentLang = 'HEAD_NAV.LANG_ZH_CN';
      break;
    }
    this.translateService.use(lang);
  }

  changLanguage(lang: string) {
    this._assertLanguage(lang);
  }

  doSearch(event) {
    this.searchContent = event.target.value;
    if(this.hasSignedIn) {
      this.router.navigate(['/search' ], { queryParams: { q: this.searchContent, token: this.appInitService.token }});
    } else {
      this.router.navigate(['/search' ], { queryParams: { q: this.searchContent }});
    }
  }

  clickLogoAction() {
    if(!this.hasSignedIn) {
      this.router.navigate(['/sign-in']);
    }
  }

  logOut() {
    this.accountService.signOut(this.appInitService.currentUser.user_name).subscribe(() => {
      this.cookieService.remove('token');
      this.appInitService.token = '';
      this.appInitService.currentUser = null;
      if (this.authMode === 'indata_auth') {
        window.location.href = this.redirectionURL;
        return;
      }
      this.router.navigate(['/sign-in']).then();
    }, () => this.messageService.showAlert('ACCOUNT.FAILED_TO_SIGN_OUT', {alertType: 'alert-danger'}));
  }
}