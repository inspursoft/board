import { NgModule } from '@angular/core';
import { RouterModule } from "@angular/router";
import { ImageListComponent } from './image-list/image-list.component';
import { ImageDetailComponent } from "./image-detail/image-detail.component";
import { SharedModule } from "../shared/shared.module";
import { ImageService } from "./image-service/image-service";
import { CreateImageComponent } from "./image-create/image-create.component";
import { CoreModule } from "../core/core.module";

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path:'', component: ImageListComponent}])
  ],
  providers: [
    ImageService
  ],
  entryComponents:[
    CreateImageComponent
  ],
  declarations: [
    ImageListComponent,
    CreateImageComponent,
    ImageDetailComponent
  ]
})
export class ImageModule {
}
