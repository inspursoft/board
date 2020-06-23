import { Component, Injector } from '@angular/core';
import { ServiceStepComponentBase } from '../service-step';


@Component({
  templateUrl: './testing.component.html',
  styleUrls: ['./testing.component.css']
})
export class TestingComponent extends ServiceStepComponentBase {
  inTesting = false;

  constructor(protected injector: Injector) {
    super(injector);
  }

  forward(): void {
    this.k8sService.stepSource.next({index: 5, isBack: false});
  }
}

