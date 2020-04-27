import { AfterViewInit, Component, ElementRef, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
import { SearchAddon } from 'xterm-addon-search';
import { Service } from '../../../service';
import { AppInitService } from '../../../../shared.service/app-init.service';
import { K8sService } from '../../../service.k8s';

@Component({
  selector: 'app-console',
  templateUrl: './console.component.html',
  styleUrls: ['./console.component.css']
})
export class ConsoleComponent implements OnInit, AfterViewInit, OnDestroy {
  @Input() service: Service;
  @Input() podName = '';
  @Input() containerName = '';
  @ViewChild('terminalContainer') terminalContainer: ElementRef;
  term: Terminal;
  fitAddon: FitAddon;
  webLinkAddon: WebLinksAddon;
  searchAddon: SearchAddon;
  ws: WebSocket;

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService) {
    this.term = new Terminal({
      cursorBlink: true,
      disableStdin: false,
      cursorStyle: 'block'
    });
    this.fitAddon = new FitAddon();
    this.searchAddon = new SearchAddon();
    this.webLinkAddon = new WebLinksAddon(this.webLinksHandle);
    // this.ws = new WebSocket(this.wsUrl);
  }

  ngOnInit() {
    this.k8sService.getServiceDetail(this.service.service_id).subscribe(
      res => {
        console.log(res);
      }
    );
    this.term.loadAddon(this.webLinkAddon);
    this.term.loadAddon(this.searchAddon);
    this.term.loadAddon(this.fitAddon);
    this.term.open(this.terminalContainer.nativeElement);
    this.resizeListener = this.resizeListener.bind(this);
    // this.mountWebSocket();
  }

  ngOnDestroy(): void {
    window.removeEventListener('resize', this.resizeListener);
  }

  ngAfterViewInit(): void {
    this.term.focus();
    this.fitAddon.fit();
    window.addEventListener('resize', this.resizeListener);
  }

  get wsUrl(): string {
    const host = `ws://${this.appInitService.systemInfo.board_host}`;
    const path = `/api/v1/pods/${this.service.service_project_name}/${this.podName}/shell`;
    const params = `?token=${this.appInitService.token}&container=${this.containerName}`;
    return `${host}${path}${params}`;
  }

  mountWebSocket() {
    this.ws.onopen = (ev: Event): any => {
      const msg = {type: 'resize', rows: this.term.rows, cols: this.term.cols};
      this.ws.send(JSON.stringify(msg));
    };

    this.ws.onclose = (ev: CloseEvent): any => {
      console.log('ws closed');
    };

    this.ws.onmessage = (ev: MessageEvent): any => {
      this.term.write(ev.data);
    };

    this.ws.onerror = (ev: Event): any => {
      console.log(`ws error:${ev}`);
    };

    this.term.onData((arg1: string, arg2: any): any => {
      const msg = {type: 'input', input: arg1};
      this.ws.send(JSON.stringify(msg));
    });
  }

  webLinksHandle(event: MouseEvent, uri: string): void {
    // console.log(`onWebLinks event event:${event}`);
    // console.log(`onWebLinks event uri:${uri}`);
  }

  resizeListener(event: Event) {
    this.fitAddon.fit();
    // const msg = {type: 'resize', rows: this.term.rows, cols: this.term.cols};
    // this.ws.send(JSON.stringify(msg));
  }
}
