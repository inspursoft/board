import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceComponent } from './resource.component';
import { ClarityModule } from '@clr/angular';
import { RouterModule } from '@angular/router';
import { ComputeComponent } from './compute/compute.component';
import { ResourceService } from './services/resource.service';
import { WebsocketService } from './services/websocket.service';
import { CustomHttpProvider } from './services/custom-http.service';
import { NodeAddRemoveComponent } from './node-add-remove/node-add-remove.component';
import { TranslateModule } from '@ngx-translate/core';
import { BoardComponentsLibraryModule } from 'board-components-library';

@NgModule({
  declarations: [ResourceComponent, ComputeComponent, NodeAddRemoveComponent],
  imports: [
    CommonModule,
    ClarityModule,
    RouterModule.forChild([{path: '', component: ResourceComponent}]),
    TranslateModule,
    BoardComponentsLibraryModule
  ],
  entryComponents: [
    NodeAddRemoveComponent
  ],
  providers: [
    ResourceService,
    WebsocketService,
    CustomHttpProvider
  ]
})
export class ResourceModule {
}
