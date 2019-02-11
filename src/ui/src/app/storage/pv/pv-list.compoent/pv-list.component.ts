import { Component, ComponentFactoryResolver, ViewContainerRef } from "@angular/core";
import { ClrDatagridStateInterface } from "@clr/angular";
import { TranslateService } from "@ngx-translate/core";
import { Message, PersistentVolume, RETURN_STATUS } from "../../../shared/shared.types";
import { CsModalParentBase } from "../../../shared/cs-modal-base/cs-modal-parent-base";
import { CreatePvComponent } from "../create-pv.component/create-pv.component";
import { StorageService } from "../../storage.service";
import { PvDetailComponent } from "../pv-detail.compoent/pv-detail.component";
import { MessageService } from "../../../shared/message-service/message.service";

@Component({
  templateUrl: './pv-list.component.html',
  styleUrls: ['./pv-list.component.css']
})
export class PvListComponent extends CsModalParentBase {
  isInLoadWip = false;
  pageIndex = 1;
  pageSize = 15;
  oldStateInfo: ClrDatagridStateInterface;
  pvList: Array<PersistentVolume>;

  constructor(private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private storageService: StorageService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    super(resolver, view);
    this.pvList = Array<PersistentVolume>();
  }

  refreshList(status: ClrDatagridStateInterface) {
    setTimeout(() => {
      this.isInLoadWip = true;
      this.oldStateInfo = status;
      this.storageService.getPvList('', this.pageIndex, this.pageSize).subscribe(
        (res: Array<PersistentVolume>) => this.pvList = res,
        () => this.isInLoadWip = false,
        () => this.isInLoadWip = false
      )
    })
  }

  createNewPv() {
    this.createNewModal(CreatePvComponent).onAfterCommit.subscribe(
      () => this.refreshList(this.oldStateInfo)
    );
  }

  showPvDetail(pvId: number) {
    this.storageService.getPvDetailInfo(pvId).subscribe((res: PersistentVolume) => {
      let instance = this.createNewModal(PvDetailComponent);
      instance.curPersistentVolume = res;
    });
  }

  deletePv(pvName: string, pvId: number) {
    this.translateService.get('STORAGE.PV_DELETE_CONFIRM', [pvName]).subscribe(res => {
      this.messageService.showDeleteDialog(res).subscribe((message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.storageService.deletePv(pvId).subscribe(
            () => this.messageService.showAlert(`STORAGE.PV_DELETE_SUCCESS`),
            () => this.messageService.showAlert(`STORAGE.PV_DELETE_FAILED`),
            () => this.refreshList(this.oldStateInfo))
        }
      })
    })
  }
}