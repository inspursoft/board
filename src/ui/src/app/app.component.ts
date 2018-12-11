import { AfterViewInit, Component, ComponentFactoryResolver, ViewChild, ViewContainerRef } from '@angular/core';
import { AppInitService } from './app.init.service';
import { TranslateService } from '@ngx-translate/core';
import { CookieService } from "ngx-cookie";
import { MessageService } from "./shared/message-service/message.service";

@Component({
  selector: 'board-app',
  templateUrl: './app.component.html',
  styleUrls:['./app.component.css']
})
export class AppComponent implements AfterViewInit {
  @ViewChild('messageContainer', {read: ViewContainerRef}) messageContainer;
  cookieExpiry: Date = new Date(Date.now() + 60 * 60 * 24 * 365 * 1000);

  constructor(private appInitService: AppInitService,
              private cookieService: CookieService,
              private messageService: MessageService,
              private resolver: ComponentFactoryResolver,
              private translateService: TranslateService){
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
  }

  ngAfterViewInit() {
    this.messageService.registerDialogHandle(this.messageContainer, this.resolver);
  }
}
