import { Component, HostBinding, OnInit } from "@angular/core";
import { GrafanaService } from "../grafana.service";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  templateUrl: './grafana.component.html'
})
export class GrafanaComponent implements OnInit {
  grafanaUrl = '';
  errorMessage = '';

  constructor(private grafanaService: GrafanaService) {
  };

  ngOnInit() {
    const url = '/grafana/dashboard/db/kubernetes/';
    this.grafanaService.testGrafana(url).subscribe(
      () => this.grafanaUrl = url,
      (err: HttpErrorResponse) => this.errorMessage = err.message
    )
  }

  @HostBinding('style.height') get height() {
    return '100%';
  }
}