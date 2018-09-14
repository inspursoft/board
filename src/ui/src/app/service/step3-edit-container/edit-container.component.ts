import { AfterContentChecked, Component, Injector, OnInit, QueryList, ViewChildren } from '@angular/core';
import { Container, EnvStruct, PHASE_CONFIG_CONTAINERS, ServiceStepPhase, UIServiceStep3 } from '../service-step.component';
import { EnvType } from "../../shared/environment-value/environment-value.component";
import { CsInputArrayComponent } from "../../shared/cs-components-library/cs-input-array/cs-input-array.component";
import { VolumeOutPut } from "./volume-mounts/volume-mounts.component";
import { ServiceStepBase } from "../service-step";
import { BuildImageDockerfileData } from "../../image/image";

@Component({
  templateUrl: './edit-container.component.html',
  styleUrls: ["./edit-container.component.css"]
})
export class EditContainerComponent extends ServiceStepBase implements OnInit, AfterContentChecked {
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  patternContainerName: RegExp = /^[a-zA-Z\d_-]+$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  step3TypeStatus: Map<Container, boolean>;
  showVolumeMounts = false;
  showEnvironmentValue = false;
  isInputComponentsValid = false;
  fixedEnvKeys: Array<string>;
  fixedContainerPort: Map<string, Array<number>>;
  curContainerIndex: number;

  constructor(protected injector: Injector) {
    super(injector);
    this.step3TypeStatus = new Map<Container, boolean>();
    this.fixedEnvKeys = Array<string>();
    this.fixedContainerPort = new Map<string, Array<number>>();
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(this.stepPhase).then(res => {
      this.uiBaseData = res;
      this.uiData.containerList.forEach((container: Container) => {
        this.step3TypeStatus.set(container, false);
        this.setDefaultContainerInfo(container);
      });
    });
  }

  ngAfterContentChecked() {
    this.isInputComponentsValid = true;
    if (this.inputArrayComponents) {
      this.inputArrayComponents.forEach(item => {
        if (!item.valid) {
          this.isInputComponentsValid = false;
        }
      });
    }
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_CONFIG_CONTAINERS;
  }

  get uiData(): UIServiceStep3 {
    return this.uiBaseData as UIServiceStep3;
  }

  setDefaultContainerInfo(container: Container): void {
    let isNew = !this.isBack;
    this.k8sService.getContainerDefaultInfo(container.image.image_name, container.image.image_tag, container.image.project_name)
      .then((res: BuildImageDockerfileData) => {
        this.step3TypeStatus.set(container, true);
        if (res.image_cmd && isNew) {
          container.command = res.image_cmd;//copy cmd
        }
        if (res.image_env) {
          res.image_env.forEach(value => {//copy env
            this.fixedEnvKeys.push(value.dockerfile_envname);
            if (isNew) {
              let env = new EnvStruct();
              env.dockerfile_envname = value.dockerfile_envname;
              env.dockerfile_envvalue = value.dockerfile_envvalue;
              container.env.push(env);
            }
          });
        }
        if (res.image_expose) {
          let fixedPorts: Array<number> = Array();
          res.image_expose.forEach(value => {//copy port
            let port: number = Number(value).valueOf();
            fixedPorts.push(port);
            container.container_port.push(port);
          });
          this.fixedContainerPort.set(container.name, fixedPorts);
        }
      }).catch(() => this.messageService.cleanNotification());
  }

  get isCanNextStep(): boolean {
    return this.uiData.containerList.length > 0 && this.isInputComponentsValid;
  }

  getVolumesDescription(container: Container): string {
    let volume = container.volume_mount;
    let storageServer = volume.target_storage_service == "" ? "" : volume.target_storage_service.concat(":");
    let result = `${volume.container_path}:${storageServer}${volume.target_path}`;
    return result == ":" ? "" : result;
  }

  getEnvsDescription(container: Container): string {
    let envsArr = container.env;
    let result: string = "";
    envsArr.forEach((value: EnvStruct) => {
      result += `${value.dockerfile_envname}=${value.dockerfile_envvalue};`
    });
    return result;
  }

  getDefaultEnvsData(index: number) {
    let result = Array<EnvType>();
    this.uiData.containerList[index].env.forEach((value: EnvStruct) => {
      result.push(new EnvType(value.dockerfile_envname, value.dockerfile_envvalue))
    });
    return result;
  }

  setEnvironment(index: number, envsData: Array<EnvType>) {
    let envsArray = this.uiData.containerList[index].env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      let env = new EnvStruct();
      env.dockerfile_envname = value.envName;
      env.dockerfile_envvalue = value.envValue;
      envsArray.push(env);
    });
  }

  setVolumeMount(data: VolumeOutPut, index: number) {
    let volume = this.uiData.containerList[index].volume_mount;
    volume.target_storage_service = data.out_medium;
    volume.target_path = data.out_path;
    volume.container_path = data.out_mountPath;
    volume.volume_name = data.out_name;
  }

  getVolumeMountData(index: number): VolumeOutPut {
    let volume = this.uiData.containerList[index].volume_mount;
    return {
      out_name: volume.volume_name,
      out_mountPath: volume.container_path,
      out_path: volume.target_path,
      out_medium: volume.target_storage_service
    };
  }

  toggleShowStatus(container: Container): void {
    let status = this.step3TypeStatus.get(container);
    this.step3TypeStatus.set(container, !status);
  }

  shieldEnter($event: KeyboardEvent) {
    if ($event.charCode == 13) {
      (<any>$event.target).blur();
    }
  }

  backStep(): void {
    this.k8sService.stepSource.next({index: 2, isBack: true});
  }

  forward(): void {
    if (this.verifyInputValid()) {
      this.k8sService.setServiceConfig(this.uiData.uiToServer()).then(() =>
        this.k8sService.stepSource.next({index: 4, isBack: false})
      )
    }
  }
}