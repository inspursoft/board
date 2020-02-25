import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceComponent } from './resource.component';
import { ClarityModule } from '@clr/angular';
import { RouterModule } from '@angular/router';
import { ResourceService } from './services/resource.service';
import { CustomHttpProvider } from './services/custom-http.service';
import { TranslateModule } from '@ngx-translate/core';
import { BoardComponentsLibraryModule } from 'board-components-library';
import { NodeListComponent } from './compute/node-list/node-list.component';
import { NodeLogsComponent } from './compute/node-logs/node-logs.component';
import { NodeDetailComponent } from './compute/node-detail/node-detail.component';

@NgModule({
  declarations: [
    ResourceComponent,
    NodeListComponent,
    NodeLogsComponent,
    NodeDetailComponent
  ],
  imports: [
    CommonModule,
    ClarityModule,
    RouterModule.forChild([
      {path: '', redirectTo: 'node-list', pathMatch: 'full'},
      {
        path: '', component: ResourceComponent, children: [
          {path: 'node-list', component: NodeListComponent},
          {path: 'node-logs', component: NodeLogsComponent}
        ]
      }]),
    TranslateModule,
    BoardComponentsLibraryModule
  ],
  entryComponents: [NodeDetailComponent],
  providers: [
    ResourceService,
    CustomHttpProvider
  ]
})
export class ResourceModule {
}
