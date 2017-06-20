import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';

@Component({
    selector: 'board-app',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.scss']
})
export class AppComponent {
    constructor(private translateService: TranslateService) {
      let currentLang = translateService.getBrowserCultureLang();
      translateService.setDefaultLang(currentLang || 'zh-cn');
      this.translateService.onLangChange.subscribe(()=>{
        console.log('Changed lang is :' + this.translateService.currentLang);
      });
    }
}
