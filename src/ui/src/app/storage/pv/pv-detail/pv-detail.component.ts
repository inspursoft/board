import { Component } from '@angular/core';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { NFSPersistentVolume, PersistentVolume, RBDPersistentVolume } from '../../sotrage.types';

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
