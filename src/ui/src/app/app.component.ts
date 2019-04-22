import { AfterViewInit, Component, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { AppInitService } from './shared.service/app-init.service';
import { LangChangeEvent, TranslateService } from '@ngx-translate/core';
import { CookieService } from 'ngx-cookie';
import { MessageService } from './shared.service/message.service';
import { registerLocaleData } from '@angular/common';
import localeZhHans from '@angular/common/locales/zh-Hans';

@Component({
  selector: 'board-app',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterViewInit {
  @ViewChild('messageContainer', {read: ViewContainerRef}) messageContainer;
  cookieExpiry: Date = new Date(Date.now() + 60 * 60 * 24 * 365 * 1000);

  constructor(private appInitService: AppInitService,
              private cookieService: CookieService,
              private messageService: MessageService,
              private resolver: ComponentFactoryResolver,
              private translateService: TranslateService) {
    if (!cookieService.get('currentLang')) {
      console.log('No found cookie for current lang, will use the default browser language.');
      cookieService.put('currentLang', this.translateService.getBrowserCultureLang(), {expires: this.cookieExpiry});
    }
    this.appInitService.currentLang = cookieService.get('currentLang') || 'en-us';
    translateService.use(this.appInitService.currentLang);
    this.translateService.onLangChange.subscribe((res: LangChangeEvent) => {
      const oldLang = this.appInitService.currentLang;
      this.appInitService.currentLang = this.translateService.currentLang;
      cookieService.put('currentLang', this.appInitService.currentLang, {expires: this.cookieExpiry});
      if (res.lang.toLocaleLowerCase() !== oldLang.toLocaleLowerCase()) {
        window.location.reload(true);
      }
    });
    if (appInitService.currentLang === 'zh-cn') {
      registerLocaleData(localeZhHans, 'zh-Hans');
    }
  }

  ngAfterViewInit() {
    this.messageService.registerDialogHandle(this.messageContainer, this.resolver);
  }
}
