/**
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, OnInit, Output, QueryList, ViewChildren } from "@angular/core"
import { CsInputComponent } from "../cs-components-library/cs-input/cs-input.component";
import { ValidatorFn, Validators } from "@angular/forms";
import { CsModalChildBase } from "../cs-modal-base/cs-modal-child-base";
import { MessageService } from "../message-service/message.service";
import { ConfigMapDetail, ConfigMapList } from "../../resource/resource.types";
import { ResourceService } from "../../resource/resource.service";
import { HttpErrorResponse } from "@angular/common/http";

export class EnvType {
  public envName = '';
  public envValue = '';
  public envConfigMapName = '';
  public envConfigMapKey = '';

  constructor(name, value: string) {
    this.envName = name;
    this.envValue = value;
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
  inputValidatorMsg: Array<{validatorKey: string, validatorMessage: string}>;
  configMapList: Array<ConfigMapList>;
  configMapDetail: Map<number, ConfigMapDetail>;
  bindConfigMap: Map<number, boolean>;

  @ViewChildren(CsInputComponent) inputComponentList: QueryList<CsInputComponent>;
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
    this.inputValidatorMsg = Array<{validatorKey: string, validatorMessage: string}>();
    this.configMapList = Array<ConfigMapList>();
    this.configMapDetail = new Map<number, ConfigMapDetail>();
    this.bindConfigMap = new Map<number, boolean>();
  }

  ngOnInit() {
    this.inputValidator.push(Validators.required);
    this.inputValidatorMsg.push({validatorKey: "required", validatorMessage: "SERVICE.ENV_REQUIRED"});
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
      (res: Array<ConfigMapList>) => this.configMapList = res);
    this.modalOpened = true;
  }

  getEnvConfigMapDefaltText(index: number): string {
    return this.envsData[index].envConfigMapName != '' ?
      this.envsData[index].envConfigMapName :
      'SERVICE.ENV_CONFIG_MAP_DEFAULT_TEXT'
  }

  getEnvConfigMapDefaltKeyText(index: number): string {
    return this.envsData[index].envConfigMapKey != '' ?
      this.envsData[index].envConfigMapKey :
      'SERVICE.ENV_CONFIG_MAP_KEY_DEFAULT_TEXT'
  }

  addNewEnv() {
    this.envsData.push(new EnvType("", ""));
  }

  confirmEnvInfo() {
    if (this.verifyInputValid() && this.verifyDropdownValid()) {
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
    let patternEnv = this.patternEnv;
    let envTypes: Array<EnvType>;
    try {
      envTypes = this.envsText.split(";").map((str: string) => {
        let envStrPair = str.split("=");
        if (!patternEnv.test(envStrPair[0]) || !patternEnv.test(envStrPair[1])) {
          throw new Error()
        }
        return new EnvType(envStrPair[0], envStrPair[1]);
      });
    } catch (e) {
      this.messageService.showAlert('SERVICE.TXT_ALERT_MESSAGE', {alertType: 'alert-warning', view: this.alertView});
      return;
    }
    this.envsData = this.envsData.concat(envTypes);
  }

  changeConfigMap(index: number, envInfo: EnvType, configMap: ConfigMapList) {
    envInfo.envConfigMapName = configMap.name;
    this.resourceService.getConfigMapDetail(configMap.name, this.projectName).subscribe(
      (res: ConfigMapDetail) => this.configMapDetail.set(index, res),
      (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {alertType: "alert-danger", view: this.alertView})
    )
  }

  changeConfigMapKey(envInfo: EnvType, data: {key: string, value: string}) {
    envInfo.envValue = data.value;
    envInfo.envConfigMapKey = data.key;
  }

  changeBindConfigMap(index: number, envInfo: EnvType, checked: boolean) {
    this.bindConfigMap.set(index, checked);
    if (!checked) {
      envInfo.envConfigMapName = '';
      envInfo.envConfigMapKey = '';
    }
  }
}
