import { Component, AfterViewInit, OnDestroy, 
         ComponentFactoryResolver, 
         Input, Output, EventEmitter, ViewChild } from '@angular/core';

import { Subscription } from 'rxjs/Subscription';        
import { ServiceHostDirective } from './service-host.directive';
import { ServiceStep } from './service-step';
import { ServiceStepComponent } from './service-step.component';
import { K8sService } from './service.k8s';

@Component({
  selector: 'service-group',
  templateUrl: './service-group.component.html'
})
export class ServiceGroupComponent implements AfterViewInit, OnDestroy {

  @Input() serviceSteps: ServiceStep[];
  
  _subscription: Subscription;

  currentStepIndex: number = 0;

  @ViewChild(ServiceHostDirective) serviceHost: ServiceHostDirective;

  constructor(
    private componentFactoryResolver: ComponentFactoryResolver,
    private k8sService: K8sService
  ){}

  ngAfterViewInit(): void {
    this.loadComponent();
    this._subscription = this.k8sService.step$.subscribe((index: number)=>{
      this.currentStepIndex = index;
      this.loadComponent();
    });
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  loadComponent() {
    let currentStep = this.serviceSteps[this.currentStepIndex];
    let componentFactory = this.componentFactoryResolver.resolveComponentFactory(currentStep.component);
    let viewContainerRef = this.serviceHost.viewContainerRef;
    viewContainerRef.clear();
    let componentRef = viewContainerRef.createComponent(componentFactory);
    (<ServiceStepComponent>componentRef.instance).data = currentStep.data;
  }
}
