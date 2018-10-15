import { Component, OnDestroy } from "@angular/core"
import { Router } from "@angular/router";

@Component({
  templateUrl: "./bad-gateway.component.html"
})
export class BadGatewayComponent implements OnDestroy {
  intervalHandle: any;

  constructor(private route: Router) {
    this.intervalHandle = setInterval(() => this.route.navigate(["/sign-in"]), 5000);
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandle);
  }
}