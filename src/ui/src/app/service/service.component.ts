import {
  Component, OnInit, Output, EventEmitter, OnDestroy,
  ComponentFactoryResolver, ViewChild
} from '@angular/core';

import { ServiceStep } from './service-step';
import { K8sService } from './service.k8s';
import { StepService } from './service-step.service';
import { ServiceStepComponent } from "./service-step.component";
import { ServiceHostDirective } from "./service-host.directive";
import { Subscription } from 'rxjs/Subscription';

@Component({
  selector: 'service',
  templateUrl: './service.component.html'
})
export class ServiceComponent implements OnInit, OnDestroy {

  serviceSteps: ServiceStep[];

  @ViewChild(ServiceHostDirective) serviceHost: ServiceHostDirective;

  currentStepIndex: number = 0;
  _subscription: Subscription;

  constructor(
    private k8sService: K8sService,
    private componentFactoryResolver: ComponentFactoryResolver,
    private stepService: StepService
  ){}

  ngOnInit(): void {
    this._subscription = this.k8sService.step$.subscribe((index: number)=>{
      this.currentStepIndex = index;
      this.loadComponent();
    });
    this.serviceSteps = this.stepService.getServiceSteps();
    this.k8sService.stepSource.next(0);
  }

  ngOnDestroy(){
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  loadComponent(){
    let currentStep = this.serviceSteps[this.currentStepIndex];
    let componentFactory = this.componentFactoryResolver.resolveComponentFactory(currentStep.component);
    let viewContainerRef = this.serviceHost.viewContainerRef;
    viewContainerRef.clear();
    let componentRef = viewContainerRef.createComponent(componentFactory);
    (<ServiceStepComponent>componentRef.instance).data = currentStep.data;
  }
}
