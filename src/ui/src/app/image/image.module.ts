import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';

import { ImageListComponent } from './image-list/image-list.component';
import { ImageDetailComponent } from "./image-detail/image-detail.component";
import { SharedModule } from "../shared/shared.module";
import { ImageService } from "./image-service/image-service";

@NgModule({
  imports: [SharedModule],
  providers: [ImageService],

  declarations: [
    ImageListComponent,
    ImageDetailComponent]
})
export class ImageModule {
}