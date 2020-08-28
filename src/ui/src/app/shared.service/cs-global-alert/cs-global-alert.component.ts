import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { Observable, Subject } from 'rxjs';
import { Router } from '@angular/router';
import { GlobalAlertMessage } from '../../shared/shared.types';
import { RouteSignIn } from '../../shared/shared.const';
import { AppInitService } from '../app-init.service';

@Component({
  templateUrl: './cs-global-alert.component.html',
  styleUrls: ['./cs-global-alert.component.css']
})
export class CsGlobalAlertComponent implements OnInit {
  isOpenValue = false;
  curMessage: GlobalAlertMessage;
  onCloseEvent: Subject<any>;
  detailModalOpen = false;
  curErrorDetailMsg = '';


  constructor(private route: Router,
              private changeRef: ChangeDetectorRef,
              private appInitService: AppInitService) {
    this.onCloseEvent = new Subject<any>();
  }

  ngOnInit(): void {
    this.getErrorDetailMsg();
  }

  get isOpen(): boolean {
    return this.isOpenValue;
  }

  set isOpen(value: boolean) {
    this.isOpenValue = value;
    if (!value) {
      this.onCloseEvent.next();
    }
  }

  getErrorDetailMsg() {
    if (this.curMessage.errorObject && this.curMessage.errorObject instanceof HttpErrorResponse) {
      const err = (this.curMessage.errorObject as HttpErrorResponse).error;
      if (typeof err === 'object') {
        if (err instanceof Blob) {
          const reader = new FileReader();
          reader.addEventListener('loadend', () => {
            this.curErrorDetailMsg = reader.result as string;
            this.changeRef.detectChanges();
          });
          reader.readAsText(err);
        } else {
          this.curErrorDetailMsg = err ? err.message : (this.curMessage.errorObject as HttpErrorResponse).message;
        }
      } else {
        this.curErrorDetailMsg = err;
      }
    } else if (this.curMessage.errorObject) {
      this.curErrorDetailMsg = (this.curMessage.errorObject as Error).message;
    }
  }

  public openAlert(message: GlobalAlertMessage): Observable<any> {
    this.curMessage = message;
    this.isOpen = true;
    return this.onCloseEvent.asObservable();
  }

  loginClick() {
    const authMode = this.appInitService.systemInfo.authMode;
    const redirectionURL = this.appInitService.systemInfo.redirectionUrl;
    if (authMode === 'indata_auth') {
      window.location.href = redirectionURL;
      this.isOpen = false;
      return;
    }
    this.isOpen = false;
    this.route.navigate([RouteSignIn]).then();
  }
}
