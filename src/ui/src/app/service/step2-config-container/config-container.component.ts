import { Component, Injector, OnInit } from '@angular/core';
import {
  Container,
  ContainerType,
  EnvStruct,
  PHASE_CONFIG_CONTAINERS,
  PHASE_CONFIG_INIT_CONTAINERS,
  UiServiceFactory,
  UIServiceStep2,
  VolumeStruct
} from '../service-step.component';
import { BuildImageDockerfileData, Image, ImageDetail } from "../../image/image";
import { ServiceStepBase } from "../service-step";
import { EnvType } from "../../shared/environment-value/environment-value.component";
import { ValidationErrors } from "@angular/forms";
import { NodeAvailableResources } from "../../shared/shared.types";
import { VolumeMountsComponent } from "./volume-mounts/volume-mounts.component";
import { concat, Observable, of } from "rxjs";
import { map } from "rxjs/operators";
import { RouteImages } from "../../shared/shared.const";
import { AppTokenService } from "../../shared.service/app-token.service";

@Component({
  templateUrl: './config-container.component.html',
  styleUrls: ["./config-container.component.css"]
})
export class ConfigContainerComponent extends ServiceStepBase implements OnInit {
  patternContainerName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  patternCpuRequest: RegExp = /^[0-9]*m$/;
  patternCpuLimit: RegExp = /^[0-9]*m$/;
  patternMemRequest: RegExp = /^[0-9]*Mi$/;
  patternMemLimit: RegExp = /^[0-9]*Mi$/;
  imageSourceList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageTagNotReadyList: Map<Container, boolean>;
  containerIsInEdit: Map<Container, boolean>;
  fixedContainerPort: Map<Container, Array<number>>;
  fixedContainerEnv: Map<string, Array<EnvStruct>>;
  serviceStep2Data: UIServiceStep2;
  serviceStep2DataInit: UIServiceStep2;
  showEnvironmentValue = false;
  showVolumeMounts = false;
  curEditEnvContainer: Container;
  curContainerType: ContainerType = ContainerType.runContainer;

