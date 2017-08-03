import { OnInit, Component } from '@angular/core';
import { Image } from "../image";
import { ImageService } from "../image-service/image-service"
import { MessageService } from "../../shared/message-service/message.service";

@Component({
  selector: 'image-list',
  templateUrl: './image-list.component.html',
  styleUrls: ["./image-list.component.css"]
})
export class ImageListComponent implements OnInit {
  curPage: number = 1;
  curImage: Image;
  isShowDetail: boolean = false;
  imageListErrMsg: string = "";
  imageList: Image[] = Array<Image>();
  imageCountPerPage: number = 2;

  constructor(private imageService: ImageService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.imageService.getImages("", 0, 0)
      .then(res => this.imageList = res)
      .catch(err => this.messageService.dispatchError(err));
  }

  showImageDetail(image: Image) {
    //need add get one Image from server
    this.curImage = image;
    this.isShowDetail = true;
  }

  pageChange(pageIndex) {
    this.curPage = pageIndex;
  }
}