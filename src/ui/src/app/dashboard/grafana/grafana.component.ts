import { Component, HostBinding, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { DashboardService } from '../dashboard.service';

@Component({
  selector: 'app-grafana',
  templateUrl: './grafana.component.html'
})
export class GrafanaComponent implements OnInit {
  grafanaUrl = '';
  errorMessage = '';

  constructor(private dashboardService: DashboardService) {
  }

  ngOnInit() {
    const url = '/grafana/dashboard/db/kubernetes/';
    this.dashboardService.testGrafana(url).subscribe(
      () => this.grafanaUrl = url,
      (err: HttpErrorResponse) => this.errorMessage = err.message
    );
  }

  @HostBinding('style.height') get height() {
    return '100%';
  }
}
