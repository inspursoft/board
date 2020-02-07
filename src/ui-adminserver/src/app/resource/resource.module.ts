import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceComponent } from './resource.component';
import { ClarityModule } from '@clr/angular';
import { RouterModule } from '@angular/router';

@NgModule({
  declarations: [ResourceComponent],
  imports: [
    CommonModule,
    ClarityModule,
    RouterModule.forChild([{ path: '', component: ResourceComponent }])
  ]
})
export class ResourceModule { }
