import { NgModule, LOCALE_ID, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SignUpComponent } from './sign-up/sign-up.component';
import { SignInComponent } from './sign-in/sign-in.component';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../shared/shared.module';
import { AccountService } from './account.service';
import { BoardComponentsLibraryModule } from 'board-components-library';
import { TranslateModule } from '@ngx-translate/core';
import { ClarityModule } from '@clr/angular';
import { HeaderComponent } from '../shared/header/header.component';
import { InstallationComponent } from './installation/installation.component';
import { SharedServiceModule } from '../shared.service/shared-service.module';
import { ConfigurationModule } from '../configuration/configuration.module';

@NgModule({
  declarations: [
    SignUpComponent,
    SignInComponent,
    InstallationComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    SharedModule,
    ClarityModule,
    BoardComponentsLibraryModule,
    TranslateModule,
    SharedServiceModule,
    ConfigurationModule,
  ],
  providers: [
    AccountService,
    HeaderComponent,
    [{ provide: LOCALE_ID, useFactory: GetLocale }]
  ]
})
export class AccountModule { }

function GetLocale(): string {
  const currentLang = (window.localStorage.getItem('currentLang') === 'zh-cn' || window.localStorage.getItem('currentLang') === 'zh');
  if (currentLang) {
    return 'zh-Hans';
  } else {
    return 'en-US';
  }
}
