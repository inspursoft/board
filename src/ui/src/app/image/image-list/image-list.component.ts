import { OnInit, Component } from '@angular/core';
import { Image } from "../image";
import { ImageService } from "../image-service/image-service"

@Component({
  selector: 'image-list',
  templateUrl: './image-list.component.html',
  styleUrls:["./image-list.component.css"]
})
export class ImageListComponent implements OnInit {
  curPage: number = 1;
  curImage: Image;
  isShowDetail: boolean = false;
  imageListErrMsg: string = "";
  imageList: Image[] = Array<Image>();
  imageCountPerPage: number = 10;

  constructor(private imageService: ImageService) {
  }

  ngOnInit() {
    this.imageService.getImages("", 0, 0)
      .then(res => this.imageList = res)
      .catch((reason: string) => this.imageListErrMsg = reason);
  }

  showImageDetail(image: Image) {
    //need add get one Image from server
    this.curImage = image;
    this.isShowDetail = true;
  }


  pageChange(pageIndex: number) {

  }
}