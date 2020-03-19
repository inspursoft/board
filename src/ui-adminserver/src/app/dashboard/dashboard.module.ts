import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PreviewerComponent } from './previewer/previewer.component';
import { RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { DashboardService } from './dashboard.service';

@NgModule({
  declarations: [
    PreviewerComponent
  ],
  providers: [
    DashboardService
  ],
  imports: [
    CommonModule,
    ClarityModule,
    RouterModule.forChild([{ path: '', component: PreviewerComponent }])
  ],
})
export class DashboardModule { }
