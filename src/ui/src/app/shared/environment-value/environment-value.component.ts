/**
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core"
import { ValidatorFn, Validators } from "@angular/forms";
import { CsModalChildBase } from "../cs-modal-base/cs-modal-child-base";
import { ConfigMapDetail, ConfigMap } from "../../resource/resource.types";
import { ResourceService } from "../../resource/resource.service";
import { HttpErrorResponse } from "@angular/common/http";
import { MessageService } from "../../shared.service/message.service";

export class EnvType {
  public envName = '';
  public envValue = '';
  public envConfigMapName = '';
  public envConfigMapKey = '';

  constructor(name, value: string) {
    this.envName = name.trim();
    this.envValue = value.trim();
  }
}

@Component({
  selector: "environment-value",
  templateUrl: "./environment-value.component.html",
  styleUrls: ["./environment-value.component.css"],
  providers: [ResourceService]
})
export class EnvironmentValueComponent extends CsModalChildBase implements OnInit {
  patternEnv = /^[\w-$/\\=\"[\]{}@&:,'`\t. ?]+$/;
  envsData: Array<EnvType>;
  envsText = "";
  inputValidator: Array<ValidatorFn>;
  configMapList: Array<ConfigMap>;
  configMapDetail: Map<number, ConfigMapDetail>;
  bindConfigMap: Map<number, boolean>;

  @Input() inputEnvsData: Array<EnvType>;
  @Input() inputFixedKeyList: Array<string>;
  @Input() isProvideBindConfigMap = false;
  @Input() projectName = '';
  @Output() onConfirm: EventEmitter<Array<EnvType>>;

  constructor(private messageService: MessageService,
              private resourceService: ResourceService) {
    super();
    this.envsData = Array<EnvType>();
    this.onConfirm = new EventEmitter<Array<EnvType>>();
    this.inputValidator = Array<ValidatorFn>();
    this.configMapList = Array<ConfigMap>();
    this.configMapDetail = new Map<number, ConfigMapDetail>();
    this.bindConfigMap = new Map<number, boolean>();
  }

  ngOnInit() {
    this.inputValidator.push(Validators.required);
    if (this.inputEnvsData && this.inputEnvsData.length > 0) {
      this.envsData = this.envsData.concat(this.inputEnvsData);
      this.envsData.forEach((value: EnvType, index: number) => {
        let detail = new ConfigMapDetail();
        detail.dataList.push({key: value.envConfigMapKey, value: value.envConfigMapKey})
        this.configMapDetail.set(index, detail);
        this.bindConfigMap.set(index, value.envConfigMapKey != '')
      })
    }
    this.resourceService.getConfigMapList(this.projectName, 0, 0).subscribe(
      (res: Array<ConfigMap>) => this.configMapList = res);
    this.modalOpened = true;
  }

  addNewEnv() {
    this.envsData.push(new EnvType("", ""));
  }

  confirmEnvInfo() {
    if (this.verifyInputExValid() && this.verifyDropdownExValid()) {
      this.onConfirm.emit(this.envsData);
      this.modalOpened = false;
    }
  }

  envMinusClick(index: number) {
    this.envsData.splice(index, 1);
    this.bindConfigMap.delete(index);
    this.configMapDetail.delete(index);
  }

  envTextAddClick() {
    try {
      let envTypes = this.envsText.split(";").map((str: string) => {
        let envStrPair = str.split("=");
        if (!this.patternEnv.test(envStrPair[0]) || !this.patternEnv.test(envStrPair[1])) {
          throw new Error()
        }
        return new EnvType(envStrPair[0], envStrPair[1]);
      });
      this.envsData = this.envsData.concat(envTypes);
    } catch (e) {
      this.messageService.showAlert('SERVICE.TXT_ALERT_MESSAGE', {alertType: 'warning', view: this.alertView});
      return;
    }
  }

  changeConfigMap(index: number, envInfo: EnvType, configMap: ConfigMap) {
    envInfo.envConfigMapName = configMap.name;
    this.resourceService.getConfigMapDetail(configMap.name, this.projectName).subscribe(
      (res: ConfigMapDetail) => this.configMapDetail.set(index, res),
      (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {alertType: "danger", view: this.alertView})
    )
  }

  changeConfigMapKey(envInfo: EnvType, data: {key: string, value: string}) {
    envInfo.envValue = data.value;
    envInfo.envConfigMapKey = data.key;
  }

  changeBindConfigMap(index: number, envInfo: EnvType, event: Event) {
    const checked = (event.target as HTMLInputElement).checked;
    this.bindConfigMap.set(index, checked);
    if (!checked) {
      envInfo.envConfigMapName = '';
      envInfo.envConfigMapKey = '';
    }
  }
}
