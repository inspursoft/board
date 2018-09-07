import { NgModule } from '@angular/core';

import { ImageListComponent } from './image-list/image-list.component';
import { ImageDetailComponent } from "./image-detail/image-detail.component";
import { SharedModule } from "../shared/shared.module";
import { ImageService } from "./image-service/image-service";
import { CreateImageComponent } from "./image-create/image-create.component";

@NgModule({
  imports: [SharedModule],
  providers: [ImageService],
  entryComponents:[CreateImageComponent],
  declarations: [
    ImageListComponent,
    CreateImageComponent,
    ImageDetailComponent]
})
export class ImageModule {
}