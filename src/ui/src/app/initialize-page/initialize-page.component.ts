import {
  Component, ElementRef,
  OnDestroy,
  OnInit,
  TemplateRef,
  ViewChild,
  ViewContainerRef
} from '@angular/core';
import { AppInitService } from '../shared.service/app-init.service';
import { Router } from '@angular/router';
import { RouteSignIn } from '../shared/shared.const';
import { MessageService } from '../shared.service/message.service';
import { HttpErrorResponse } from '@angular/common/http';
import { GlobalAlertType } from '../shared/shared.types';
import { interval, Subscription } from 'rxjs';

@Component({
  selector: 'app-initialize-page',
  templateUrl: './initialize-page.component.html',
  styleUrls: ['./initialize-page.component.css']
})
export class InitializePageComponent implements OnInit, OnDestroy {
  @ViewChild('logTemplate') logTemplate: TemplateRef<any>;
  @ViewChild('divElement') divElement: ElementRef;
  @ViewChild('logContainer', {read: ViewContainerRef}) logContainer: ViewContainerRef;
  messageMap: Map<string, string>;
  subscription: Subscription;
  tailPoints = '.';
  totalSteps = 1;
  curStep = 1;
  curStepMessage = '';

  constructor(private appInitService: AppInitService,
              private route: Router,
              private messageService: MessageService) {
    this.messageMap = new Map<string, string>();
  }

  ngOnInit() {
    this.reloadData();
    this.messageMap.set('UPDATE_ADMIN_PASSWORD', 'InitializeInfo.UpdateAdminPassword');
    this.messageMap.set('INIT_PROJECT_REPO', 'InitializeInfo.InitProjectRepo');
    this.messageMap.set('PREPARE_KVM_HOST', 'InitializeInfo.PrepareKvmHost');
    this.messageMap.set('INIT_KUBERNETES_INFO', 'InitializeInfo.InitKubernetesInfo');
    this.messageMap.set('SYNC_UP_K8S', 'InitializeInfo.SyncUpK8s');
    this.messageMap.set('READY', 'InitializeInfo.Ready');
    this.subscription = interval(2000).subscribe(
      (res: number) => {
        if (res % 3 === 0) {
          this.tailPoints = '..';
        } else if (res % 3 === 1) {
          this.tailPoints = '...';
        } else {
          this.tailPoints = '.';
        }
        this.reloadData();
      }
    );
  }

  ngOnDestroy(): void {
    this.subscription.unsubscribe();
  }

  reloadData() {
    this.appInitService.getSystemInfo().subscribe(
      () => this.route.navigate([RouteSignIn]).then(),
      (err: HttpErrorResponse) => {
        if (err.status === 406) {
          this.messageService.cleanNotification();
          const res: { status: number, message: string } = err.error;
          if (res.message.includes('_')) {
            const messageArr = res.message.split('_');
            this.totalSteps = Number(messageArr[0]);
            this.curStep = Number(messageArr[1]);
            const messageKey = res.message.substr(res.message.indexOf(messageArr[1]) + 2);
            this.curStepMessage = this.messageMap.get(messageKey);
          } else {
            this.curStepMessage = this.messageMap.get(res.message);
          }
        } else {
          this.messageService.showGlobalMessage(err.message, {
            globalAlertType: GlobalAlertType.gatShowDetail,
            errorObject: err
          });
        }
      }
    );
  }
}
