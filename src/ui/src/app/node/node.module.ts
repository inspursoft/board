import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { NodeComponent } from './node.component';
import { NodeService } from './node.service';
import { NodeGroupComponent } from './node-group/node-group.component';
import { NodeListComponent } from './node-list/node-list.component';
import { NodeDetailComponent } from './node-detail/node-detail.component';
import { NodeCreateGroupComponent } from './node-create-group/node-create-group.component';
import { NodeControlComponent } from './node-control/node-control.component';
import { CoreModule } from '../core/core.module';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: NodeComponent}])
  ],
  declarations: [
    NodeComponent,
    NodeDetailComponent,
    NodeGroupComponent,
    NodeListComponent,
    NodeCreateGroupComponent,
    NodeControlComponent
  ],
  entryComponents: [
    NodeCreateGroupComponent
  ],
  providers: [
    NodeService
  ]
})

export class NodeModule {
}
