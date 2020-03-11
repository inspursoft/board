import { Component, HostBinding, OnInit, ViewChild } from "@angular/core";
import { KibanaService } from "../kibana.service";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  templateUrl: './kibana.component.html'
})
export class KibanaComponent implements OnInit {
  @ViewChild('frame') frame;
  errorMessage = '';
  kibanaUrl: string = '';

  constructor(private kibanaService: KibanaService) {

  }

  ngOnInit() {
    const url = '/kibana/';
    this.kibanaService.testKibana(url).subscribe(
      () => this.kibanaUrl = url,
      (err: HttpErrorResponse) => this.errorMessage = err.message)
  }

  @HostBinding('style.height') get height() {
    return '100%';
  }
}