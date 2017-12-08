import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { Image, ImageDetail } from "../image"
import { ImageService } from "../image-service/image-service";
import { MessageService } from "../../shared/message-service/message.service";
import { MESSAGE_TARGET, BUTTON_STYLE, MESSAGE_TYPE } from '../../shared/shared.const';
import { Message } from '../../shared/message-service/message';

@Component({
  selector: "image-detail",
  templateUrl: "./image-detail.component.html",
  styleUrls: ["./image-detail.component.css"]
})

export class ImageDetailComponent implements OnInit {
  _isOpen: boolean;
  alertClosed: boolean;
  @Input() curImage: Image;
  showDeleteAlert: boolean[];
  imageDetailPageSize: number = 10;
  imageDetailErrMsg: string = "";
  imageDetailList: ImageDetail[] = Array<ImageDetail>();

  loadingWIP: boolean;
  @Output() reload = new EventEmitter<boolean>();

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
      this.loadingWIP = true;
      this.imageService.getImageDetailList(this.curImage.image_name)
        .then((res: ImageDetail[]) => {
          this.loadingWIP = false;
          for (let item of res || []) {
            item['image_detail'] = JSON.parse(item['image_detail']);
            item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
            item['image_size_unit'] = 'MB';
          }
          this.showDeleteAlert = new Array(this.imageDetailList.length);
          this.imageDetailList = res || [];
        })
        .catch(err => {
          this.loadingWIP = false;
          this.messageService.dispatchError(err)
        });
    }
  }

  deleteTag(tagName: string) {
    let m: Message = new Message();
    this.imageService
      .deleteImageTag(this.curImage.image_name, tagName)
      .then(res => {
        m.message = 'IMAGE.SUCCESSFUL_DELETED_TAG';
        this.messageService.inlineAlertMessage(m);
        this.reload.emit(true);
        this.isOpen = false;
      })
      .catch(err => {
        m.message = 'IMAGE.FAILED_TO_DELETE_TAG';
        m.type = MESSAGE_TYPE.COMMON_ERROR;
        this.messageService.inlineAlertMessage(m);
      });
  }
}