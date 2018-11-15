import { Component } from "@angular/core";
import { NFSPersistentVolume, PersistentVolume, RBDPersistentVolume } from "../../../shared/shared.types";
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";

@Component({
  templateUrl: './pv-detail.component.html',
  styleUrls: ['./pv-detail.component.css']
})
export class PvDetailComponent extends CsModalChildBase {
  curPersistentVolume: PersistentVolume;

  get curNFSPersistentVolume(): NFSPersistentVolume {
    return this.curPersistentVolume as NFSPersistentVolume;
  }

  get curRBDPersistentVolume(): RBDPersistentVolume {
    return this.curPersistentVolume as RBDPersistentVolume;
  }
}