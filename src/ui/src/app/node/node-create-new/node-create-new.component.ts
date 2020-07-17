import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { EdgeNode, NodeStatus } from '../node.types';
import { interval, Observable, of, Subject, Subscription, TimeoutError } from 'rxjs';
import { ValidationErrors } from '@angular/forms';
import { catchError, map } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';
import { NodeService } from '../node.service';
import { MessageService } from '../../shared.service/message.service';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-node-create-new',
  templateUrl: './node-create-new.component.html',
  styleUrls: ['./node-create-new.component.css']
})
export class NodeCreateNewComponent extends CsModalChildBase implements OnInit, OnDestroy {
  nodeList: Array<NodeStatus>;
  patternNodeName: RegExp = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*[a-zA-Z0-9]*$/;
  patternNodeIp: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  edgeNode: EdgeNode;
  cpuTypes: Array<{ description: string }>;
  masters: Array<string>;
  registryMode: Array<{ key: string, value: string }>;
  isActionWip = false;
  subRefreshInterval: Subscription;
  checkTimes = 0;

  constructor(private nodeService: NodeService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    super();
    this.edgeNode = new EdgeNode();
    this.cpuTypes = new Array<{ description: string }>();
    this.masters = new Array<string>();
    this.nodeList = new Array<NodeStatus>();
    this.registryMode = new Array<{ key: string, value: string }>();
  }

  ngOnInit() {
    this.cpuTypes.push({description: 'auto-detect'});
    this.cpuTypes.push({description: 'X86'});
    this.cpuTypes.push({description: 'ARM64'});
    this.cpuTypes.push({description: 'ARM32'});
    this.translateService.get(['NodeCreateNew.Auto', 'NodeCreateNew.Manual']).subscribe(
      res => {
        this.registryMode.push({key: Reflect.get(res, 'NodeCreateNew.Auto'), value: 'auto'});
        this.registryMode.push({key: Reflect.get(res, 'NodeCreateNew.Manual'), value: 'manual'});
      }
    );
    this.isActionWip = true;
    this.nodeList.forEach(node => {
      if (node.nodeType === 'master') {
        this.masters.push(node.nodeName);
      }
    });
    this.isActionWip = false;
  }

  ngOnDestroy() {
    if (this.subRefreshInterval) {
      this.subRefreshInterval.unsubscribe();
    }
    super.ngOnDestroy();
  }

  get checkEdgeNodeIpFun() {
    return this.checkEdgeNodeIp.bind(this);
  }

  get checkEdgeNodeNameFun() {
    return this.checkEdgeNodeName.bind(this);
  }

  cpuTypeDisableFun(item: { description: string }): boolean {
    return item.description !== 'auto-detect';
  }

  checkEdgeNodeIp(control: HTMLInputElement): Observable<ValidationErrors | null> {
    const existsNode = this.nodeList.find(node => node.nodeName === control.value);
    return existsNode ? of({edgeNodeIpExists: 'NodeCreateNew.IpExists'}) : of(null);
  }

  checkEdgeNodeName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.nodeService.checkNodeGroupExist(control.value)
      .pipe(
        map(() => null),
        catchError((err: HttpErrorResponse) => {
          this.messageService.cleanNotification();
          if (err.status === 409) {
            return of({edgeNodeNameExists: 'NodeCreateNew.NameExists'});
          } else {
            return of(null);
          }
        })
      );
  }

  setRegisterMode(register: { key: string, value: string }) {
    this.edgeNode.registryMode = register.value;
  }

  setCpuType(type: { description: string }) {
    this.edgeNode.cpuType = type.description;
  }

  checkNodeBuildStatus() {
    this.nodeService.getNodes().subscribe(
      (res: Array<NodeStatus>) => {
        res.forEach(value => {
          if (value.nodeName === this.edgeNode.name) {
            this.messageService.showAlert('NodeCreateNew.AddSuccessfully');
            this.modalOpened = false;
          }
          if (this.checkTimes > 5) {
            this.messageService.showAlert('NodeCreateNew.TimeOutMessage', {alertType: 'danger'});
            this.modalOpened = false;
          }
          this.checkTimes += 1;
        });
      },
      () => this.modalOpened = false
    );
  }

  addEdgeNode() {
    if (this.verifyInputExValid() && this.verifyDropdownExValid()) {
      this.isActionWip = true;
      this.nodeService.addEdgeNode(this.edgeNode).subscribe(
        () => this.subRefreshInterval = interval(10000).subscribe(() => this.checkNodeBuildStatus()),
        (err: HttpErrorResponse | TimeoutError) => {
          this.translateService.get('NodeCreateNew.ParamsErrorMessage').subscribe(msg => {
              if (err instanceof HttpErrorResponse) {
                if (err.status === 400) {
                  this.messageService.cleanNotification();
                  this.messageService.showAlert(msg, {alertType: 'danger'});
                }
              } else {
                this.messageService.showAlert(msg, {alertType: 'danger'});
              }
              this.modalOpened = false;
            }
          );
        }
      );
    }
  }
}
