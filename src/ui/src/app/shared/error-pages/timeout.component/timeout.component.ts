import { Component, OnDestroy } from "@angular/core"
import { Router } from "@angular/router";

@Component({
  templateUrl:"./timeout.component.html"
})
export class TimeoutComponent implements OnDestroy {
  intervalHandle: any;

  constructor(private router: Router) {
    this.intervalHandle = setInterval(() => this.router.navigate(["/sign-in"]), 5000);
  }

  ngOnDestroy() {
    clearInterval(this.intervalHandle)
  }
}