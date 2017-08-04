import { Directive, ViewContainerRef } from '@angular/core';

@Directive({
  selector: '[service-host]'
})
export class ServiceHostDirective {
  constructor(public viewContainerRef: ViewContainerRef) {}
}