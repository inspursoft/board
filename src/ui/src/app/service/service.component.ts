import { Component, OnInit, Output, EventEmitter } from '@angular/core';

import { ServiceStep } from './service-step';
import { K8sService } from './service.k8s';
import { StepService } from './service-step.service';

@Component({
  selector: 'service',
  templateUrl: './service.component.html'
})
export class ServiceComponent implements OnInit {

  serviceSteps: ServiceStep[];

  hideNavigationButton: boolean = true;

  constructor(
    private k8sService: K8sService,
    private stepService: StepService){}

  ngOnInit(): void {
    this.serviceSteps = this.stepService.getServiceSteps();
    this.k8sService.stepSource.next(0);
  }
}
