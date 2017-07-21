import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { Image, ImageDetail } from "../image"
import { ImageService } from "../image-service/image-service";

@Component({
  selector: "image-detail",
  templateUrl: "./image-detail.component.html",
  styleUrls: ["./image-detail.component.css"]
})

export class ImageDetailComponent implements OnInit {
  _isOpen: boolean;
  @Input() curImage: Image;
  curPage: number = 1;
  imageDetailPageSize: number = 10;
  imageDetailErrMsg: string;
  imageDetailList: ImageDetail[] = Array<ImageDetail>();

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(value: boolean) {
    this._isOpen = value;
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private imageService: ImageService) {
  }

  ngOnInit() {
    this.getImageDetailList();
  }

  getImageDetailList() {
    if (this.curImage && this.curImage.image_name) {
      this.imageService.getImageDetailList(this.curImage.image_name)
        .then(res => this.imageDetailList = res)
        .catch((reason: string) => this.imageDetailErrMsg = reason);
    }
  }

  pageChange(pageIndex: number) {

  }
}