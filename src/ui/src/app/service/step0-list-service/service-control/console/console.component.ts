import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef, EventEmitter,
  Input,
  OnDestroy,
  OnInit, Output,
  ViewChild
} from '@angular/core';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SearchAddon } from 'xterm-addon-search';
import { AppInitService } from '../../../../shared.service/app-init.service';
import { K8sService } from '../../../service.k8s';
import { Service, ServiceContainer, ServiceDetailInfo } from '../../../service.types';

@Component({
  selector: 'app-console',
  templateUrl: './console.component.html',
  styleUrls: ['./console.component.css']
})
export class ConsoleComponent implements OnInit, AfterViewInit, OnDestroy {
  @Input() service: Service;
  @ViewChild('terminalContainer') terminalContainer: ElementRef;
  @Output() actionIsEnabledEvent: EventEmitter<boolean>;
  term: Terminal;
  fitAddon: FitAddon;
  webLinkAddon: WebLinksAddon;
  searchAddon: SearchAddon;
  curPodName = '';
  curContainerName = '';
  ws: WebSocket;
  serviceDetailInfo: ServiceDetailInfo;
  curActiveIndex = -1;

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService,
              private changeRef: ChangeDetectorRef) {
    this.fitAddon = new FitAddon();
    this.searchAddon = new SearchAddon();
    this.webLinkAddon = new WebLinksAddon(this.webLinksHandle);
    this.serviceDetailInfo = new ServiceDetailInfo();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.k8sService.getServiceDetail(this.service.serviceId).subscribe(
      (res: ServiceDetailInfo) => {
        this.serviceDetailInfo = res;
        this.changeRef.detectChanges();
        if (this.serviceDetailInfo.serviceContainers.length > 0) {
          this.buildSocketConnect(this.serviceDetailInfo.serviceContainers[0], 0);
        }
      }
    );
    this.resizeListener = this.resizeListener.bind(this);
    this.actionIsEnabledEvent.emit(true);
  }

  ngOnDestroy(): void {
    this.ws.close();
    window.removeEventListener('resize', this.resizeListener);
  }

  ngAfterViewInit(): void {
    window.addEventListener('resize', this.resizeListener);
  }

  get wsUrl(): string {
    const host = `wss://${this.appInitService.systemInfo.board_host}`;
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
      cols: 59,
      rows: 25,
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
      this.createTerm();
      this.initTerm();
      this.mountTerm();
      const msg = {type: 'resize', rows: this.term.rows, cols: this.term.cols};
      this.ws.send(JSON.stringify(msg));
    };

    this.ws.onclose = (ev: CloseEvent): any => {
      this.curActiveIndex = -1;
    };

    this.ws.onmessage = (ev: MessageEvent): any => {
      this.term.write(ev.data);
    };

    this.ws.onerror = (ev: Event): any => {
      this.curActiveIndex = -1;
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
}
