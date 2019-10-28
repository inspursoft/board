/**
 * Created by liyanq on 9/1/17.
 */
import { Component, ComponentFactoryResolver, EventEmitter, Input, OnInit, Output, ViewContainerRef } from "@angular/core"
import { HttpErrorResponse } from "@angular/common/http";
import { CsModalChildBase, CsModalChildMessage } from "../../shared/cs-modal-base/cs-modal-child-base";
import { CreatePvcComponent } from "../../shared/create-pvc/create-pvc.component";
import { JobVolumeMounts } from "../job.type";
import { JobService } from "../job.service";
import { PersistentVolumeClaim } from "../../shared/shared.types";
import { MessageService } from "../../shared.service/message.service";

@Component({
  selector: "job-volume-mounts",
  templateUrl: "./job-volume-mounts.component.html",
  styleUrls: ["./job-volume-mounts.component.css"]
})
export class JobVolumeMountsComponent extends CsModalChildMessage implements OnInit {
  patternName: RegExp = /^[a-z0-9A-Z_]+$/;
  patternMountPath: RegExp = /^[a-z0-9A-Z_/.:]+$/;
  patternIP: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  patternPath: RegExp = /^[a-z0-9A-Z_/.:]+$/;
  volumeTypes: Array<{name: 'nfs' | 'pvc', value: number}>;
  volumeList: Array<JobVolumeMounts>;
  pvcList: Array<PersistentVolumeClaim>;

  @Input() set volumeDataList(value: Array<JobVolumeMounts>) {
    value.forEach(volumeData => {
      let tempVolumeData = new JobVolumeMounts();
      Object.assign(tempVolumeData, volumeData);
      this.volumeList.push(tempVolumeData);
    })
  }

  @Output() onConfirmEvent: EventEmitter<Array<JobVolumeMounts>>;

  constructor(private jobService: JobService,
              protected messageService: MessageService,
              private factoryResolver: ComponentFactoryResolver,
              private selfView: ViewContainerRef) {
    super(messageService);
    this.onConfirmEvent = new EventEmitter<Array<JobVolumeMounts>>();
    this.volumeTypes = Array<{name: 'nfs' | 'pvc', value: number}>();
    this.volumeList = Array<JobVolumeMounts>();
    this.pvcList = Array<PersistentVolumeClaim>();
  }

  ngOnInit() {
    this.volumeTypes.push({name: "nfs", value: 1});
    this.volumeTypes.push({name: "pvc", value: 2});
    this.jobService.getPvcNameList().subscribe(
      (res: Array<PersistentVolumeClaim>) => this.pvcList = res,
      (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {alertType: "warning", view: this.alertView})
    );
  }

  checkInputValid(): boolean {
    let validInput = true;
    this.volumeList.forEach((volume: JobVolumeMounts, index: number) => {
      if (validInput && this.volumeList.find((value, i) =>
        value.volume_name == volume.volume_name &&
        index != i) != undefined) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_NAME', {alertType: "warning", view: this.alertView});
        validInput = false;
      }
      if (validInput && this.volumeList.find((value, i) =>
        value.volume_type === 'nfs' &&
        value.target_path == volume.target_path &&
        index != i) != undefined) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_PATH', {alertType: "warning", view: this.alertView});
        validInput = false;
      }
      if (validInput && this.volumeList.find((value, i) =>
        value.container_path == volume.container_path &&
        index != i) != undefined) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_CONTAINER_PATH', {alertType: "warning", view: this.alertView});
        validInput = false;
      }
    });
    return validInput
  }

  setContainerPathFlag(volume: JobVolumeMounts, event: Event) {
    const checked = (event.target as HTMLInputElement).checked;
    volume.container_path_flag = checked ? 1 : 0;
  }

  confirmVolumeInfo() {
    if (this.verifyInputExValid() && this.checkInputValid()) {
      this.onConfirmEvent.emit(this.volumeList);
      this.modalOpened = false;
    }
  }

  changeSelectPVC(index: number, pvc: PersistentVolumeClaim) {
    this.volumeList[index].target_pvc = pvc.name;
  }

  createNewPvc(index: number) {
    let factory = this.factoryResolver.resolveComponentFactory(CreatePvcComponent);
    let componentRef = this.selfView.createComponent(factory);
    componentRef.instance.openModal().subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    componentRef.instance.onAfterCommit.subscribe((pvc: PersistentVolumeClaim) => {
      this.messageService.cleanNotification();
      this.pvcList.push(pvc);
      this.volumeList[index].target_pvc = pvc.name;
    })
  }

  getCurActivePvc(index: number): PersistentVolumeClaim{
    const pvcName = this.volumeList[index].target_pvc;
    return this.pvcList.find(value => value.name === pvcName);
  }

  changeSelectVolumeType(index: number, volumeType: {name: 'nfs' | 'pvc', value: number}) {
    this.volumeList[index].volume_type = volumeType.name;
  }

  deleteVolumeData(index: number) {
    this.volumeList.splice(index, 1);
  }

  addNewVolumeData() {
    let tempVolumeData = new JobVolumeMounts();
    tempVolumeData.volume_type = 'nfs';
    this.volumeList.push(tempVolumeData);
  }
}
