import { Component, Input } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';

@Component({
  templateUrl: './edit-container.component.html'
})
export class EditContainerComponent implements ServiceStepComponent {
  @Input() data: any;

  constructor(private k8sService: K8sService){}

  forward(): void {
    this.k8sService.stepSource.next(4);
  }
}