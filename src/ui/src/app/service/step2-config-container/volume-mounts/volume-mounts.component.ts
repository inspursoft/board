/**
 * Created by liyanq on 9/1/17.
 */
import {
  Component,
  ComponentFactoryResolver,
  EventEmitter,
  Input, OnDestroy,
  OnInit,
  ViewContainerRef
} from '@angular/core';
import { Volume, VolumeType } from '../../service-step.component';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { K8sService } from '../../service.k8s';
import { PersistentVolumeClaim } from '../../../shared/shared.types';
import { MessageService } from '../../../shared.service/message.service';
import { HttpErrorResponse } from '@angular/common/http';
import { CreatePvcComponent } from '../../../shared/create-pvc/create-pvc.component';

@Component({
  templateUrl: './volume-mounts.component.html',
  styleUrls: ['./volume-mounts.component.css']
})
export class VolumeMountsComponent extends CsModalChildBase implements OnInit, OnDestroy {
  patternName: RegExp = /^[a-z0-9A-Z_]+$/;
  patternMountPath: RegExp = /^[a-z0-9A-Z_/]+$/;
  patternIP: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  patternPath: RegExp = /^[a-z0-9A-Z_/.:]+$/;
  volumeTypes: Array<{ type: VolumeType, value: number }>;
  curVolumeDataList: Array<Volume>;
  pvcList: Array<PersistentVolumeClaim>;
  configMapList: Array<string>;
  onConfirmEvent: EventEmitter<Array<Volume>>;
  projectName = '';

  constructor(private k8sService: K8sService,
              private messageService: MessageService,
              private factoryResolver: ComponentFactoryResolver,
              private selfView: ViewContainerRef) {
    super();
    this.onConfirmEvent = new EventEmitter<Array<Volume>>();
    this.volumeTypes = Array<{ type: VolumeType, value: number }>();
    this.curVolumeDataList = Array<Volume>();
    this.pvcList = Array<PersistentVolumeClaim>();
    this.configMapList = new Array<string>();
  }

  ngOnInit() {
    this.volumeTypes.push({type: 'nfs', value: 1});
    this.volumeTypes.push({type: 'pvc', value: 2});
    this.volumeTypes.push({type: 'configmap', value: 3});
    this.k8sService.getPvcNameList().subscribe(
      (res: Array<PersistentVolumeClaim>) => this.pvcList = res,
      (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {
        alertType: 'warning',
        view: this.alertView
      })
    );
    this.k8sService.getConfigMapNames(this.projectName).subscribe(
      res => this.configMapList = res,
      (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {
        alertType: 'warning',
        view: this.alertView
      })
    );
  }

  ngOnDestroy() {
    this.onConfirmEvent.unsubscribe();
    delete this.onConfirmEvent;
    super.ngOnDestroy();
  }

  @Input() set volumeDataList(value: Array<Volume>) {
    value.forEach(volumeData => {
      const tempVolumeData = new Volume();
      Object.assign(tempVolumeData, volumeData);
      this.curVolumeDataList.push(tempVolumeData);
    });
  }

  getCurActivePvc(index: number): PersistentVolumeClaim {
    const pvcName = this.curVolumeDataList[index].targetPvc;
    return this.pvcList.find(value => value.name === pvcName);
  }

  getCurActiveConfigMap(index: number): string {
    return this.curVolumeDataList[index].targetConfigMap;
  }

  checkInputValid(): boolean {
    let validInput = true;
    this.curVolumeDataList.forEach((volume: Volume, index: number) => {
      if (this.curVolumeDataList.find((value, i) =>
        value.volumeName === volume.volumeName && index !== i) !== undefined && validInput) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_NAME', {alertType: 'warning', view: this.alertView});
        validInput = false;
      }
      if (this.curVolumeDataList.find((value, i) =>
        value.targetPath === volume.targetPath && index !== i && value.targetPath !== '') !== undefined && validInput) {
        this.messageService.showAlert('SERVICE.NFS_STORAGE_VALID_PATH', {alertType: 'warning', view: this.alertView});
        validInput = false;
      }
      if (this.curVolumeDataList.find((value, i) =>
        value.containerPath === volume.containerPath && index !== i) !== undefined && validInput) {
        this.messageService.showAlert('SERVICE.VOLUME_VALID_CONTAINER_PATH', {
          alertType: 'warning',
          view: this.alertView
        });
        validInput = false;
      }
    });
    return validInput;
  }

  confirmVolumeInfo() {
    if (this.verifyInputExValid() && this.checkInputValid() && this.verifyDropdownExValid()) {
      this.onConfirmEvent.emit(this.curVolumeDataList);
      this.modalOpened = false;
    }
  }

  changeSelectPVC(index: number, pvc: PersistentVolumeClaim) {
    this.curVolumeDataList[index].targetPvc = pvc.name;
  }

  createNewPvc(index: number) {
    const factory = this.factoryResolver.resolveComponentFactory(CreatePvcComponent);
    const componentRef = this.selfView.createComponent(factory);
    componentRef.instance.openModal().subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    componentRef.instance.onAfterCommit.subscribe((pvc: PersistentVolumeClaim) => {
      this.messageService.cleanNotification();
      this.curVolumeDataList[index].targetPvc = pvc.name;
      this.pvcList.push(pvc);
    });
  }

  changeSelectVolumeType(index: number, volumeType: { type: VolumeType, value: number }) {
    this.curVolumeDataList[index].volumeType = volumeType.type;
  }

  deleteVolumeData(index: number) {
    this.curVolumeDataList.splice(index, 1);
  }

  addNewVolumeData() {
    const tempVolumeData = new Volume();
    this.curVolumeDataList.push(tempVolumeData);
  }
}
