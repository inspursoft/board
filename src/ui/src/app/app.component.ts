import { AfterViewInit, Component, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { AppInitService } from './shared.service/app-init.service';
import { LangChangeEvent, TranslateService } from '@ngx-translate/core';
import { CookieService } from 'ngx-cookie';
import { MessageService } from './shared.service/message.service';
import { registerLocaleData } from '@angular/common';
import localeZhHans from '@angular/common/locales/zh-Hans';

@Component({
  selector: 'app-board',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements AfterViewInit {
  @ViewChild('messageContainer', {read: ViewContainerRef}) messageContainer;

  constructor(private appInitService: AppInitService,
              private cookieService: CookieService,
              private messageService: MessageService,
              private resolver: ComponentFactoryResolver,
              private translateService: TranslateService) {
    let currentLang = localStorage.getItem('currentLang');
    if (!currentLang) {
      console.log('No found cookie for current lang, will use the default browser language.');
      currentLang = this.translateService.getBrowserCultureLang();
      localStorage.setItem('currentLang', currentLang || 'en-us');
    }
    this.appInitService.currentLang = currentLang || 'en-us';
    translateService.use(currentLang || 'en-us');
    this.translateService.onLangChange.subscribe((res: LangChangeEvent) => {
      const oldLang = this.appInitService.currentLang;
      currentLang = this.translateService.currentLang;
      this.appInitService.currentLang = currentLang;
      localStorage.setItem('currentLang', currentLang);
      if (res.lang.toLocaleLowerCase() !== oldLang.toLocaleLowerCase()) {
        window.location.reload(true);
      }
    });
    if (appInitService.currentLang.toLowerCase() === 'zh-cn') {
      registerLocaleData(localeZhHans, 'zh-Hans');
    }
  }

  ngAfterViewInit() {
    this.messageService.registerDialogHandle(this.messageContainer, this.resolver);
  }
}
