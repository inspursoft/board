import { Component, HostBinding } from "@angular/core";

@Component({
  templateUrl: './kibana.component.html'
})
export class KibanaComponent {
  kibanaUrl: string = '/kibana/';

  @HostBinding('style.height') get height() {
    return '100%';
  }
}