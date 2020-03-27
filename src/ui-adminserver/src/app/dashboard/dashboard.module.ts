import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PreviewerComponent } from './previewer/previewer.component';
import { RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { DashboardService } from './dashboard.service';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

@NgModule({
  declarations: [
    PreviewerComponent
  ],
  providers: [
    DashboardService
  ],
  imports: [
    CommonModule,
    FormsModule,
    ClarityModule,
    TranslateModule,
    RouterModule.forChild([{ path: '', component: PreviewerComponent }])
  ],
})
export class DashboardModule { }
