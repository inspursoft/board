import { Component, Input } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
@Component({
  templateUrl: './deploy-testing.component.html'
})
export class DeployTestingComponent implements ServiceStepComponent {
  @Input() data: any;
}