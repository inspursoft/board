import { NgModule } from '@angular/core';
import { SharedModule } from '../shared.module';
import { DynamicFormComponent } from './dynamic-form.component';
import { DynamicFormControlComponent } from './dynamic-form-control.component';
import { DynamicFormService } from './dynamic-form.service';
import { CoreModule } from "../../core/core.module";

@NgModule({
  imports: [
    CoreModule,
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
