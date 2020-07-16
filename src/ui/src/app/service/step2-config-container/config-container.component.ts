import { Component, Injector, OnInit } from '@angular/core';
import {
  Container,
  ContainerType,
  EnvStruct,
  PHASE_CONFIG_CONTAINERS,
  PHASE_CONFIG_INIT_CONTAINERS,
  ServiceStep2Data,
  ServiceStep2DataInit
} from '../service-step.component';
import { concat, Observable, of } from 'rxjs';
import { map } from 'rxjs/operators';
import { RouteImages } from '../../shared/shared.const';
import { AppTokenService } from '../../shared.service/app-token.service';
import { ConfigParamsComponent } from './config-params/config-params.component';
import { ServiceStepComponentBase } from '../service-step';
import { ServiceDockerfileData, ServiceImage, ServiceImageDetail } from '../service.types';

@Component({
  templateUrl: './config-container.component.html',
  styleUrls: ['./config-container.component.css']
})
export class ConfigContainerComponent extends ServiceStepComponentBase implements OnInit {
  imageSourceList: Array<ServiceImage>;
  imageDetailSourceList: Map<string, Array<ServiceImageDetail>>;
  imageTagNotReadyList: Map<Container, boolean>;
  fixedContainerPort: Map<Container, Array<number>>;
  fixedContainerEnv: Map<Container, Array<EnvStruct>>;
  serviceStep2Data: ServiceStep2Data;
  serviceStep2DataInit: ServiceStep2DataInit;
  curContainerType: ContainerType = ContainerType.runContainer;

