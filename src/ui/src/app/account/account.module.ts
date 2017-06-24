import { NgModule } from '@angular/core';

import { SharedModule } from '../shared/shared.module';
import { SignInComponent } from './sign-in/sign-in.component';

@NgModule({
  imports: [ 
    SharedModule 
  ],
  declarations: [ 
    SignInComponent 
  ]
})
export class AccountModule {}