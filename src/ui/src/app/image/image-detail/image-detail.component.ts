import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { Image, ImageDetail } from "../image"
import { ImageService } from "../image-service/image-service";
import { MessageService } from "../../shared/message-service/message.service";

@Component({
  selector: "image-detail",
  templateUrl: "./image-detail.component.html",
  styleUrls: ["./image-detail.component.css"]
})

export class ImageDetailComponent implements OnInit {
  _isOpen: boolean;
  alertClosed: boolean;
  @Input() curImage: Image;
  showDeleteAlert: boolean = false;
  imageDetailPageSize: number = 1;
  imageDetailErrMsg: string = "";
  imageDetailList: ImageDetail[] = Array<ImageDetail>();

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(value: boolean) {
    this._isOpen = value;
    this.alertClosed = true;
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private imageService: ImageService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.getImageDetailList();
  }

  getImageDetailList() {
    if (this.curImage && this.curImage.image_name) {
      this.imageService.getImageDetailList(this.curImage.image_name)
        .then((res: ImageDetail[]) => {
          for (let item of res) {
            item['image_detail'] = JSON.parse(item['image_detail']);
            item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
            item['image_size_unit'] = 'MB';
          }
          this.imageDetailList = res;
        })
        .catch(err => this.messageService.dispatchError(err));
    }
  }
}