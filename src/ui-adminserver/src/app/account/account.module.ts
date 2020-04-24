import { NgModule } from '@angular/core';
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

@NgModule({
  declarations: [
    SignUpComponent,
    SignInComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    SharedModule,
    ClarityModule,
    BoardComponentsLibraryModule,
    TranslateModule,
  ],
  providers: [
    AccountService,
    HeaderComponent,
  ]
})
export class AccountModule { }
