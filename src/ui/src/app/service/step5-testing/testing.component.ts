import { Component, Injector } from '@angular/core';
import { ServiceStepBase } from "../service-step";
@Component({
  templateUrl: './testing.component.html',
  styleUrls: ["./testing.component.css"]
})
export class TestingComponent extends ServiceStepBase {
  inTesting: boolean = false;

  constructor(protected injector: Injector) {
    super(injector);
  }

  forward(): void {
    this.k8sService.stepSource.next({index: 6, isBack: false});
  }
}