/**
 * Created by liyanq on 9/4/17.
 */

import { Component, EventEmitter, Input, OnInit, Output, QueryList, ViewChildren } from "@angular/core"
import { CsInputComponent } from "../cs-components-library/cs-input/cs-input.component";
import { ValidatorFn, Validators } from "@angular/forms";
import { CsComponentBase } from "../cs-components-library/cs-component-base";

export class EnvType {
  constructor(public envName: string,
              public envValue: string) {
  }
}
@Component({
  selector: "environment-value",
  templateUrl: "./environment-value.component.html",
  styleUrls: ["./environment-value.component.css"]
})
export class EnvironmentValueComponent extends CsComponentBase implements OnInit {
  _isOpen: boolean = false;
  patternEnv:RegExp = /^[\w-$/\\=\"[\]{}@&:,'`\t. ?]+$/;
  envAlertMessage: string;
  envsData: Array<EnvType>;
  envsText: string = "";
  isAlertOpen: boolean = false;
  inputValidator: Array<ValidatorFn>;
  inputValidatorMsg: Array<{validatorKey: string, validatorMessage: string}>;
  @ViewChildren(CsInputComponent) inputComponentList: QueryList<CsInputComponent>;
  @Input() inputEnvsData: Array<EnvType>;
  @Input() inputFixedKeyList: Array<string>;

  constructor() {
    super();
    this.envsData = Array<EnvType>();
    this.inputValidator = Array<ValidatorFn>();
    this.inputValidatorMsg = Array<{validatorKey: string, validatorMessage: string}>();
  }

  ngOnInit() {
    this.inputValidator.push(Validators.required);
    this.inputValidatorMsg.push({validatorKey: "required", validatorMessage: "SERVICE.ENV_REQUIRED"});
    if (this.inputEnvsData && this.inputEnvsData.length > 0) {
      this.envsData = this.envsData.concat(this.inputEnvsData);
    }
  }

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  @Output() onConfirm: EventEmitter<Array<EnvType>> = new EventEmitter<Array<EnvType>>();

  addNewEnv() {
    this.envsData.push(new EnvType("", ""));
  }

  confirmEnvInfo() {
    if (this.verifyInputValid()) {
      this.onConfirm.emit(this.envsData);
      this.isOpen = false;
    }
  }

  envMinusClick(index: number) {
    this.envsData.splice(index, 1);
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
      this.isAlertOpen = true;
      this.envAlertMessage = "SERVICE.TXT_ALERT_MESSAGE";
      return;
    }
    this.envsData = this.envsData.concat(envTypes);
    this.isAlertOpen = false;
  }
}