import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';

import { SharedModule } from '../shared/shared.module';
import { NodeComponent } from './node.component';
import { NodeDetailComponent } from './node-detail.component';
import { NodeService } from './node.service';

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [ 
    NodeComponent,
    NodeDetailComponent
  ],

  providers: [
    NodeService
  ]
})
export class NodeModule {}