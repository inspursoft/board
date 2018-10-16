import { ChangeDetectorRef, Component, Injector, OnInit, ViewChild } from '@angular/core';
import {
  ConfigCardData,
  Container,
  PHASE_CONFIG_CONTAINERS,
  PHASE_EXTERNAL_SERVICE,
  ServiceStepPhase,
  UIServiceStep3,
  UIServiceStep4
} from '../service-step.component';
import { ServiceStepBase } from "../service-step";
import { ValidationErrors } from "@angular/forms/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs/Observable";
import { DragStatus } from "../../shared/shared.types";
import { SetExternalComponent } from "./set-external-port/set-external.component";
import { ConfigCardListComponent } from "./config-card-list/config-card-list.component";

@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent extends ServiceStepBase implements OnInit {
  @ViewChild('external') externalList: ConfigCardListComponent;
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  containerSourceDataList: Array<ConfigCardData>;
  affineSourceDataList: Array<ConfigCardData>;
  nodeSelectorCardList: Array<ConfigCardData>;
  uiPreData: UIServiceStep3;
  noPortForExtent = false;
  tabBaseActive = true;
  tabAdvanceActive = false;
  isActionWip: boolean = false;

  constructor(protected injector: Injector,
              private changeDetectorRef: ChangeDetectorRef) {
    super(injector);
    this.containerSourceDataList = Array<ConfigCardData>();
    this.affineSourceDataList = Array<ConfigCardData>();
    this.nodeSelectorCardList = Array<ConfigCardData>();
    this.uiPreData = new UIServiceStep3();
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS).subscribe(res => {
      this.uiPreData = res as UIServiceStep3;
      this.uiPreData.containerList.forEach((container: Container) => {
        container.container_port.forEach(port => {
          let card = new ConfigCardData();
          card.cardName = container.name;
          card.containerPort = port;
          card.status = DragStatus.dsReady;
          this.containerSourceDataList.push(card);
        });
      });
      this.noPortForExtent = this.uiPreData.containerList.every(value => !value.isHavePort())
    });
    this.k8sService.getServiceConfig(this.stepPhase).subscribe(res => {
      this.uiBaseData = res;
      this.setServiceName(this.uiData.serviceName);
      this.changeDetectorRef.detectChanges();
    });
    this.k8sService.getNodeSelectors().subscribe((res:Array<string>)=>{
      res.forEach(value => {
        let card = new ConfigCardData();
        card.cardName = value;
        card.status = DragStatus.dsReady;
        this.nodeSelectorCardList.push(card);
      });
    });
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_EXTERNAL_SERVICE
  }

  get uiData(): UIServiceStep4 {
    return this.uiBaseData as UIServiceStep4;
  }

  get checkServiceNameFun() {
    return this.checkServiceName.bind(this);
  }

  checkServiceName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.checkServiceExist(this.uiData.projectName, control.value)
      .map(() => null)
      .catch((err:HttpErrorResponse) => {
        if (err.status == 409) {
          this.messageService.cleanNotification();
          return Observable.of({serviceExist: "SERVICE.STEP_3_SERVICE_NAME_EXIST"});
        } else if (err.status == 404) {
          this.messageService.cleanNotification();
        }
        return Observable.of(null);
      });
  }

  setServiceName(serviceName: string): void {
    this.uiData.serviceName = serviceName;
    this.k8sService.getCollaborativeService(serviceName, this.uiData.projectName).subscribe((res: Array<string>) => {
        res.forEach(value => {
          let card = new ConfigCardData();
          card.cardName = value;
          card.status = DragStatus.dsReady;
          this.affineSourceDataList.push(card);
        });
      },
      (err: HttpErrorResponse) => {
        if (err.status == 404) {
          this.messageService.cleanNotification();
        }
      });
  }

  addNewAffinity() {
    this.uiData.affinityList.push({flag: 0, services: Array<ConfigCardData>()})
  }

  deleteAffinity(index: number) {
    this.uiData.affinityList[index].services.forEach(value => {
      value.status = DragStatus.dsReady;
      this.affineSourceDataList.push(value);
    });
    this.uiData.affinityList.splice(index, 1);
  }

  setExternalPort(data: ConfigCardData): void {
    let factory = this.factoryResolver.resolveComponentFactory(SetExternalComponent);
    let componentRef = this.selfView.createComponent(factory);
    componentRef.instance.openSetModal(data).subscribe(() => {
      if (!componentRef.instance.alreadySet) {
        this.externalList.removeContainerCard(data);
      }
      this.selfView.remove(this.selfView.indexOf(componentRef.hostView))
    });
  }

  forward(): void {
    let funExecute = () => {
      if (this.uiData.externalServiceList.length == 0) {
        this.tabBaseActive = true;
        this.messageService.showAlert(`SERVICE.STEP_3_WARNING_MESSAGE`, {alertType: "alert-warning"});
      } else {
        this.isActionWip = true;
        this.k8sService.setServiceConfig(this.uiData.uiToServer()).subscribe(
          () => this.k8sService.stepSource.next({index: 5, isBack: false})
        );
      }
    };
    if (this.tabAdvanceActive) {
      if (this.uiData.serviceName === '') {
        this.tabBaseActive = true;
        this.messageService.showAlert(`SERVICE.STEP_3_SERVICE_NAME_EMPTY`, {alertType: "alert-warning"});
      } else {
        funExecute();
      }
    } else if (this.verifyInputValid()) {
      funExecute();
    }
  }

  backUpStep(): void {
    this.k8sService.stepSource.next({index: 2, isBack: true});
  }
}