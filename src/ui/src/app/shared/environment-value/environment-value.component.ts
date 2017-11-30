/**
 * Created by liyanq on 9/4/17.
 */

import {
  Component,
  EventEmitter,
  Input,
  Output,
  OnInit,
  ViewChildren,
  QueryList,
  AfterContentChecked
} from "@angular/core"
import { CsInputComponent } from "../cs-components-library/cs-input/cs-input.component";
import { ValidatorFn, Validators } from "@angular/forms";

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
export class EnvironmentValueComponent implements OnInit, AfterContentChecked {
  _isOpen: boolean = false;
  patternEnv:RegExp = /^[\w-$/\\=\"[\]{}@:,'`\t. ?]+$/;
  isCanConfirm: boolean = false;
  envAlertMessage: string;
  envsData: Array<EnvType>;
  isAlertOpen: boolean = false;
  afterCommitErr: string = "";
  inputValidator: Array<ValidatorFn>;
  inputValidatorMsg: Array<{validatorKey: string, validatorMessage: string}>;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @Input() inputEnvsData: Array<EnvType>;
  @Input() inputFixedKeyList: Array<string>;

  constructor() {
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

  ngAfterContentChecked() {
    if (this.inputComponents){
      let componentArr = this.inputComponents.toArray();
      for (let i = 0; i < componentArr.length; i++) {
        if (!componentArr[i].valid) {
          this.isCanConfirm = false;
          return
        }
      }
      this.isCanConfirm = true;
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
    this.inputComponents.forEach(value => {
      value.checkValueByHost();
    });
    this.onConfirm.emit(this.envsData);
    this.isOpen = false;
  }

  envMinusClick(index: number) {
    this.envsData.splice(index, 1);
  }

}