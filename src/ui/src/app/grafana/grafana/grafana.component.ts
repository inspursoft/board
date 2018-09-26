import { Component, HostBinding } from "@angular/core";

@Component({
  templateUrl: './grafana.component.html'
})
export class GrafanaComponent {
  grafanaUrl: string = '/grafana/dashboard/db/kubernetes/';

  @HostBinding('style.height') get height() {
    return '100%';
  }
}