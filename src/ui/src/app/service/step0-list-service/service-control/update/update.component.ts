import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ImageIndex } from "../../../service-step.component";
import { HttpErrorResponse } from "@angular/common/http";
import { ImageDetail } from "../../../../image/image";
import { K8sService } from "../../../service.k8s";
import { Service } from "../../../service";
import { Message } from "../../../../shared/message-service/message";
import "rxjs/add/observable/of"

@Component({
  selector: 'update',
  templateUrl: './update.component.html',
  styleUrls: ['./update.component.css']
})
export class UpdateComponent implements OnInit {
  @Input() service: Service;
  @Input('isActionInWIP') isActionInWIP: boolean;
  @Output("onMessage") onMessage: EventEmitter<Message>;
  @Output("onError") onError: EventEmitter<any>;
  @Output("onAlertMessage") onAlertMsg: EventEmitter<string>;
  @Output("onActionIsEnabled") onActionIsEnabled: EventEmitter<boolean>;
  imageList: Array<ImageIndex>;
  imageTagList: Map<string, Array<ImageDetail>>;
  imageTagSelected: Map<string, string>;

  constructor(private k8sService: K8sService) {
    this.imageList = Array<ImageIndex>();
    this.onError = new EventEmitter<any>();
    this.onMessage = new EventEmitter<Message>();
    this.onAlertMsg = new EventEmitter<string>();
    this.onActionIsEnabled = new EventEmitter<boolean>();
    this.imageTagList = new Map<string, Array<ImageDetail>>();
    this.imageTagSelected = new Map<string, string>();
  }

  ngOnInit() {
    this.onActionIsEnabled.emit(false);
    this.getServiceImages();
  }

  getServiceImages() {
    this.imageTagSelected.clear();
    this.imageList.splice(0, this.imageList.length);
    this.k8sService.getServiceImages(this.service.service_project_name, this.service.service_name)
      .then(res => {
        this.imageList = res;
        this.imageList.forEach(value => {
          this.imageTagSelected.set(value.image_name, value.image_tag);
          this.k8sService.getImageDetailList(value.image_name)
            .then((res: Array<ImageDetail>) => {
              if (res.length == 0) {
                let tag = new ImageDetail();
                let tagList = Array<ImageDetail>();
                tag.image_tag = value.image_tag;
                tagList.push(tag);
                this.imageTagList.set(value.image_name, tagList);
              } else {
                this.imageTagList.set(value.image_name, res);
              }
              this.setActionEnabled();
            })
            .catch(err => this.onError.next(err));
        });
      })
      .catch((err: HttpErrorResponse) => {
        if (err.status == 500) {
          this.onAlertMsg.emit("SERVICE.SERVICE_CONTROL_NOT_UPDATE");
        } else {
          this.onError.next(err);
        }
      })
  }

  changeImageTag(imageName: string, imageDetail: ImageDetail) {
    this.imageTagSelected.set(imageName, imageDetail.image_tag);
    this.setActionEnabled();
  }

  actionExecute(): void {
    this.imageList.map(value => value.image_tag = this.imageTagSelected.get(value.image_name));
    this.k8sService.updateServiceImages(this.service.service_project_name, this.service.service_name, this.imageList)
      .then(() => {
        let msg: Message = new Message();
        msg.message = "SERVICE.SERVICE_CONTROL_UPDATE_SUCCESSFUL";
        msg.params = [this.service.service_name];
        this.onMessage.emit(msg);
      })
      .catch((err) => this.onError.emit(err));
  }

  setActionEnabled(): void {
    let isEnable: boolean = false;
    this.imageList.forEach(value => {
      if (!isEnable) {
        let tag = this.imageTagSelected.get(value.image_name);
        isEnable = tag != value.image_tag;
      }
    });
    this.onActionIsEnabled.emit(isEnable);
  }
}
