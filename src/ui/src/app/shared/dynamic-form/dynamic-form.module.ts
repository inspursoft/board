import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { SharedModule } from '../shared.module';
import { DynamicFormComponent } from './dynamic-form.component';
import { DynamicFormControlComponent } from './dynamic-form-control.component';
import { DynamicFormService } from './dynamic-form.service';

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [
    DynamicFormComponent,
    DynamicFormControlComponent
  ],

  exports: [
    DynamicFormComponent,
    DynamicFormControlComponent
  ],
  providers: [ DynamicFormService ]
})
export class DynamicFormModule {}