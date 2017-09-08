import { Component, Input, OnInit } from '@angular/core';
import { AppInitService } from '../../app.init.service';
import { K8sService } from '../service.k8s';
import { Service } from '../service';

@Component({
  templateUrl: './list-service.component.html'
})
export class ListServiceComponent implements OnInit {
  @Input() data: any;
  currentUser: {[key: string]: any};
  services: Service[];

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService) {
  }

  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser;
    this.retrieve();
  }


  createService(): void {
    this.k8sService.stepSource.next(1);
  }

  retrieve(): void {
    this.k8sService.getServices().then(res => this.services = res);
  }

  editService(s: Service) {

  }

  toggleServiceStatus(s: Service) {

  }

  confirmToDeleteService(s: Service) {

  }
}