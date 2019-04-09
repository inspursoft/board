import {
  Component, OnInit, OnDestroy,
  ComponentFactoryResolver, ViewChild, Type
} from '@angular/core';

import { ServiceStepBase } from './service-step';
import { K8sService } from './service.k8s';
import { StepService } from './service-step.service';
import { ServiceHostDirective } from "./service-host.directive";
import { Subscription } from "rxjs";

@Component({
  selector: 'service',
  templateUrl: './service.component.html'
})
export class ServiceComponent implements OnInit, OnDestroy {

  @ViewChild(ServiceHostDirective) serviceHost: ServiceHostDirective;
  serviceSteps: Array<Type<ServiceStepBase>>;
  currentStepIndex: number = 0;
  _subscription: Subscription;

  constructor(private k8sService: K8sService,
              private componentFactoryResolver: ComponentFactoryResolver) {
  }

  ngOnInit(): void {
    this._subscription = this.k8sService.step$.subscribe((stepInfo: {index: number, isBack: boolean}) => {
      this.currentStepIndex = stepInfo.index;
      this.loadComponent(stepInfo.isBack);
    });
    this.serviceSteps = StepService.getServiceSteps();
    this.k8sService.stepSource.next({index: 0, isBack: false});
  }

  ngOnDestroy() {
    if (this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  loadComponent(isBack: boolean) {
    let currentStep = this.serviceSteps[this.currentStepIndex];
    let componentFactory = this.componentFactoryResolver.resolveComponentFactory(currentStep);
    let viewContainerRef = this.serviceHost.viewContainerRef;
    viewContainerRef.clear();
    let nextComponent = viewContainerRef.createComponent(componentFactory);
    (nextComponent.instance as ServiceStepBase).isBack = isBack;
    (nextComponent.instance as ServiceStepBase).selfView = viewContainerRef;
    (nextComponent.instance as ServiceStepBase).factoryResolver = this.componentFactoryResolver;
  }
}
