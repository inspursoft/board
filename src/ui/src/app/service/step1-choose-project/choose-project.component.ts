import { Component, Input } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';

@Component({
  templateUrl: './choose-project.component.html'
})
export class ChooseProjectComponent implements ServiceStepComponent {
  @Input() data: any;

  constructor(private k8sService: K8sService){}

  forward() {
    this.k8sService.stepSource.next(2);
  }
}