import { Component, EventEmitter, OnInit } from "@angular/core";
import { Project } from "../../project/project";
import { CsModalChildBase } from "../cs-modal-base/cs-modal-child-base";
import { PersistentVolume, PersistentVolumeClaim, PvcAccessMode } from "../shared.types";
import { SharedService } from "../shared.service";
import { MessageService } from "../message-service/message.service";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  templateUrl: "./create-pvc.component.html",
  styleUrls: ["./create-pvc.component.css"]
})
export class CreatePvcComponent extends CsModalChildBase implements OnInit {
  onAfterCommit: EventEmitter<PersistentVolumeClaim>;
  projectsList: Array<Project>;
  accessModeList: Array<PvcAccessMode>;
  pvList: Array<PersistentVolume>;
  newPersistentVolumeClaim: PersistentVolumeClaim;
  isCreateWip = false;

  constructor(private sharedService: SharedService,
              private messageService: MessageService) {
    super();
    this.newPersistentVolumeClaim = new PersistentVolumeClaim();
    this.onAfterCommit = new EventEmitter<PersistentVolumeClaim>();
    this.projectsList = Array<Project>();
    this.accessModeList = Array<PvcAccessMode>();
    this.pvList = Array<PersistentVolume>();
  }

  ngOnInit() {
    this.accessModeList.push(PvcAccessMode.ReadWriteOnce);
    this.accessModeList.push(PvcAccessMode.ReadWriteMany);
    this.accessModeList.push(PvcAccessMode.ReadOnlyMany);
    this.sharedService.getAllProjects().subscribe((res: Array<Project>) => this.projectsList = res);
    this.sharedService.getAllPvList().subscribe((res: Array<PersistentVolume>) => this.pvList = res);
  }

  changeSelectProject(project: Project) {
    this.newPersistentVolumeClaim.projectId = project.project_id;
  }

  changeDesignatePv(pv: PersistentVolume) {
    this.newPersistentVolumeClaim.designatedPv = pv.name;
  }

  createNewPvc() {
    if (this.verifyInputValid() && this.verifyDropdownValid()) {
      this.isCreateWip = true;
      this.sharedService.createNewPvc(this.newPersistentVolumeClaim).subscribe(
        () => this.messageService.showAlert('STORAGE.PVC_CREATE_SUCCESS'),
        (error: HttpErrorResponse) => {
          this.messageService.showAlert(error.message, {alertType: "alert-warning", view: this.alertView});
          this.isCreateWip = false;
        },
        () => {
          this.onAfterCommit.emit(this.newPersistentVolumeClaim);
          this.modalOpened = false;
        }
      )
    }
  }
}