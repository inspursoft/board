import { Component, HostBinding, Input, OnDestroy, OnInit } from '@angular/core';
import { HttpProgressEvent } from '@angular/common/http';
import { interval, Subscription } from 'rxjs';

@Component({
  selector: 'cs-progress',
  styleUrls: ['./cs-progress.component.css'],
  templateUrl: './cs-progress.component.html'
})
export class CsProgressComponent implements OnInit, OnDestroy {
  @Input() progressData: HttpProgressEvent;
  @HostBinding('style.width') width = '100%';
  @HostBinding('style.height') height = '100%';
  private previousValue = 0;
  private subscription: Subscription;
  progressValue = 0;
  speed = 0;

  ngOnInit() {
    this.subscription = interval(500).subscribe(() => {
      if (this.progressData) {
        if (this.progressData.loaded < this.progressData.total) {
          if (this.previousValue > 0) {
            this.speed = this.progressData.loaded - this.previousValue;
          }
          this.previousValue = this.progressData.loaded;
          this.progressValue = Math.round(this.progressData.loaded / this.progressData.total * 1000) / 10;
        } else if (this.progressData.loaded === this.progressData.total) {
          this.progressValue = 100;
        }
      }
    });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  get isShowSelf(): boolean {
    const selfVisible = this.progressData && this.progressData.loaded < this.progressData.total;
    if (!selfVisible) {
      this.progressValue = 0;
    }
    return selfVisible;
  }
}
