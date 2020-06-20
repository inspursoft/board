import {
  Component, OnInit, OnDestroy,
  ComponentFactoryResolver, ViewChild, Type, ViewContainerRef
} from '@angular/core';
import { Subscription } from 'rxjs';
import { K8sService } from './service.k8s';
import { StepService } from './service-step.service';
import { ServiceStepComponentBase } from './service-step';

@Component({
  selector: 'app-service',
  templateUrl: './service.component.html'
})
export class ServiceComponent implements OnInit, OnDestroy {
  @ViewChild('serviceHost', {read: ViewContainerRef}) serviceHostView: ViewContainerRef;
  serviceSteps: Array<Type<ServiceStepComponentBase>>;
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
    (nextComponent.instance as ServiceStepComponentBase).isBack = isBack;
    (nextComponent.instance as ServiceStepComponentBase).selfView = this.serviceHostView;
    (nextComponent.instance as ServiceStepComponentBase).factoryResolver = this.componentFactoryResolver;
  }
}
