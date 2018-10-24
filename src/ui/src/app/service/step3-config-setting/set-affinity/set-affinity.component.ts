import { Component } from "@angular/core";
import { Observable } from "rxjs/Rx";
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";
import { AffinityCardData, UIServiceStep4 } from "../../service-step.component";
import { DragStatus } from "../../../shared/shared.types";
import { HttpErrorResponse } from "@angular/common/http";
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared/message-service/message.service";
import { Service } from "../../service";
import { SERVICE_STATUS } from "../../../shared/shared.const";

@Component({
  templateUrl: './set-affinity.component.html',
  styleUrls: ['./set-affinity.component.css']
})
export class SetAffinityComponent extends CsModalChildBase {
  isActionWip = false;
  affinitySourceDataList: Array<AffinityCardData>;
  uiData: UIServiceStep4;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    super();
    this.affinitySourceDataList = Array<AffinityCardData>();
  }

  addNewAffinity() {
    this.uiData.affinityList.push({flag: true, services: Array<AffinityCardData>()})
  }

  deleteAffinity(index: number) {
    if (!this.isActionWip) {
      this.uiData.affinityList[index].services.forEach(value => {
        value.status = DragStatus.dsReady;
        this.affinitySourceDataList.push(value);
      });
      this.uiData.affinityList.splice(index, 1);
    }
  }

  initAffinity() {
    this.isActionWip = true;
    this.k8sService.getCollaborativeService(this.uiData.serviceName, this.uiData.projectName).subscribe((res: Array<Service>) => {
      this.isActionWip = false;
      res.forEach(value => {
        let serviceInUsed = this.uiData.affinityList.find(
          value1 => value1.services.find(
            card => card.serviceName === value.service_name) !== undefined);
        if (!serviceInUsed) {
          let card = new AffinityCardData();
          card.serviceName = value.service_name;
          card.serviceStatus = value.service_status;
          card.status = DragStatus.dsReady;
          this.affinitySourceDataList.push(card);
        }
      });
      this.uiData.affinityList.forEach(card => {
        card.services.forEach((affinity: AffinityCardData) => {
          let source = res.find(source => source.service_name === affinity.key);
          affinity.serviceStatus = source ? source.service_status : SERVICE_STATUS.DELETED
        })
      })
    }, (err: HttpErrorResponse) => {
      if (err.status == 404) {
        this.messageService.cleanNotification();
      }
    });
  }

  openSetModal(uiData: UIServiceStep4): Observable<any> {
    this.uiData = uiData;
    if (this.uiData.affinityList.length == 0) {
      this.addNewAffinity();
    }
    this.initAffinity();
    return super.openModal();
  }
}