import { Component, Input } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
import { K8sService } from "../service.k8s";
@Component({
  templateUrl: './testing.component.html',
  styleUrls: ["./testing.component.css"]
})
export class TestingComponent implements ServiceStepComponent {
  @Input() data: any;
  inTesting:boolean = false;

  constructor(private k8sService: K8sService) {
  }


  forward(): void {
    this.k8sService.stepSource.next(6);
  }
}