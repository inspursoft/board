import { AfterViewInit, Component, ComponentFactoryResolver, HostListener, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { AppInitService } from './app.init.service';
import { TranslateService } from '@ngx-translate/core';
import { CookieService } from "ngx-cookie";
import { MessageService } from "./shared/message-service/message.service";

@Component({
  selector: 'board-app',
  templateUrl: './app.component.html',
  styleUrls:['./app.component.css']
})
export class AppComponent implements AfterViewInit, OnInit {
  @ViewChild('messageContainer', {read: ViewContainerRef}) messageContainer;
  cookieExpiry: Date = new Date(Date.now() + 60 * 60 * 24 * 365 * 1000);
  monthNameMap: Map<string, string>;
  monthBriefNameMap: Map<string, string>;

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
    this.translateService.onLangChange.subscribe(() => {
      this.appInitService.currentLang = this.translateService.currentLang;
      cookieService.put('currentLang', this.appInitService.currentLang, {expires: this.cookieExpiry});
      console.log('Change lang to:' + this.appInitService.currentLang);
    });
    this.monthNameMap = new Map<string, string>();
    this.monthBriefNameMap = new Map<string, string>();
  }

  ngOnInit() {
    this.monthNameMap.set('January', '一月');
    this.monthNameMap.set('February', '二月');
    this.monthNameMap.set('March', '三月');
    this.monthNameMap.set('April', '四月');
    this.monthNameMap.set('May', '五月');
    this.monthNameMap.set('June', '六月');
    this.monthNameMap.set('July', '七月');
    this.monthNameMap.set('August', '八月');
    this.monthNameMap.set('September', '九月');
    this.monthNameMap.set('October', '十月');
    this.monthNameMap.set('November', '十一月');
    this.monthNameMap.set('December', '十二月');

    this.monthBriefNameMap.set('Jan', '一月');
    this.monthBriefNameMap.set('Feb', '二月');
    this.monthBriefNameMap.set('Mar', '三月');
    this.monthBriefNameMap.set('Apr', '四月');
    this.monthBriefNameMap.set('May', '五月');
    this.monthBriefNameMap.set('Jun', '六月');
    this.monthBriefNameMap.set('Jul', '七月');
    this.monthBriefNameMap.set('Aug', '八月');
    this.monthBriefNameMap.set('Sep', '九月');
    this.monthBriefNameMap.set('Oct', '十月');
    this.monthBriefNameMap.set('Nov', '十一月');
    this.monthBriefNameMap.set('Dec', '十二月');
  }
  ngAfterViewInit() {
    this.messageService.registerDialogHandle(this.messageContainer, this.resolver);
  }

  @HostListener('click', ['$event.target']) clickEvent(element: HTMLElement) {
    let btnTrigger = document.getElementsByClassName('calendar-btn monthpicker-trigger');
    if (btnTrigger.length > 0 && this.appInitService.currentLang == 'zh-cn') {
      let oldText = (btnTrigger[0] as HTMLButtonElement).innerText;
      (btnTrigger[0] as HTMLButtonElement).innerText = this.monthBriefNameMap.get(oldText)
    }
    let btnMonths = document.getElementsByClassName('calendar-btn month ng-star-inserted');
    if (btnMonths.length > 0 && this.appInitService.currentLang == 'zh-cn') {
      for (let i = 0; i < btnMonths.length; i++) {
        let oldText = (btnMonths[i] as HTMLButtonElement).innerText;
        (btnMonths[i] as HTMLButtonElement).innerText = this.monthNameMap.get(oldText)
      }
    }
  }
}
