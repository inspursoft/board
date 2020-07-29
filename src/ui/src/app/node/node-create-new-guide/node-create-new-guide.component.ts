import {
  Component,
  ComponentFactoryResolver,
  ElementRef,
  EventEmitter,
  Input,
  OnInit,
  Output,
  ViewChild,
  ViewContainerRef
} from '@angular/core';
import { AddNodeType, NodeStatus } from '../node.types';
import { HttpErrorResponse } from '@angular/common/http';
import { AppInitService } from '../../shared.service/app-init.service';
import { MessageService } from '../../shared.service/message.service';
import { NodeCreateNewComponent } from '../node-create-new/node-create-new.component';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';

@Component({
  selector: 'app-node-create-new-guide',
  templateUrl: './node-create-new-guide.component.html',
  styleUrls: ['./node-create-new-guide.component.css']
})
export class NodeCreateNewGuideComponent extends CsModalParentBase implements OnInit {
  @ViewChild('addNormalLink') addNormalLink: ElementRef;
  @Output() createNodeFinished: EventEmitter<any>;
  @Input() nodeList: Array<NodeStatus>;
  addNodeType = AddNodeType.normal;
  addNormalNodeUrl = '';
  isAdminSeverWorking = false;

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              public factoryResolver: ComponentFactoryResolver,
              public selfView: ViewContainerRef) {
    super(factoryResolver, selfView);
    this.nodeList = Array<NodeStatus>();
    this.createNodeFinished = new EventEmitter();
  }

  ngOnInit() {
    this.addNormalNodeUrl = `http://${this.appInitService.systemInfo.board_host}:8082/resource/node-list`;
    this.checkAdminServer();
  }

  get btnOkDisable(): boolean {
    if (this.addNodeType === AddNodeType.normal) {
      return !this.isAdminSeverWorking;
    } else {
      return false;
    }
  }

  checkAdminServer() {
    this.appInitService.getIsShowAdminServer().subscribe(
      () => this.isAdminSeverWorking = !this.appInitService.isArmSystem &&
        !this.appInitService.isMipsSystem,
      (err: HttpErrorResponse) => {
        this.messageService.cleanNotification();
        if (err.status === 401) {
          this.isAdminSeverWorking = !this.appInitService.isArmSystem &&
            !this.appInitService.isMipsSystem;
        }
      }
    );
  }

  setAddNodeType(type: AddNodeType) {
    this.addNodeType = type;
  }

  cancelAdd() {
    this.createNodeFinished.emit(true);
  }

  addNode() {
    if (this.addNodeType === AddNodeType.edge) {
      const instance = this.createNewModal(NodeCreateNewComponent);
      instance.nodeList = this.nodeList;
      instance.closeNotification.subscribe(() => this.createNodeFinished.emit(true));
    } else {
      (this.addNormalLink.nativeElement as HTMLLinkElement).click();
    }
  }

}
