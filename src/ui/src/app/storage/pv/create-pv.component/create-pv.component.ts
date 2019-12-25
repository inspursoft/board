import { Component, EventEmitter, OnInit, Type } from '@angular/core';
import { ValidationErrors } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import {
  NFSPersistentVolume,
  PersistentVolume,
  PvAccessMode,
  PvReclaimMode,
  RBDPersistentVolume
} from '../../../shared/shared.types';
import { StorageService } from '../../storage.service';
import { MessageService } from '../../../shared.service/message.service';

@Component({
  templateUrl: './create-pv.component.html',
  styleUrls: ['./create-pv.component.css']
})
export class CreatePvComponent extends CsModalChildBase implements OnInit {
  onAfterCommit: EventEmitter<PersistentVolume>;
  storageTypeList: Array<{name: string, value: number, classType: Type<PersistentVolume>}>;
  accessModeList: Array<PvAccessMode>;
  reclaimModeList: Array<PvReclaimMode>;
  newPersistentVolume: PersistentVolume;
  patternKeyring: RegExp = /^\/(\w+\/?)+$/;
  isEditPvMonitors = false;
  isCreateWip = false;
  namePattern: RegExp = /^[a-z0-9][(.a-z0-9?)]*$/;

  constructor(private storageService: StorageService,
              private messageService: MessageService) {
    super();
    this.storageTypeList = Array<{name: string, value: number, classType: Type<PersistentVolume>}>();
    this.newPersistentVolume = new NFSPersistentVolume();
    this.accessModeList = Array<PvAccessMode>();
    this.reclaimModeList = Array<PvReclaimMode>();
    this.onAfterCommit = new EventEmitter<PersistentVolume>();
  }

  ngOnInit() {
    this.accessModeList.push(PvAccessMode.ReadWriteOnce);
    this.accessModeList.push(PvAccessMode.ReadOnlyMany);
    this.accessModeList.push(PvAccessMode.ReadWriteMany);
    this.reclaimModeList.push(PvReclaimMode.Retain);
    this.reclaimModeList.push(PvReclaimMode.Recycle);
    this.reclaimModeList.push(PvReclaimMode.Delete);
    this.storageTypeList = [
      {name: 'NFS', value: 1, classType: NFSPersistentVolume},
      {name: 'Ceph rbd', value: 2, classType: RBDPersistentVolume}
    ];
  }

  changeAccessMode(mode: PvAccessMode) {
    this.newPersistentVolume.accessMode = mode;
  }

  createNewPv() {
    if (this.verifyInputExValid()) {
      this.isCreateWip = true;
      this.storageService.createNewPv(this.newPersistentVolume).subscribe(
        () => this.onAfterCommit.emit(this.newPersistentVolume),
        () => this.modalOpened = false,
        () => this.modalOpened = false
      );
    }
  }

  get nfsPersistentVolume(): NFSPersistentVolume {
    return this.newPersistentVolume as NFSPersistentVolume;
  }

  get rbdPersistentVolume(): RBDPersistentVolume {
    return this.newPersistentVolume as RBDPersistentVolume;
  }

  get checkPvNameFun() {
    return this.checkPvName.bind(this);
  }

  checkPvName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.storageService.checkPvNameExist(control.value)
      .pipe(map(() => null), catchError((err: HttpErrorResponse) => {
        if (err.status === 409) {
          this.messageService.cleanNotification();
          return of({pvNameExists: 'STORAGE.PV_NAME_EXIST'});
        } else if (err.status === 404) {
          this.messageService.cleanNotification();
        }
        return of(null);
      }));
  }

  changeSelectType(event: { name: string, value: number, classType: Type<PersistentVolume> }) {
    const newPv = new event.classType();
    newPv.name = this.newPersistentVolume.name;
    newPv.capacity = this.newPersistentVolume.capacity;
    newPv.reclaim = this.newPersistentVolume.reclaim;
    this.newPersistentVolume = newPv;
    this.newPersistentVolume.type = event.value;
    this.accessModeList = Array<PvAccessMode>();
    this.accessModeList.push(PvAccessMode.ReadWriteOnce);
    this.accessModeList.push(PvAccessMode.ReadOnlyMany);
    if (this.newPersistentVolume.type === 1) {
      this.accessModeList.push(PvAccessMode.ReadWriteMany);
    }
  }

  editMonitors() {
    this.isEditPvMonitors = true;
  }
}
