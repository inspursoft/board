import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SignUpComponent } from './sign-up/sign-up.component';
import { SignInComponent } from './sign-in/sign-in.component';
import { FormsModule } from '@angular/forms';
import { SharedModule } from '../shared/shared.module';
import { AccountService } from './account.service';
import { BoardComponentsLibraryModule } from 'board-components-library';
import { TranslateModule } from '@ngx-translate/core';

@NgModule({
  declarations: [
    SignUpComponent,
    SignInComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    SharedModule,
    BoardComponentsLibraryModule,
    TranslateModule,
  ],
  providers: [
    AccountService,
  ]
})
export class AccountModule { }
