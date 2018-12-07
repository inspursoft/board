import { Component } from "@angular/core";
import { PersistentVolumeClaim } from "../../../shared/shared.types";
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";

@Component({
  templateUrl:"./pvc-detail.component.html",
  styleUrls:["./pvc-detail.component.css"]
})
export class PvcDetailComponent extends CsModalChildBase  {
  curPersistentVolumeClaim: PersistentVolumeClaim;

  get eventDescription(): string{
    return this.curPersistentVolumeClaim.events.join(';')
  }
}