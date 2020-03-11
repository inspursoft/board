import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core"
import { Image, ImageDetail } from "../image"
import { ImageService } from "../image-service/image-service";
import { MessageService } from "../../shared.service/message.service";

@Component({
  selector: "image-detail",
  templateUrl: "./image-detail.component.html",
  styleUrls: ["./image-detail.component.css"]
})

export class ImageDetailComponent implements OnInit {
  _isOpen: boolean;
  @Input() curImage: Image;
  showDeleteAlert: boolean[];
  imageDetailPageSize: number = 10;
  imageDetailList: ImageDetail[] = Array<ImageDetail>();

  loadingWIP: boolean;
  @Output() reload = new EventEmitter<boolean>();

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(value: boolean) {
    this._isOpen = value;
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
      this.loadingWIP = true;
      this.imageService.getImageDetailList(this.curImage.image_name).subscribe((res: ImageDetail[]) => {
          this.loadingWIP = false;
          for (let item of res) {
            if (item['image_detail'] && item['image_detail'] != ''){
              item['image_detail'] = JSON.parse(item['image_detail']);
            }
            item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
            item['image_size_unit'] = 'MB';
          }
          this.showDeleteAlert = new Array(this.imageDetailList.length);
          this.imageDetailList = res || [];
        }, () => this.loadingWIP = false
      );
    }
  }

  deleteTag(tagName: string) {
    this.imageService.deleteImageTag(this.curImage.image_name, tagName).subscribe(() => {
        this.reload.emit(true);
        this.isOpen = false;
        this.messageService.showAlert('IMAGE.SUCCESSFUL_DELETED_TAG');
      },() => this.isOpen = false
    )
  }
}
