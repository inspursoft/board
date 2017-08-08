import { Component, Input } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';

@Component({
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent implements ServiceStepComponent {
  @Input() data: any;

  constructor(private k8sService: K8sService){}

  forward(): void {
    this.k8sService.stepSource.next(5);
  }
}