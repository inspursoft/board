import {
  Component, OnInit, OnDestroy,
  ComponentFactoryResolver, ViewChild, Type, ViewContainerRef
} from '@angular/core';
import { Subscription } from 'rxjs';
import { ServiceStepBase } from './service-step';
import { K8sService } from './service.k8s';
import { StepService } from './service-step.service';

@Component({
  selector: 'app-service',
  templateUrl: './service.component.html'
})
export class ServiceComponent implements OnInit, OnDestroy {
  @ViewChild('serviceHost', {read: ViewContainerRef}) serviceHostView: ViewContainerRef;
  serviceSteps: Array<Type<ServiceStepBase>>;
  currentStepIndex = 0;
  subscription: Subscription;

  constructor(private k8sService: K8sService,
              private componentFactoryResolver: ComponentFactoryResolver) {
  }

  ngOnInit(): void {
    this.subscription = this.k8sService.step$.subscribe(
      (stepInfo: { index: number, isBack: boolean }) => {
        this.currentStepIndex = stepInfo.index;
        this.loadComponent(stepInfo.isBack);
      });
    this.serviceSteps = StepService.getServiceSteps();
    this.k8sService.stepSource.next({index: 0, isBack: false});
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  loadComponent(isBack: boolean) {
    this.serviceHostView.clear();
    const currentStep = this.serviceSteps[this.currentStepIndex];
    const componentFactory = this.componentFactoryResolver.resolveComponentFactory(currentStep);
    const nextComponent = this.serviceHostView.createComponent(componentFactory);
    (nextComponent.instance as ServiceStepBase).isBack = isBack;
    (nextComponent.instance as ServiceStepBase).selfView = this.serviceHostView;
    (nextComponent.instance as ServiceStepBase).factoryResolver = this.componentFactoryResolver;
  }
}