  constructor(protected injector: Injector,
              private tokenService: AppTokenService) {
    super(injector);
    this.imageSourceList = Array<Image>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.imageTagNotReadyList = new Map<Container, boolean>();
    this.containerIsInEdit = new Map<Container, boolean>();
    this.fixedContainerPort = new Map<Container, Array<number>>();
    this.fixedContainerEnv = new Map<string, Array<EnvStruct>>();
    this.serviceStep2Data = UiServiceFactory.getInstance(PHASE_CONFIG_CONTAINERS) as UIServiceStep2;
    this.serviceStep2DataInit = UiServiceFactory.getInstance(PHASE_CONFIG_INIT_CONTAINERS) as UIServiceStep2;
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS).subscribe((res: UIServiceStep2) => {
      this.serviceStep2Data = res;
      this.serviceStep2Data.containerList.forEach((container: Container) => {
        this.getImageDetailList(container).subscribe();
        this.containerIsInEdit.set(container, false);
        this.setContainerFixedInfo(container);
      });
      if (this.serviceStep2Data.containerList.length == 0) {
        this.addEmptyWorkItem();
      }
    });
    this.k8sService.getServiceConfig(PHASE_CONFIG_INIT_CONTAINERS).subscribe((res: UIServiceStep2) => {
      this.serviceStep2DataInit = res;
      this.serviceStep2DataInit.containerList.forEach((container: Container) => {
        this.getImageDetailList(container).subscribe();
        this.containerIsInEdit.set(container, false);
        this.setContainerFixedInfo(container);
      });
    });
    this.k8sService.getImages("", 0, 0).subscribe(res => this.imageSourceList = res)
  }

  get isCanNextStep(): boolean {
    const runContainerValid = this.serviceStep2Data.containerList
      .filter(value => value.image.image_name !== "")
      .length == this.serviceStep2Data.containerList.length;
    const initContainerValid = this.serviceStep2DataInit.containerList
      .filter(value => value.image.image_name !== "")
      .length == this.serviceStep2DataInit.containerList.length;
    return runContainerValid && initContainerValid;
  }

  get checkSetCpuRequestFun() {
    return this.checkSetCpuRequest.bind(this);
  }

  get checkSetMemRequestFun() {
    return this.checkSetMemRequest.bind(this);
  }

  get canChangeSelectImageFun() {
    return this.canChangeSelectImage.bind(this);
  }

  get curStep2Data(): UIServiceStep2 {
    return this.curContainerType === ContainerType.runContainer ? this.serviceStep2Data : this.serviceStep2DataInit;
  }

  getActiveImage(index: number): Image {
    const imageName = this.curStep2Data.containerList[index].image.image_name;
    return this.imageSourceList.find(value => value.image_name === imageName);
  }

  changeContainerType(containerType: ContainerType) {
    this.curContainerType = containerType;
  }

  changeSelectImage(index: number, image: Image) {
    let container = this.curStep2Data.containerList[index];
    container.name = image.image_name.substr(image.image_name.indexOf('/') + 1);
    container.image.project_name = this.curStep2Data.projectName;
    container.image.image_name = image.image_name;
    this.toggleContainerEditStatus(container);
    if (this.imageDetailSourceList.has(image.image_name)) {
      let detailList: Array<ImageDetail> = this.imageDetailSourceList.get(image.image_name);
      container.image.image_tag = detailList[0].image_tag;
      this.setDefaultContainerInfo(container);
      this.setContainerFixedInfo(container);
    } else {
      this.getImageDetailList(container).subscribe((res: ImageDetail[]) => {
        container.image.image_tag = res[0].image_tag;
        this.setDefaultContainerInfo(container);
        this.setContainerFixedInfo(container);
      })
    }
  }

  changeSelectImageDetail(imageName: string, imageDetail: ImageDetail) {
    let container = this.curStep2Data.containerList.find(value => value.image.image_name == imageName);
    container.image.image_tag = imageDetail.image_tag;
    this.setDefaultContainerInfo(container);
    this.setContainerFixedInfo(container);
  }

  getImageDetailList(container: Container): Observable<Array<ImageDetail>> {
    const imageName = container.image.image_name;
    this.imageTagNotReadyList.set(container, false);
    return this.k8sService.getImageDetailList(imageName)
      .pipe(map((res: Array<ImageDetail>) => {
        if (res && res.length > 0) {
          for (let item of res) {
            item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
            item['image_size_unit'] = 'MB';
          }
          this.imageDetailSourceList.set(imageName, res);
        } else {
          this.imageTagNotReadyList.set(container, true);
        }
        return res;
      }));
  }

  setContainerFixedInfo(container: Container): void {
    const imageIndex = container.image;
    this.k8sService.getContainerDefaultInfo(imageIndex.image_name, imageIndex.image_tag, imageIndex.project_name).subscribe(
      (res: BuildImageDockerfileData) => {
        if (res.image_env) {
          let fixedEnvs: Array<EnvStruct> = Array<EnvStruct>();
          res.image_env.forEach(value => {
            let env = new EnvStruct();
            env.dockerfile_envname = value.dockerfile_envname;
            env.dockerfile_envvalue = value.dockerfile_envvalue;
            fixedEnvs.push(env);
          });
          this.fixedContainerEnv.set(imageIndex.image_name, fixedEnvs);
        }
        if (res.image_expose) {
          let fixedPorts: Array<number> = Array();
          res.image_expose.forEach(value => {
            let port: number = Number(value).valueOf();
            fixedPorts.push(port);
          });
          this.fixedContainerPort.set(container, fixedPorts);
        }
      }, () => this.messageService.cleanNotification()
    );
  }


  setDefaultContainerInfo(container: Container): void {
    const imageIndex = container.image;
    this.k8sService.getContainerDefaultInfo(imageIndex.image_name, imageIndex.image_tag, imageIndex.project_name).subscribe(
      (res: BuildImageDockerfileData) => {
        if (res.image_cmd) {
          container.command = res.image_cmd;
        }
        if (res.image_env) {
          res.image_env.forEach(value => {
            let env = new EnvStruct();
            env.dockerfile_envname = value.dockerfile_envname;
            env.dockerfile_envvalue = value.dockerfile_envvalue;
            container.env.push(env);
          });
        }
        if (res.image_expose) {
          res.image_expose.forEach(value => {
            let port: number = Number(value).valueOf();
            container.container_port.push(port);
          });
        }
      }, () => this.messageService.cleanNotification());
  }

  isValidContainerNames(step2Data: UIServiceStep2): {invalid: boolean, invalidIndex: number} {
    let invalidIndex: number = -1;
    let everyValid = step2Data.containerList.every((container, index: number) => {
      invalidIndex = index;
      let findRepeat = step2Data.containerList.find((findValue: Container, findIndex: number) =>
        findValue.name == container.name && findIndex != index);
      return this.patternContainerName.test(container.name) && findRepeat == undefined;
    });
    return {invalid: !everyValid, invalidIndex: invalidIndex};
  }

  isValidContainerPorts(step2Data: UIServiceStep2): {invalid: boolean, invalidIndex: number} {
    let invalidIndex: number = -1;
    let valid = true;
    let portBuf = new Set<number>();
    step2Data.containerList.forEach((container, index) => {
      container.container_port.forEach(port => {
        if (portBuf.has(port)) {
          invalidIndex = index;
          valid = false
        } else {
          portBuf.add(port);
        }
      })
    });
    return {invalid: !valid, invalidIndex: invalidIndex};
  }

  isValidContainerCpuAndMem(): {invalid: boolean, invalidIndex: number} {
    let containerList = this.serviceStep2Data.containerList;
    let invalidIndex: number = -1;
    let everyValid = containerList.every((container: Container, index: number) => {
      invalidIndex = index;
      let cpuValid = true;
      let memValid = true;
      if (container.cpu_request != '' && container.cpu_limit != '') {
        cpuValid = Number.parseFloat(container.cpu_request) < Number.parseFloat(container.cpu_limit);
      }
      if (container.mem_request != '' && container.mem_limit != '') {
        memValid = Number.parseFloat(container.mem_request) < Number.parseFloat(container.mem_limit)
      }
      return cpuValid && memValid;
    });
    return {invalid: !everyValid, invalidIndex: invalidIndex};
  }

  forward(): void {
    let funShowInvalidContainer = (invalidIndex: number) => {
      let iterator: IterableIterator<Container> = this.containerIsInEdit.keys();
      let key = iterator.next();
      while (!key.done) {
        this.containerIsInEdit.set(key.value, false);
        key = iterator.next();
      }
      this.containerIsInEdit.set(this.curStep2Data.containerList[invalidIndex], true);
    };
    let checkContainerName = this.isValidContainerNames(this.serviceStep2Data);
    if (checkContainerName.invalid) {
      this.curContainerType = ContainerType.runContainer;
      funShowInvalidContainer(checkContainerName.invalidIndex);
      if (this.verifyInputExValid()) {
        this.messageService.showAlert('SERVICE.STEP_2_CONTAINER_NAME_REPEAT', {alertType: "warning"});
      }
      return;
    }
    let checkInitContainerName = this.isValidContainerNames(this.serviceStep2DataInit);
    if (checkInitContainerName.invalid) {
      this.curContainerType = ContainerType.initContainer;
      funShowInvalidContainer(checkInitContainerName.invalidIndex);
      if (this.verifyInputExValid()) {
        this.messageService.showAlert('SERVICE.STEP_2_CONTAINER_NAME_REPEAT', {alertType: "warning"});
      }
      return;
    }
    let checkContainerPort = this.isValidContainerPorts(this.serviceStep2Data);
    if (checkContainerPort.invalid) {
      this.curContainerType = ContainerType.runContainer;
      funShowInvalidContainer(checkContainerPort.invalidIndex);
      this.messageService.showAlert('SERVICE.STEP_2_CONTAINER_PORT_REPEAT', {alertType: "warning"});
      return;
    }
    let checkInitContainerPort = this.isValidContainerPorts(this.serviceStep2DataInit);
    if (checkInitContainerPort.invalid) {
      this.curContainerType = ContainerType.initContainer;
      funShowInvalidContainer(checkInitContainerPort.invalidIndex);
      this.messageService.showAlert('SERVICE.STEP_2_CONTAINER_PORT_REPEAT', {alertType: "warning"});
      return;
    }
    let checkRequest = this.isValidContainerCpuAndMem();
    if (checkRequest.invalid) {
      this.curContainerType = ContainerType.runContainer;
      funShowInvalidContainer(checkRequest.invalidIndex);
      this.messageService.showAlert('SERVICE.STEP_2_CONTAINER_REQUEST_ERROR', {alertType: "warning"});
      return;
    }
    if (this.verifyInputExValid() && this.verifyInputArrayExValid()) {
      const obsRunContainer = this.k8sService.setServiceConfig(this.serviceStep2Data.uiToServer());
      const obsInitContainer = this.k8sService.setServiceConfig(this.serviceStep2DataInit.uiToServer());
      concat(obsRunContainer, obsInitContainer).subscribe(
        ()=>this.k8sService.stepSource.next({index: 3, isBack: false})
      );
    }
  }

  canChangeSelectImage(image: Image): Observable<boolean> {
    if (this.curStep2Data.containerList.find(value => value.image.image_name == image.image_name)) {
      this.messageService.showAlert('IMAGE.CREATE_IMAGE_EXIST', {alertType: "warning"});
      return of(false);
    }
    return of(true);
  }

  checkSetCpuRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.getNodesAvailableSources()
      .pipe(map((res: Array<NodeAvailableResources>) => {
        let isInValid = res.every(value => Number.parseInt(control.value) > Number.parseInt(value.cpu_available) * 1000);
        if (isInValid) {
          return {beyondMaxLimit: 'SERVICE.STEP_2_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      }));
  }

  checkSetMemRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.getNodesAvailableSources()
      .pipe(map((res: Array<NodeAvailableResources>) => {
        let isInValid = res.every(value => Number.parseInt(control.value) > Number.parseInt(value.mem_available) / (1024 * 1024));
        if (isInValid) {
          return {beyondMaxLimit: 'SERVICE.STEP_2_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      }));
  }

  createNewCustomImage(index: number) {
    this.router.navigate([`/${RouteImages}`], {
        fragment: 'createImage',
        queryParams: {token: this.tokenService.token}
      }
    ).then();
  }

  minusSelectImage(index: number) {
    if (index > 0 || this.curContainerType === ContainerType.initContainer) {
      this.curStep2Data.containerList.splice(index, 1);
    }
  }

  addEmptyWorkItem() {
    let container = new Container();
    this.containerIsInEdit.set(container, false);
    this.curStep2Data.containerList.push(container);
  }

  getVolumesDescription(index: number, container: Container): string {
    let volume = container.volume_mounts;
    if (volume.length > index) {
      let storageServer = volume[index].target_storage_service == "" ? "" : volume[index].target_storage_service.concat(":");
      let result = `${volume[index].container_path}:${storageServer}${volume[index].target_path}`;
      return result == ":" ? "" : result;
    } else {
      return ""
    }
  }

  getEnvsDescription(container: Container): string {
    let envsArr = container.env;
    let result: string = "";
    envsArr.forEach((value: EnvStruct) => {
      result += `${value.dockerfile_envname}=${value.dockerfile_envvalue};`
    });
    return result;
  }

  toggleContainerEditStatus(container: Container): void {
    let oldStatus = this.containerIsInEdit.get(container);
    let iterator: IterableIterator<Container> = this.containerIsInEdit.keys();
    let key = iterator.next();
    while (!key.done) {
      this.containerIsInEdit.set(key.value, false);
      key = iterator.next();
    }
    this.containerIsInEdit.set(container, !oldStatus);
  }

  editEnvironment(container: Container) {
    this.curEditEnvContainer = container;
    this.showEnvironmentValue = true;
  }

  setEnvironment(envsData: Array<EnvType>) {
    let envsArray = this.curEditEnvContainer.env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      let env = new EnvStruct();
      env.dockerfile_envname = value.envName;
      env.dockerfile_envvalue = value.envValue;
      env.configmap_name = value.envConfigMapName;
      env.configmap_key = value.envConfigMapKey;
      envsArray.push(env);
    });
  }

  editVolumeMount(container: Container) {
    this.curEditEnvContainer = container;
    this.showVolumeMounts = true;
    let component = this.createNewModal(VolumeMountsComponent);
    component.volumeDataList = this.curEditEnvContainer.volume_mounts;
    component.onConfirmEvent.subscribe((res: Array<VolumeStruct>) => this.curEditEnvContainer.volume_mounts = res);
  }

  getDefaultEnvsData() {
    let result = Array<EnvType>();
    this.curEditEnvContainer.env.forEach((value: EnvStruct) => {
      let env = new EnvType(value.dockerfile_envname, value.dockerfile_envvalue);
      env.envConfigMapKey = value.configmap_key;
      env.envConfigMapName = value.configmap_name;
      result.push(env)
    });
    return result;
  }

  getDefaultEnvsFixedData(): Array<string> {
    let result = Array<string>();
    if (this.fixedContainerEnv.has(this.curEditEnvContainer.image.image_name)) {
      let fixedEnvs: Array<EnvStruct> = this.fixedContainerEnv.get(this.curEditEnvContainer.image.image_name);
      fixedEnvs.forEach(value => result.push(value.dockerfile_envname));
    }
    return result;
  }

  backStep() {
    this.k8sService.stepSource.next({index: 1, isBack: true});
  }
}
