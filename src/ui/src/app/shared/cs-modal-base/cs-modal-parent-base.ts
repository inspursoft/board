import { ComponentFactoryResolver, Type, ViewContainerRef } from '@angular/core';
import { CsModalChildBase } from './cs-modal-child-base';
import { CsComponentBase } from '../cs-components-library/cs-component-base';

export class CsModalParentBase extends CsComponentBase {
  constructor(protected factoryResolver: ComponentFactoryResolver,
              protected selfView: ViewContainerRef) {
    super();
  }

  createNewModal<T extends CsModalChildBase>(newComponent: Type<T>): T {
    let factory = this.factoryResolver.resolveComponentFactory(newComponent);
    let componentRef = this.selfView.createComponent(factory);
    componentRef.instance.openModal().subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    return componentRef.instance;
  }
}