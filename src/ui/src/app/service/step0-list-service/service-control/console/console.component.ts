import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
  ViewContainerRef
} from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SearchAddon } from 'xterm-addon-search';
import { HttpErrorResponse, HttpEvent, HttpEventType, HttpProgressEvent, HttpResponse } from '@angular/common/http';
import { AppInitService } from '../../../../shared.service/app-init.service';
import { K8sService } from '../../../service.k8s';
import { Service, ServiceContainer, ServiceDetailInfo } from '../../../service.types';
import { MessageService } from '../../../../shared.service/message.service';
import { GlobalAlertType } from '../../../../shared/shared.types';

@Component({
  selector: 'app-console',
  templateUrl: './console.component.html',
  styleUrls: ['./console.component.css']
})
export class ConsoleComponent implements OnInit, AfterViewInit, OnDestroy {
  @Input() service: Service;
  @ViewChild('terminalContainer') terminalContainer: ElementRef;
  @ViewChild('messageContainer', {read: ViewContainerRef}) messageContainer: ViewContainerRef;
  @Output() actionIsEnabledEvent: EventEmitter<boolean>;
  @Output() updateProgressEvent: EventEmitter<HttpProgressEvent>;
  @Input() isActionInWIP: boolean;
  @Output() isActionInWIPChange: EventEmitter<boolean>;
  @Output() errorEvent: EventEmitter<any>;
  term: Terminal;
  fitAddon: FitAddon;
  webLinkAddon: WebLinksAddon;
  searchAddon: SearchAddon;
  curPodName = '';
  curContainerName = '';
  ws: WebSocket;
  serviceDetailInfo: ServiceDetailInfo;
  curActiveIndex = -1;
  curDownLoadPath = '';
  curUploadFile: File;
  curUploadPath = '';
  curReadyState = 0;

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService,
              private messageService: MessageService,
              private changeRef: ChangeDetectorRef,
              private translateService: TranslateService) {
    this.fitAddon = new FitAddon();
    this.searchAddon = new SearchAddon();
    this.webLinkAddon = new WebLinksAddon(this.webLinksHandle);
    this.serviceDetailInfo = new ServiceDetailInfo();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
    this.updateProgressEvent = new EventEmitter<HttpProgressEvent>();
    this.isActionInWIPChange = new EventEmitter<boolean>();
    this.errorEvent = new EventEmitter<any>();
  }

  ngOnInit() {
    this.k8sService.getServiceDetail(this.service.serviceId).subscribe(
      (res: ServiceDetailInfo) => {
        this.serviceDetailInfo = res;
        const containers = this.serviceDetailInfo.serviceContainers;
        if (containers.length > 0) {
          this.serviceDetailInfo.serviceContainers = containers.filter(value =>
            value.securityContext === false && value.initContainer === false);
        }
        if (this.serviceDetailInfo.serviceContainers.length > 0) {
          this.buildSocketConnect(this.serviceDetailInfo.serviceContainers[0], 0);
        }
        this.changeRef.detectChanges();
      },
      (err) => this.errorEvent.emit(err)
    );
    this.resizeListener = this.resizeListener.bind(this);
    this.actionIsEnabledEvent.emit(true);
  }

  ngOnDestroy(): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.close();
    }
    window.removeEventListener('resize', this.resizeListener);
  }

  ngAfterViewInit(): void {
    window.addEventListener('resize', this.resizeListener);
  }

  get hasContainerToConnect(): boolean {
    return this.serviceDetailInfo.serviceContainers.length > 0;
  }

  get noneContainerToConnect(): boolean {
    return this.serviceDetailInfo.serviceContainers.length === 0;
  }

  get wsUrl(): string {
    const boardHost = this.appInitService.systemInfo.boardHost;
    const host = `${this.appInitService.getWebsocketPrefix}://${boardHost}:${window.location.port}`;
    const path = `/api/v1/pods/${this.service.serviceProjectId}/${this.curPodName}/shell`;
    const params = `?token=${this.appInitService.token}&container=${this.curContainerName}`;
    return `${host}${path}${params}`;
  }

  get status(): string {
    if (this.ws) {
      switch (this.ws.readyState) {
        case WebSocket.OPEN:
          return 'ServiceControlConsole.Open';
        case WebSocket.CLOSED:
          return 'ServiceControlConsole.Closed';
        case WebSocket.CLOSING:
          return 'ServiceControlConsole.Closing';
        case WebSocket.CONNECTING:
          return 'ServiceControlConsole.Connecting';
        default:
          return 'ServiceControlConsole.Unknown';
      }
    } else {
      return 'ServiceControlConsole.Unknown';
    }
  }

  get statusStyle(): { [key: string]: string } {
    if (this.ws) {
      switch (this.ws.readyState) {
        case WebSocket.OPEN:
          return {color: 'green'};
        case WebSocket.CLOSED:
          return {color: 'red'};
        case WebSocket.CLOSING:
          return {color: 'yellow'};
        case WebSocket.CONNECTING:
          return {color: 'lightgreen'};
        default:
          return {color: 'black'};
      }
    } else {
      return {color: 'black'};
    }
  }

  buildSocketConnect(serviceContainer: ServiceContainer, index: number) {
    if (this.isActionInWIP) {
      return;
    }
    this.curActiveIndex = index;
    this.curPodName = serviceContainer.podName;
    this.curContainerName = serviceContainer.containerName;
    this.ws = new WebSocket(this.wsUrl);
    this.mountWebSocket();
  }

  createTerm() {
    this.term = new Terminal({
      cursorBlink: true,
      disableStdin: false,
      cursorStyle: 'block',
      cols: 35,
      rows: 20,
    });

  }

  initTerm() {
    this.term.loadAddon(this.webLinkAddon);
    this.term.loadAddon(this.searchAddon);
    this.term.loadAddon(this.fitAddon);
    const terminalContainerElement = (this.terminalContainer.nativeElement as HTMLElement);
    while (terminalContainerElement.firstChild) {
      terminalContainerElement.firstChild.remove();
    }
    this.term.open(terminalContainerElement);
    this.term.focus();
    this.fitAddon.fit();
  }

  mountTerm() {
    this.term.onData((arg1: string, arg2: any): any => {
      if (this.ws.readyState === WebSocket.OPEN) {
        const msg = {type: 'input', input: arg1};
        this.ws.send(JSON.stringify(msg));
      }
    });
  }

  mountWebSocket() {
    this.ws.onopen = (ev: Event): any => {
      this.curReadyState = 1;
      this.createTerm();
      this.initTerm();
      this.mountTerm();
      const msg = {type: 'resize', rows: this.term.rows, cols: this.term.cols};
      this.ws.send(JSON.stringify(msg));
    };

    this.ws.onclose = (ev: CloseEvent): any => {
      this.curActiveIndex = -1;
      this.curReadyState = 2;
    };

    this.ws.onmessage = (ev: MessageEvent): any => {
      this.term.write(ev.data);
    };

    this.ws.onerror = (ev: Event): any => {
      this.curActiveIndex = -1;
      this.curReadyState = 3;
      this.createTerm();
      this.initTerm();
      this.mountTerm();
      this.translateService.get('ServiceControlConsole.WebsocketConnectionError').subscribe(
        res => this.messageService.showGlobalMessage(res, {alertType: 'danger', view: this.messageContainer})
      );
    };
  }

  webLinksHandle(event: MouseEvent, uri: string): void {
    // Todo: enhancement
  }

  resizeListener(event: Event) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.fitAddon.fit();
      const msg = {type: 'resize', rows: this.term.rows, cols: this.term.cols};
      this.ws.send(JSON.stringify(msg));
    }
  }

  changeUploadFile(event: Event) {
    if (this.ws.readyState !== WebSocket.OPEN) {
      return;
    }
    const fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      this.curUploadFile = fileList[0];
    }
  }

  uploadFile() {
    if (this.isActionInWIP) {
      return;
    }
    if (this.ws.readyState !== WebSocket.OPEN) {
      return;
    }
    if (this.curUploadPath === '') {
      this.translateService.get('ServiceControlConsole.UploadFilePath').subscribe(
        res => this.messageService.showAlert(res, {alertType: 'warning', view: this.messageContainer})
      );
      return;
    }
    if (!this.curUploadFile) {
      this.translateService.get('ServiceControlConsole.UploadFile').subscribe(
        res => this.messageService.showAlert(res, {alertType: 'warning', view: this.messageContainer})
      );
      return;
    }
    this.isActionInWIPChange.emit(true);
    this.messageService.cleanNotification();
    const formData = new FormData();
    formData.append('upload_file', this.curUploadFile, this.curUploadFile.name);
    this.k8sService.uploadFile(this.service.serviceProjectId, this.curPodName, this.curContainerName, this.curUploadPath, formData)
      .subscribe((res: HttpEvent<any>) => {
          if (res.type === HttpEventType.UploadProgress || res.type === HttpEventType.DownloadProgress) {
            this.updateProgressEvent.emit(res);
          } else if (res.type === HttpEventType.Response) {
            this.translateService.get('ServiceControlConsole.UploadFileSuccess').subscribe(
              msg => this.messageService.showAlert(msg, {view: this.messageContainer})
            );
            this.isActionInWIPChange.emit(false);
          }
        },
        () => this.isActionInWIPChange.emit(false)
      );
  }

  downloadFile() {
    if (this.isActionInWIP) {
      return;
    }
    if (this.ws.readyState !== WebSocket.OPEN) {
      return;
    }
    if (this.curDownLoadPath === '') {
      this.translateService.get('ServiceControlConsole.downloadFilePath').subscribe(
        res => this.messageService.showAlert(res, {alertType: 'warning', view: this.messageContainer})
      );
      return;
    }
    this.isActionInWIPChange.emit(true);
    this.messageService.cleanNotification();
    this.k8sService.downloadFile(this.service.serviceProjectId, this.curPodName, this.curContainerName, this.curDownLoadPath).subscribe(
      (res: HttpEvent<any>) => {
        if (res.type === HttpEventType.DownloadProgress) {
          this.updateProgressEvent.emit(res);
        } else if (res instanceof HttpResponse) {
          const file = new Blob([res.body], {type: res.body.type});
          const url = URL.createObjectURL(file);
          const link = document.createElement('a');
          const downLoadPaths = this.curDownLoadPath.split('/');
          link.setAttribute('style', 'display:none');
          link.setAttribute('href', url);
          link.setAttribute('target', '_blank');
          link.setAttribute('download', downLoadPaths[downLoadPaths.length - 1]);
          link.click();
          this.isActionInWIPChange.emit(false);
        }
      },
      (err: HttpErrorResponse) => {
        this.messageService.cleanNotification();
        this.isActionInWIPChange.emit(false);
        this.translateService.get('ServiceControlConsole.downloadFileFailed').subscribe(
          res => this.messageService.showGlobalMessage(res, {
            globalAlertType: GlobalAlertType.gatShowDetail,
            view: this.messageContainer,
            errorObject: err
          })
        );
      }
    );
  }
}
