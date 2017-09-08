/**
 * Created by liyanq on 9/1/17.
 */

import { Component, Input, Output, EventEmitter } from "@angular/core"

export enum ContainerFieldStatus{cfsView, cfsEdit}
export enum FieldType{ftWithInput, ftWithNoInput}
export type ContainerFieldType = ContainerField<string> | ContainerField<number>;
export class ContainerField<T> {
  constructor(public status: ContainerFieldStatus,
              public defaultValue: T,
              public value: T) {
  }

  get typeName(): string {
    return typeof this.value;
  }
}

@Component({
  selector: "container-field",
  templateUrl: "./container-field.component.html",
  styleUrls: ["./container-field.component.css"]
})
export class ContainerFieldComponent {
  @Input("Label") labelText: string = "";
  @Input("Field") curField: ContainerFieldType;
  @Input("Type") curFieldType: FieldType = FieldType.ftWithInput;
  @Output("OnEdit") onEditEvent: EventEmitter<any> = new EventEmitter<any>();

  onEditIconClick() {
    if (this.curFieldType == FieldType.ftWithInput){
      this.curField.status = 1;
    } else{
      this.onEditEvent.emit();
    }
  }
}
