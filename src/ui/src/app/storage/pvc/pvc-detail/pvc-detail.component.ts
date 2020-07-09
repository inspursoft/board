import { Component } from '@angular/core';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { PersistentVolumeClaimDetail } from '../../sotrage.types';

@Component({
  templateUrl: './pvc-detail.component.html',
  styleUrls: ['./pvc-detail.component.css']
})
export class PvcDetailComponent extends CsModalChildBase {
  curDetail: PersistentVolumeClaimDetail;

  get eventDescription(): string {
    return this.curDetail.events.join(';');
  }
}

