import { Component, ComponentFactoryResolver, ViewContainerRef } from "@angular/core";
import { ClrDatagridStateInterface } from "@clr/angular";
import { TranslateService } from "@ngx-translate/core";
import { Message, PersistentVolumeClaim, RETURN_STATUS } from "../../../shared/shared.types";
import { MessageService } from "../../../shared/message-service/message.service";
import { StorageService } from "../../storage.service";
import { CsModalParentBase } from "../../../shared/cs-modal-base/cs-modal-parent-base";
import { CreatePvcComponent } from "../../../shared/create-pvc/create-pvc.component";

@Component({
  templateUrl: './pvc-list.component.html',
  styleUrls: ['./pvc-list.component.css']
})
export class PvcListComponent extends CsModalParentBase{
  isInLoadWip = false;
  pageIndex = 1;
  pageSize = 15;
  oldStateInfo: ClrDatagridStateInterface;
  pvcList: Array<PersistentVolumeClaim>;

  constructor(private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private storageService: StorageService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    super(resolver, view);
    this.pvcList = Array<PersistentVolumeClaim>();
  }

  refreshList(status: ClrDatagridStateInterface) {
    setTimeout(() => {
      this.isInLoadWip = true;
      this.oldStateInfo = status;
      this.storageService.getPvcList('', this.pageIndex, this.pageSize).subscribe(
        (res: Array<PersistentVolumeClaim>) => this.pvcList = res,
        () => this.isInLoadWip = false,
        () => this.isInLoadWip = false
      )
    })
  }

  createNewPvc(){
    this.createNewModal(CreatePvcComponent).onAfterCommit.subscribe(
      () => this.refreshList(this.oldStateInfo)
    );
  }

  deletePvc(pvcName: string, pvcId: number){
    this.translateService.get('STORAGE.PVC_DELETE_CONFIRM', [pvcName]).subscribe(res => {
      this.messageService.showDeleteDialog(res).subscribe((message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.storageService.deletePvc(pvcId).subscribe(
            () => this.messageService.showAlert(`STORAGE.PVC_DELETE_SUCCESS`),
            () => this.messageService.showAlert(`STORAGE.PVC_DELETE_FAILED`),
            () => this.refreshList(this.oldStateInfo))
        }
      })
    })
  }

  showPvcDetail(pvcId: number){

  }
}