  constructor(protected injector: Injector,
              private tokenService: AppTokenService) {
    super(injector);
    this.imageSourceList = Array<ServiceImage>();
    this.imageDetailSourceList = new Map<string, Array<ServiceImageDetail>>();
    this.imageTagNotReadyList = new Map<Container, boolean>();
    this.fixedContainerPort = new Map<Container, Array<number>>();
    this.fixedContainerEnv = new Map<Container, Array<EnvStruct>>();
    this.serviceStep2Data = new ServiceStep2Data();
    this.serviceStep2DataInit = new ServiceStep2DataInit();
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS, ServiceStep2Data).subscribe((res: ServiceStep2Data) => {
      this.serviceStep2Data = res;
      this.serviceStep2Data.containerList.forEach((container: Container) => {
        this.getImageDetailList(container).subscribe();
        this.setContainerFixedInfo(container);
      });
      if (this.serviceStep2Data.containerList.length === 0) {
        this.addEmptyWorkItem();
      }
    });
    this.k8sService.getServiceConfig(PHASE_CONFIG_INIT_CONTAINERS, ServiceStep2DataInit).subscribe((res: ServiceStep2DataInit) => {
      this.serviceStep2DataInit = res;
      this.serviceStep2DataInit.containerList.forEach((container: Container) => {
        this.getImageDetailList(container).subscribe();
        this.setContainerFixedInfo(container);
      });
    });
    this.k8sService.getImages('', 0, 0).subscribe(res => this.imageSourceList = res);
  }

  get isCanNextStep(): boolean {
    const runContainerValid = this.serviceStep2Data.containerList
      .filter(value => value.image.imageName !== '')
      .length === this.serviceStep2Data.containerList.length;
    const initContainerValid = this.serviceStep2DataInit.containerList
      .filter(value => value.image.imageName !== '')
      .length === this.serviceStep2DataInit.containerList.length;
    return runContainerValid && initContainerValid;
  }

  get canChangeSelectImageFun() {
    return this.canChangeSelectImage.bind(this);
  }

  get curStep2Data(): ServiceStep2Data {
    return this.curContainerType === ContainerType.runContainer ? this.serviceStep2Data : this.serviceStep2DataInit;
  }

  getActiveImage(index: number): ServiceImage {
    const imageName = this.curStep2Data.containerList[index].image.imageName;
    return this.imageSourceList.find(value => value.imageName === imageName);
  }

  changeContainerType(containerType: ContainerType) {
    this.curContainerType = containerType;
  }

  changeSelectImage(index: number, image: ServiceImage) {
    const container = this.curStep2Data.containerList[index];
    container.name = image.imageName.substr(image.imageName.indexOf('/') + 1);
    container.image.projectName = this.curStep2Data.projectName;
    container.image.imageName = image.imageName;
    if (this.imageDetailSourceList.has(image.imageName)) {
      const detailList: Array<ServiceImageDetail> = this.imageDetailSourceList.get(image.imageName);
      container.image.imageTag = detailList[0].imageTag;
      this.setContainerFixedInfo(container);
    } else {
      this.getImageDetailList(container).subscribe((res: Array<ServiceImageDetail>) => {
        container.image.imageTag = res[0].imageTag;
        this.setContainerFixedInfo(container);
      });
    }
  }

  changeSelectImageDetail(imageName: string, imageDetail: ServiceImageDetail) {
    const container = this.curStep2Data.containerList.find(value => value.image.imageName === imageName);
    container.image.imageTag = imageDetail.imageTag;
    this.setContainerFixedInfo(container);
  }

  getImageDetailList(container: Container): Observable<Array<ServiceImageDetail>> {
    const imageName = container.image.imageName;
    this.imageTagNotReadyList.set(container, false);
    return this.k8sService.getImageDetailList(imageName).pipe(
      map((res: Array<ServiceImageDetail>) => {
        if (res && res.length > 0) {
          for (const item of res) {
            item.imageSizeNumber = Number.parseFloat((item.imageSizeNumber / (1024 * 1024)).toFixed(2));
            item.imageSizeUnit = 'MB';
          }
          this.imageDetailSourceList.set(imageName, res);
        } else {
          this.imageTagNotReadyList.set(container, true);
        }
        return res;
      })
    );
  }

  setContainerFixedInfo(container: Container): void {
    const imageIndex = container.image;
    this.k8sService.getContainerDefaultInfo(imageIndex.imageName, imageIndex.imageTag, imageIndex.projectName).subscribe(
      (res: ServiceDockerfileData) => {
        if (res.imageEnv.length > 0) {
          const fixedEnvs: Array<EnvStruct> = Array<EnvStruct>();
          res.imageEnv.forEach(value => {
            const env = new EnvStruct();
            env.dockerFileEnvName = value.envName;
            env.dockerFileEnvValue = value.envValue;
            if (container.env.find(value1 => value1.dockerFileEnvName === value.envName) === undefined) {
              container.env.push(env);
            }
            fixedEnvs.push(env);
          });
          this.fixedContainerEnv.set(container, fixedEnvs);
        }
        if (res.imageExpose.length > 0) {
          const fixedPorts: Array<number> = Array();
          res.imageExpose.forEach(value => {
            const port: number = Number(value).valueOf();
            fixedPorts.push(port);
            if (container.containerPort.find(value1 => value1 === port) === undefined) {
              container.containerPort.push(port);
            }
          });
          this.fixedContainerPort.set(container, fixedPorts);
        }
      }, () => this.messageService.cleanNotification()
    );
  }

  forward(): void {
    if (this.verifyInputArrayExValid()) {
      const obsRunContainer = this.k8sService.setServiceStepConfig(this.serviceStep2Data);
      const obsInitContainer = this.k8sService.setServiceStepConfig(this.serviceStep2DataInit);
      concat(obsRunContainer, obsInitContainer).subscribe(
        () => this.k8sService.stepSource.next({index: 3, isBack: false})
      );
    }
  }

  canChangeSelectImage(image: ServiceImage): Observable<boolean> {
    if (this.curStep2Data.containerList.find(value => value.image.imageName === image.imageName)) {
      this.messageService.showAlert('IMAGE.CREATE_IMAGE_EXIST', {alertType: 'warning'});
      return of(false);
    }
    return of(true);
  }

  createNewCustomImage() {
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
    const container = new Container();
    this.curStep2Data.containerList.push(container);
  }

  showConfigParams(container: Container): void {
    const component = this.createNewModal(ConfigParamsComponent);
    component.container = container;
    component.fixedContainerEnv = this.fixedContainerEnv;
    component.fixedContainerPort = this.fixedContainerPort;
    component.step2Data = this.curStep2Data;
    component.curContainerType = this.curContainerType;
  }

  backStep() {
    this.k8sService.stepSource.next({index: 1, isBack: true});
  }
}
