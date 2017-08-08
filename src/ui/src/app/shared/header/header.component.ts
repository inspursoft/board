import { Component, Input, OnInit } from '@angular/core';
import { Router } from '@angular/router';

import { TranslateService, LangChangeEvent } from '@ngx-translate/core';

import { AppInitService } from '../../app.init.service';
import { AccountService } from '../../account/account.service';
import { MessageService } from '../message-service/message.service';

@Component({
  selector: 'header-content',
  templateUrl: 'header.component.html',
  styleUrls: [ 'header.component.css' ]
})
export class HeaderComponent implements OnInit {

  currentLang: string;
  @Input() isSignIn: boolean;
  @Input() hasSignedIn: boolean;
  currentUser: {[key: string]: any};
  showChangePassword:boolean = false;

  get brandLogoUrl(): string {
    return this.isSignIn ? '../../images/board-blue.jpg': '../../../images/board.png';
  }

  constructor(
    private router: Router,
    private translateService: TranslateService,
    private appInitService: AppInitService,
    private accountService: AccountService,
    private messageService: MessageService) {
    let lang: string = this.translateService.getBrowserCultureLang();
    this._assertLanguage(lang);
  }

  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser || {};
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

  logOut() {
    this.accountService
      .signOut()
      .then(res=>this.router.navigate(['/sign-in']))
      .catch(err=>this.messageService.dispatchError(err, 'ACCOUNT.FAILED_TO_SIGN_OUT'));
  }
}