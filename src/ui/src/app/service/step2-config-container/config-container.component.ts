import { Component, Injector, OnInit } from '@angular/core';
import {
  Container,
  ContainerType,
  EnvStruct,
  PHASE_CONFIG_CONTAINERS,
  PHASE_CONFIG_INIT_CONTAINERS,
  UiServiceFactory,
  UIServiceStep2
} from '../service-step.component';
import { BuildImageDockerfileData, Image, ImageDetail } from '../../image/image';
import { ServiceStepBase } from '../service-step';
import { concat, Observable, of } from 'rxjs';
import { map } from 'rxjs/operators';
import { RouteImages } from '../../shared/shared.const';
import { AppTokenService } from '../../shared.service/app-token.service';
import { ConfigParamsComponent } from "./config-params/config-params.component";

@Component({
  templateUrl: './config-container.component.html',
  styleUrls: ['./config-container.component.css']
})
export class ConfigContainerComponent extends ServiceStepBase implements OnInit {
  imageSourceList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageTagNotReadyList: Map<Container, boolean>;
  fixedContainerPort: Map<Container, Array<number>>;
  fixedContainerEnv: Map<string, Array<EnvStruct>>;
  serviceStep2Data: UIServiceStep2;
  serviceStep2DataInit: UIServiceStep2;
  curContainerType: ContainerType = ContainerType.runContainer;

  constructor(protected injector: Injector,
              private tokenService: AppTokenService) {
    super(injector);
    this.imageSourceList = Array<Image>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.imageTagNotReadyList = new Map<Container, boolean>();
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
        this.setContainerFixedInfo(container);
      });
      if (this.serviceStep2Data.containerList.length === 0) {
        this.addEmptyWorkItem();
      }
    });
    this.k8sService.getServiceConfig(PHASE_CONFIG_INIT_CONTAINERS).subscribe((res: UIServiceStep2) => {
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
      .filter(value => value.image.image_name !== '')
      .length === this.serviceStep2Data.containerList.length;
    const initContainerValid = this.serviceStep2DataInit.containerList
      .filter(value => value.image.image_name !== '')
      .length === this.serviceStep2DataInit.containerList.length;
    return runContainerValid && initContainerValid;
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
    const container = this.curStep2Data.containerList[index];
    container.name = image.image_name.substr(image.image_name.indexOf('/') + 1);
    container.image.project_name = this.curStep2Data.projectName;
    container.image.image_name = image.image_name;
    if (this.imageDetailSourceList.has(image.image_name)) {
      const detailList: Array<ImageDetail> = this.imageDetailSourceList.get(image.image_name);
      container.image.image_tag = detailList[0].image_tag;
      this.setDefaultContainerInfo(container);
      this.setContainerFixedInfo(container);
    } else {
      this.getImageDetailList(container).subscribe((res: ImageDetail[]) => {
        container.image.image_tag = res[0].image_tag;
        this.setDefaultContainerInfo(container);
        this.setContainerFixedInfo(container);
      });
    }
  }

  changeSelectImageDetail(imageName: string, imageDetail: ImageDetail) {
    const container = this.curStep2Data.containerList.find(value => value.image.image_name === imageName);
    container.image.image_tag = imageDetail.image_tag;
    this.setDefaultContainerInfo(container);
    this.setContainerFixedInfo(container);
  }

  getImageDetailList(container: Container): Observable<Array<ImageDetail>> {
    const imageName = container.image.image_name;
    this.imageTagNotReadyList.set(container, false);
    return this.k8sService.getImageDetailList(imageName).pipe(
      map((res: Array<ImageDetail>) => {
        if (res && res.length > 0) {
          for (const item of res) {
            item.image_size_number = Number.parseFloat((item.image_size_number / (1024 * 1024)).toFixed(2));
            item.image_size_unit = 'MB';
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
    this.k8sService.getContainerDefaultInfo(imageIndex.image_name, imageIndex.image_tag, imageIndex.project_name).subscribe(
      (res: BuildImageDockerfileData) => {
        if (res.image_env) {
          const fixedEnvs: Array<EnvStruct> = Array<EnvStruct>();
          res.image_env.forEach(value => {
            const env = new EnvStruct();
            env.dockerfile_envname = value.dockerfile_envname;
            env.dockerfile_envvalue = value.dockerfile_envvalue;
            fixedEnvs.push(env);
          });
          this.fixedContainerEnv.set(imageIndex.image_name, fixedEnvs);
        }
        if (res.image_expose) {
          const fixedPorts: Array<number> = Array();
          res.image_expose.forEach(value => {
            const port: number = Number(value).valueOf();
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
            const env = new EnvStruct();
            env.dockerfile_envname = value.dockerfile_envname;
            env.dockerfile_envvalue = value.dockerfile_envvalue;
            container.env.push(env);
          });
        }
        if (res.image_expose) {
          res.image_expose.forEach(value => {
            const port: number = Number(value).valueOf();
            container.container_port.push(port);
          });
        }
      }, () => this.messageService.cleanNotification());
  }

  forward(): void {
    if (this.verifyInputArrayExValid()) {
      const obsRunContainer = this.k8sService.setServiceConfig(this.serviceStep2Data.uiToServer());
      const obsInitContainer = this.k8sService.setServiceConfig(this.serviceStep2DataInit.uiToServer());
      concat(obsRunContainer, obsInitContainer).subscribe(
        () => this.k8sService.stepSource.next({index: 3, isBack: false})
      );
    }
  }

  canChangeSelectImage(image: Image): Observable<boolean> {
    if (this.curStep2Data.containerList.find(value => value.image.image_name === image.image_name)) {
      this.messageService.showAlert('IMAGE.CREATE_IMAGE_EXIST', {alertType: 'warning'});
      return of(false);
    }
    return of(true);
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
    const container = new Container();
    this.curStep2Data.containerList.push(container);
  }

  showConfigParams(container: Container): void {
    const component = this.createNewModal(ConfigParamsComponent);
    component.container = container.clone();
    component.fixedContainerEnv = this.fixedContainerEnv;
    component.fixedContainerPort = this.fixedContainerPort;
    component.step2Data = this.curStep2Data;
    component.curContainerType = this.curContainerType;
  }

  backStep() {
    this.k8sService.stepSource.next({index: 1, isBack: true});
  }
}
