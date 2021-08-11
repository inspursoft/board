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
import { NodeServiceControlComponent } from './node-service-control/node-service-control.component';
import { NodeGroupControlComponent } from './node-group-control/node-group-control.component';
import { NodeCreateNewComponent } from './node-create-new/node-create-new.component';
import { NodeCreateNewGuideComponent } from './node-create-new-guide/node-create-new-guide.component';
import { NodeGroupEditComponent } from './node-group-edit/node-group-edit.component';

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
    NodeControlComponent,
    NodeServiceControlComponent,
    NodeGroupControlComponent,
    NodeCreateNewComponent,
    NodeCreateNewGuideComponent,
    NodeGroupEditComponent
  ],
  entryComponents: [
    NodeControlComponent,
    NodeGroupControlComponent,
    NodeCreateGroupComponent,
    NodeServiceControlComponent,
    NodeCreateNewComponent,
    NodeGroupEditComponent
  ],
  providers: [
    NodeService
  ]
})

export class NodeModule {
}
