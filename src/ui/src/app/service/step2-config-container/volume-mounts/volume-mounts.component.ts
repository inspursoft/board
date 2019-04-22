/**
 * Created by liyanq on 9/1/17.
 */
import { Component, ComponentFactoryResolver, EventEmitter, Input, OnInit, Output, ViewContainerRef } from "@angular/core"
import { VolumeStruct } from "../../service-step.component";
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";
import { K8sService } from "../../service.k8s";
import { PersistentVolumeClaim } from "../../../shared/shared.types";
import { MessageService } from "../../../shared.service/message.service";
import { HttpErrorResponse } from "@angular/common/http";
import { CreatePvcComponent } from "../../../shared/create-pvc/create-pvc.component";

@Component({
  selector: "volume-mounts",
  templateUrl: "./volume-mounts.component.html",
  styleUrls: ["./volume-mounts.component.css"]
})
export class VolumeMountsComponent extends CsModalChildBase implements OnInit {
  patternName: RegExp = /^[a-z0-9A-Z_]+$/;
  patternMountPath: RegExp = /^[a-z0-9A-Z_/]+$/;
  patternIP: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  patternPath: RegExp = /^[a-z0-9A-Z_/.:]+$/;
  volumeTypes: Array<{name: 'nfs' | 'pvc', value: number}>;
  curVolumeDataList: Array<VolumeStruct>;
  pvcList: Array<PersistentVolumeClaim>;

  @Input() set volumeDataList(value: Array<VolumeStruct>) {
    value.forEach(volumeData => {
      let tempVolumeData = new VolumeStruct();
      Object.assign(tempVolumeData, volumeData);
      this.curVolumeDataList.push(tempVolumeData);
    })
  }

  @Output() onConfirmEvent: EventEmitter<Array<VolumeStruct>>;

  constructor(private k8sService: K8sService,
              private messageService: MessageService,
              private factoryResolver: ComponentFactoryResolver,
              private selfView: ViewContainerRef) {
    super();
    this.onConfirmEvent = new EventEmitter<Array<VolumeStruct>>();
    this.volumeTypes = Array<{name: 'nfs' | 'pvc', value: number}>();
    this.curVolumeDataList = Array<VolumeStruct>();
    this.pvcList = Array<PersistentVolumeClaim>();
  }

  ngOnInit() {
    this.volumeTypes.push({name: "nfs", value: 1});
    this.volumeTypes.push({name: "pvc", value: 2});
    let createNewPVC = new PersistentVolumeClaim();
    createNewPVC.name = "SERVICE.VOLUME_CREATE_PVC";
    createNewPVC["isSpecial"] = true;
    createNewPVC["OnlyClick"] = true;
    this.pvcList.push(createNewPVC);
    this.k8sService.getPvcNameList().subscribe((res: Array<PersistentVolumeClaim>) => {
        if (res && res.length > 0) {
          this.pvcList = this.pvcList.concat(res);
        }
      }, (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {alertType: "warning", view: this.alertView})
    )
  }

  pvcDropdownDefaultText(index: number): string {
    return this.curVolumeDataList[index].target_pvc == '' ? 'SERVICE.VOLUME_SELECT_PVC' : this.curVolumeDataList[index].target_pvc;
  }

  checkInputValid(): boolean {
    let validInput = true;
    this.curVolumeDataList.forEach((volume: VolumeStruct, index: number) => {
      if (this.curVolumeDataList.find((value, i) => value.volume_name == volume.volume_name && index != i) != undefined && validInput) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_NAME', {alertType: "warning", view: this.alertView});
        validInput = false;
      }
      if (this.curVolumeDataList.find((value, i) => value.target_path == volume.target_path && index != i) != undefined && validInput) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_PATH', {alertType: "warning", view: this.alertView});
        validInput = false;
      }
      if (this.curVolumeDataList.find((value, i) => value.container_path == volume.container_path && index != i) != undefined && validInput) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_CONTAINER_PATH', {alertType: "warning", view: this.alertView});
        validInput = false;
      }
    });
    return validInput
  }

  confirmVolumeInfo() {
    if (this.verifyInputValid() && this.checkInputValid()) {
      this.onConfirmEvent.emit(this.curVolumeDataList);
      this.modalOpened = false;
    }
  }

  changeSelectPVC(index: number, pvc: PersistentVolumeClaim) {
    this.curVolumeDataList[index].target_pvc = pvc.name;
  }

  createNewPvc(index: number) {
    let factory = this.factoryResolver.resolveComponentFactory(CreatePvcComponent);
    let componentRef = this.selfView.createComponent(factory);
    componentRef.instance.openModal().subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    componentRef.instance.onAfterCommit.subscribe((pvc: PersistentVolumeClaim) => {
      this.messageService.cleanNotification();
      this.curVolumeDataList[index].target_pvc = pvc.name;
      this.pvcList.push(pvc)
    })
  }

  changeSelectVolumeType(index: number, volumeType: {name: 'nfs' | 'pvc', value: number}) {
    this.curVolumeDataList[index].volume_type = volumeType.name;
  }

  deleteVolumeData(index: number) {
    this.curVolumeDataList.splice(index, 1);
  }

  addNewVolumeData() {
    let tempVolumeData = new VolumeStruct();
    this.curVolumeDataList.push(tempVolumeData);
  }
}
