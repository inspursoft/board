import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { Affinity, AffinityCardData, ServiceStep3Data } from '../../service-step.component';
import { DragStatus } from '../../../shared/shared.types';
import { K8sService } from '../../service.k8s';
import { MessageService } from '../../../shared.service/message.service';
import { SERVICE_STATUS } from '../../../shared/shared.const';
import { Service } from '../../service.types';

@Component({
  templateUrl: './set-affinity.component.html',
  styleUrls: ['./set-affinity.component.css']
})
export class SetAffinityComponent extends CsModalChildBase {
  isActionWip = false;
  affinitySourceDataList: Array<AffinityCardData>;
  uiData: ServiceStep3Data;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    super();
    this.affinitySourceDataList = Array<AffinityCardData>();
  }

  addNewAffinity() {
    this.uiData.affinityList.push(new Affinity());
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
        const serviceInUsed = this.uiData.affinityList.find(
          value1 => value1.services.find(
            card => card.serviceName === value.serviceName) !== undefined);
        if (!serviceInUsed) {
          const card = new AffinityCardData();
          card.serviceName = value.serviceName;
          card.serviceStatus = value.serviceStatus;
          card.status = DragStatus.dsReady;
          this.affinitySourceDataList.push(card);
        }
      });
      this.uiData.affinityList.forEach(card => {
        card.services.forEach((affinity: AffinityCardData) => {
          const sourceService = res.find(source => source.serviceName === affinity.key);
          affinity.serviceStatus = sourceService ? sourceService.serviceStatus : SERVICE_STATUS.DELETED;
        });
      });
    }, (err: HttpErrorResponse) => {
      if (err.status === 404) {
        this.messageService.cleanNotification();
      }
    });
  }

  openSetModal(uiData: ServiceStep3Data): Observable<any> {
    this.uiData = uiData;
    this.closeNotification.subscribe(() => {
      const list = this.uiData.affinityList;
      for (let i = list.length - 1; i >= 0; i--) {
        if (list[i].services.length === 0) {
          list.splice(i, 1);
        }
      }
    });
    if (this.uiData.affinityList.length === 0) {
      this.addNewAffinity();
    }
    this.initAffinity();
    return super.openModal();
  }
}
