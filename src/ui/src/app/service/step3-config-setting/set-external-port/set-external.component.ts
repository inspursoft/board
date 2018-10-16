import { Component } from "@angular/core";
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";
import { ConfigCardData } from "../../service-step.component";
import { Observable } from "rxjs/Rx";

@Component({
  templateUrl: './set-external.component.html',
  styleUrls: ['./set-external.component.css']
})
export class SetExternalComponent extends CsModalChildBase {
  curData: ConfigCardData;
  alreadySet = false;

  openSetModal(data: ConfigCardData): Observable<any> {
    this.curData = data;
    this.alreadySet = false;
    return super.openModal();
  }

  setExternalPort() {
    if (this.verifyInputValid()) {
      this.alreadySet = true;
      this.modalOpened = false;
    }
  }
}