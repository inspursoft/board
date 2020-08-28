/**
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ValidatorFn, Validators } from '@angular/forms';
import { CsModalChildBase } from '../cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../shared.service/message.service';
import { SharedConfigMap, SharedConfigMapDetail, SharedEnvType } from '../shared.types';
import { SharedService } from '../../shared.service/shared.service';

@Component({
  selector: 'app-environment-value',
  templateUrl: './environment-value.component.html',
  styleUrls: ['./environment-value.component.css'],
})
export class EnvironmentValueComponent extends CsModalChildBase implements OnInit {
  patternEnv = /^[\w-$/\\=\"[\]{}@&:,'`\t. ?]+$/;
  envsData: Array<SharedEnvType>;
  envsText = '';
  inputValidator: Array<ValidatorFn>;
  configMapList: Array<SharedConfigMap>;
  configMapDetail: Map<number, SharedConfigMapDetail>;
  bindConfigMap: Map<number, boolean>;

  @Input() inputEnvsData: Array<SharedEnvType>;
  @Input() inputFixedKeyList: Array<string>;
  @Input() isProvideBindConfigMap = false;
  @Input() projectName = '';
  @Output() confirm: EventEmitter<Array<SharedEnvType>>;

  constructor(private messageService: MessageService,
              private sharedService: SharedService) {
    super();
    this.envsData = new Array<SharedEnvType>();
    this.inputFixedKeyList = new Array<string>();
    this.inputEnvsData = new Array<SharedEnvType>();
    this.confirm = new EventEmitter<Array<SharedEnvType>>();
    this.inputValidator = Array<ValidatorFn>();
    this.configMapList = Array<SharedConfigMap>();
    this.configMapDetail = new Map<number, SharedConfigMapDetail>();
    this.bindConfigMap = new Map<number, boolean>();
  }

  ngOnInit() {
    this.inputValidator.push(Validators.required);
    if (this.inputEnvsData && this.inputEnvsData.length > 0) {
      this.envsData = this.envsData.concat(this.inputEnvsData);
      this.envsData.forEach((value: SharedEnvType, index: number) => {
        const detail = new SharedConfigMapDetail();
        detail.dataList.push({key: value.envConfigMapKey, value: value.envConfigMapKey});
        this.configMapDetail.set(index, detail);
        this.bindConfigMap.set(index, value.envConfigMapKey !== '');
      });
    }
    this.sharedService.getConfigMapList(this.projectName, 0, 0).subscribe(
      (res: Array<SharedConfigMap>) => this.configMapList = res);
    this.modalOpened = true;
  }

  addNewEnv() {
    this.envsData.push(new SharedEnvType());
  }

  confirmEnvInfo() {
    if (this.verifyInputExValid() && this.verifyDropdownExValid()) {
      this.confirm.emit(this.envsData);
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
      const envTypes = this.envsText.split(';').map((str: string) => {
        const envStrPair = str.split('=');
        const env = new SharedEnvType();
        if (!this.patternEnv.test(envStrPair[0]) || !this.patternEnv.test(envStrPair[1])) {
          throw new Error();
        }
        env.envName = envStrPair[0];
        env.envValue = envStrPair[1];
        return env;
      });
      this.envsData = this.envsData.concat(envTypes);
    } catch (e) {
      this.messageService.showAlert('SERVICE.TXT_ALERT_MESSAGE', {alertType: 'warning', view: this.alertView});
      return;
    }
  }

  changeConfigMap(index: number, envInfo: SharedEnvType, configMap: SharedConfigMap) {
    envInfo.envConfigMapName = configMap.name;
    this.sharedService.getConfigMapDetail(configMap.name, this.projectName).subscribe(
      (res: SharedConfigMapDetail) => this.configMapDetail.set(index, res),
      (err: HttpErrorResponse) => this.messageService.showAlert(err.message, {alertType: 'danger', view: this.alertView})
    );
  }

  changeConfigMapKey(envInfo: SharedEnvType, data: { key: string, value: string }) {
    envInfo.envValue = data.value;
    envInfo.envConfigMapKey = data.key;
  }

  changeBindConfigMap(index: number, envInfo: SharedEnvType, event: Event) {
    const checked = (event.target as HTMLInputElement).checked;
    this.bindConfigMap.set(index, checked);
    if (!checked) {
      envInfo.envConfigMapName = '';
      envInfo.envConfigMapKey = '';
    }
  }

  isFixed(env: SharedEnvType): boolean {
    return this.inputFixedKeyList.find(value => env.envName === value) !== undefined;
  }
}
