/**
 * Created by liyanq on 9/17/17.
 */

import { Component } from "@angular/core"
import { K8sService } from "../service.k8s";

@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent {

  constructor(private k8sService: K8sService) {
  }

  forward(): void {
    this.k8sService.stepSource.next(6);
  }
}