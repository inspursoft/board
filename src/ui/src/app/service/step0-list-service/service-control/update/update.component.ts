import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ImageIndex } from '../../../service-step.component';
import { ImageDetail } from '../../../../image/image';
import { K8sService } from '../../../service.k8s';
import { Service } from '../../../service';

@Component({
  selector: 'app-update',
  templateUrl: './update.component.html',
  styleUrls: ['./update.component.css']
})
export class UpdateComponent implements OnInit {
  @Input() service: Service;
  @Input() isActionInWIP: boolean;
  @Output() messageEvent: EventEmitter<string>;
  @Output() errorEvent: EventEmitter<any>;
  @Output() alertEvent: EventEmitter<string>;
  @Output() actionIsEnabledEvent: EventEmitter<boolean>;
  imageList: Array<ImageIndex>;
  imageTagList: Map<string, Array<ImageDetail>>;
  imageTagSelected: Map<string, string>;

  constructor(private k8sService: K8sService) {
    this.imageList = Array<ImageIndex>();
    this.errorEvent = new EventEmitter<any>();
    this.messageEvent = new EventEmitter<string>();
    this.alertEvent = new EventEmitter<string>();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
    this.imageTagList = new Map<string, Array<ImageDetail>>();
    this.imageTagSelected = new Map<string, string>();
  }

  ngOnInit() {
    this.actionIsEnabledEvent.emit(false);
    this.getServiceImages();
  }

  getServiceImages() {
    this.imageTagSelected.clear();
    this.imageList.splice(0, this.imageList.length);
    this.k8sService.getServiceImages(this.service.service_project_name, this.service.service_name).subscribe(
      (imageList: Array<ImageIndex>) => {
        this.imageList = imageList;
        this.imageList.forEach(value => {
          this.imageTagSelected.set(value.image_name, value.image_tag);
          this.k8sService.getImageDetailList(value.image_name).subscribe(
            (imageDetailList: Array<ImageDetail>) => {
              if (imageDetailList.length === 0) {
                const tag = new ImageDetail();
                const tagList = Array<ImageDetail>();
                tag.image_tag = value.image_tag;
                tagList.push(tag);
                this.imageTagList.set(value.image_name, tagList);
              } else {
                this.imageTagList.set(value.image_name, imageDetailList);
              }
              this.setActionEnabled();
            }, err => this.errorEvent.next(err)
          );
        });
      },
      (err: HttpErrorResponse) => {
        if (err.status === 500) {
          this.alertEvent.emit('SERVICE.SERVICE_CONTROL_NOT_UPDATE');
        } else {
          this.errorEvent.next(err);
        }
      });
  }

  changeImageTag(imageName: string, imageDetail: ImageDetail) {
    this.imageTagSelected.set(imageName, imageDetail.image_tag);
    this.setActionEnabled();
  }

  actionExecute(): void {
    this.imageList.map(value => value.image_tag = this.imageTagSelected.get(value.image_name));
    this.k8sService.updateServiceImages(this.service.service_project_name, this.service.service_name, this.imageList).subscribe(
      () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_UPDATE_SUCCESSFUL'),
      (err) => this.errorEvent.emit(err)
    );
  }

  setActionEnabled(): void {
    let isEnable = false;
    this.imageList.forEach(value => {
      if (!isEnable) {
        const tag = this.imageTagSelected.get(value.image_name);
        isEnable = tag !== value.image_tag;
      }
    });
    this.actionIsEnabledEvent.emit(isEnable);
  }
}
