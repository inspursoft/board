import { Component, EventEmitter, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ValidationErrors } from '@angular/forms';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { CsModalChildMessage } from '../cs-modal-base/cs-modal-child-base';
import { PersistentVolume, PersistentVolumeClaim, PvcAccessMode, SharedProject } from '../shared.types';
import { SharedService } from '../../shared.service/shared.service';
import { MessageService } from '../../shared.service/message.service';

@Component({
  templateUrl: './create-pvc.component.html',
  styleUrls: ['./create-pvc.component.css']
})
export class CreatePvcComponent extends CsModalChildMessage implements OnInit {
  onAfterCommit: EventEmitter<PersistentVolumeClaim>;
  projectsList: Array<SharedProject>;
  accessModeList: Array<PvcAccessMode>;
  pvList: Array<PersistentVolume>;
  newPersistentVolumeClaim: PersistentVolumeClaim;
  isCreateWip = false;
  namePattern: RegExp = /^[a-z0-9][(.a-z0-9?)]*$/;

  constructor(private sharedService: SharedService,
              public messageService: MessageService) {
    super(messageService);
    this.newPersistentVolumeClaim = new PersistentVolumeClaim();
    this.onAfterCommit = new EventEmitter<PersistentVolumeClaim>();
    this.projectsList = Array<SharedProject>();
    this.accessModeList = Array<PvcAccessMode>();
    this.pvList = Array<PersistentVolume>();
  }

  ngOnInit() {
    this.accessModeList.push(PvcAccessMode.ReadWriteOnce);
    this.accessModeList.push(PvcAccessMode.ReadWriteMany);
    this.accessModeList.push(PvcAccessMode.ReadOnlyMany);
    this.sharedService.getAllProjects().subscribe((res: Array<SharedProject>) => this.projectsList = res);
    this.sharedService.getPVList().subscribe((res: Array<PersistentVolume>) => {
      this.pvList = res;
      if (this.pvList.length > 0) {
        const pvNone = new PersistentVolume();
        pvNone.name = 'None';
        this.pvList.unshift(pvNone);
      }
    });
  }

  get checkPvcNameFun() {
    return this.checkPvcName.bind(this);
  }

  checkPvcName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.sharedService.checkPvcNameExist(this.newPersistentVolumeClaim.projectName, control.value)
      .pipe(
        map(() => null),
        catchError((err: HttpErrorResponse) => {
          this.messageService.cleanNotification();
          if (err.status === 409) {
            return of({pvNameExists: 'STORAGE.PVC_CREATE_NAME_EXIST'});
          }
          return of(null);
        })
      );
  }

  changeSelectProject(project: SharedProject) {
    this.newPersistentVolumeClaim.projectId = project.projectId;
    this.newPersistentVolumeClaim.projectName = project.projectName;
  }

  changeDesignatePv(pv: PersistentVolume) {
    this.newPersistentVolumeClaim.designatedPv = pv.name === 'None' ? '' : pv.name;
  }

  createNewPvc() {
    if (this.verifyDropdownExValid() && this.verifyInputExValid()) {
      this.isCreateWip = true;
      this.sharedService.createNewPvc(this.newPersistentVolumeClaim).subscribe(
        () => this.messageService.showAlert('STORAGE.PVC_CREATE_SUCCESS'),
        (error: HttpErrorResponse) => {
          this.messageService.showAlert(error.message, {alertType: 'warning', view: this.alertView});
          this.isCreateWip = false;
        },
        () => {
          this.onAfterCommit.emit(this.newPersistentVolumeClaim);
          this.modalOpened = false;
        }
      );
    }
  }
}
