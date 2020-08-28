import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ImageService } from '../image.service';
import { MessageService } from '../../shared.service/message.service';
import { Image, ImageDetail } from '../image.types';

@Component({
  selector: 'app-image-detail',
  templateUrl: './image-detail.component.html',
  styleUrls: ['./image-detail.component.css']
})

export class ImageDetailComponent implements OnInit {
  @Input() curImage: Image;
  isOpenValue: boolean;
  showDeleteAlert: Array<boolean>;
  imageDetailPageSize = 10;
  imageDetailList: Array<ImageDetail>;

  loadingWIP: boolean;
  @Output() reload = new EventEmitter<boolean>();

  @Input()
  get isOpen() {
    return this.isOpenValue;
  }

  set isOpen(value: boolean) {
    this.isOpenValue = value;
    this.isOpenChange.emit(this.isOpenValue);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private imageService: ImageService,
              private messageService: MessageService) {
    this.showDeleteAlert = new Array<boolean>();
    this.imageDetailList = new Array<ImageDetail>();
  }

  ngOnInit() {
    this.getImageDetailList();
  }

  getImageDetailList() {
    this.loadingWIP = true;
    this.imageService.getImageDetailList(this.curImage.imageName).subscribe(
      (res: Array<ImageDetail>) => {
        this.imageDetailList = res;
        this.loadingWIP = false;
        for (const detail of res) {
          detail.imageSizeNumber = Number.parseFloat((detail.imageSizeNumber / (1024 * 1024)).toFixed(2));
          detail.imageSizeUnit = 'MB';
        }
        this.showDeleteAlert = new Array(this.imageDetailList.length);
      }, () => this.loadingWIP = false
    );
  }

  deleteTag(tagName: string) {
    this.imageService.deleteImageTag(this.curImage.imageName, tagName).subscribe(() => {
        this.reload.emit(true);
        this.isOpen = false;
        this.messageService.showAlert('IMAGE.SUCCESSFUL_DELETED_TAG');
      }, () => this.isOpen = false
    );
  }
}
