import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ImageIndex } from '../../../service-step.component';
import { K8sService } from '../../../service.k8s';
import { Service, ServiceImage, ServiceImageDetail } from '../../../service.types';

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
  imageList: Array<ServiceImage>;
  imageTagList: Map<string, Array<ServiceImageDetail>>;
  imageTagSelected: Map<string, string>;

  constructor(private k8sService: K8sService) {
    this.imageList = Array<ServiceImage>();
    this.errorEvent = new EventEmitter<any>();
    this.messageEvent = new EventEmitter<string>();
    this.alertEvent = new EventEmitter<string>();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
    this.imageTagList = new Map<string, Array<ServiceImageDetail>>();
    this.imageTagSelected = new Map<string, string>();
  }

  ngOnInit() {
    this.actionIsEnabledEvent.emit(false);
    this.getServiceImages();
  }

  getServiceImages() {
    this.imageTagSelected.clear();
    this.imageList.splice(0, this.imageList.length);
    this.k8sService.getServiceImages(this.service.serviceProjectName, this.service.serviceName).subscribe(
      (imageList: Array<ServiceImage>) => {
        this.imageList = imageList;
        this.imageList.forEach(value => {
          this.imageTagSelected.set(value.imageName, value.imageTag);
          this.k8sService.getImageDetailList(value.imageName).subscribe(
            (imageDetailList: Array<ServiceImageDetail>) => {
              if (imageDetailList.length === 0) {
                const tag = new ServiceImageDetail();
                const tagList = Array<ServiceImageDetail>();
                tag.imageTag = value.imageTag;
                tagList.push(tag);
                this.imageTagList.set(value.imageName, tagList);
              } else {
                this.imageTagList.set(value.imageName, imageDetailList);
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

  changeImageTag(imageName: string, imageDetail: ServiceImageDetail) {
    this.imageTagSelected.set(imageName, imageDetail.imageTag);
    this.setActionEnabled();
  }

  actionExecute(): void {
    this.imageList.map(value => value.imageTag = this.imageTagSelected.get(value.imageName));
    const postBody = new Array<{ [key: string]: string }>();
    this.imageList.forEach(value => {
      const imageIndex = new ImageIndex();
      imageIndex.imageName = value.imageName;
      imageIndex.imageTag = value.imageTag;
      imageIndex.projectName = this.service.serviceProjectName;
      postBody.push(imageIndex.getPostBody());
    });
    this.k8sService.updateServiceImages(this.service.serviceProjectName, this.service.serviceName, postBody).subscribe(
      () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_UPDATE_SUCCESSFUL'),
      (err) => this.errorEvent.emit(err)
    );
  }

  setActionEnabled(): void {
    let isEnable = false;
    this.imageList.forEach(value => {
      if (!isEnable) {
        const tag = this.imageTagSelected.get(value.imageName);
        isEnable = tag !== value.imageTag;
      }
    });
    this.actionIsEnabledEvent.emit(isEnable);
  }
}
