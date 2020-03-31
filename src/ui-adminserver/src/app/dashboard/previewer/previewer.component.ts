import { Component, OnInit, OnDestroy } from '@angular/core';
import { ComponentStatus } from '../component-status.model';
import { DashboardService } from '../dashboard.service';
import { HttpErrorResponse } from '@angular/common/http';
import { User } from 'src/app/account/account.model';
import { Router } from '@angular/router';

@Component({
  selector: 'app-previewer',
  templateUrl: './previewer.component.html',
  styleUrls: ['./previewer.component.css']
})

export class PreviewerComponent implements OnInit, OnDestroy {
  componentList: ComponentStatus[];
  showDetail = false;
  modal: ComponentStatus;
  confirmModal = false;
  confirmType: ConfirmType;
  timer: any;
  user: User;
  loadingFlag = true;
  enableStop = false;
  disableApply = false;

  constructor(private dashboardService: DashboardService,
    private router: Router) {
    this.modal = new ComponentStatus();
    this.confirmType = new ConfirmType('rb');
    this.user = new User();
  }

  ngOnInit() {
    this.getMonitor();
  }

  getMonitor() {
    // 10s 刷新一次
    clearInterval(this.timer);
    this.timer = setInterval(
      () => {
        this.dashboardService.monitorContainer().subscribe(
          (res: Array<ComponentStatus>) => {
            this.componentList = res;
            this.loadingFlag = false;
            this.enableStop = this.componentList.length > 3;
            this.reflashDetail();
          },
          (err: HttpErrorResponse) => {
            this.loadingFlag = false;
            this.commonError(err);
            clearInterval(this.timer); // 销毁定时器
          }
        );
      }, 10000);
  }

  ngOnDestroy() {
    if (this.timer) {
      clearInterval(this.timer); // 销毁定时器
    }
  }

  getDetail(item: ComponentStatus) {
    this.modal = item;
    this.showDetail = true;
  }

  reflashDetail() {
    if (this.modal.id) {
      for (let item of this.componentList) {
        if (this.modal.id === item.id) {
          this.modal = item;
          break;
        }
      }
    }
  }

  confirm(type: string, containerID?: string) {
    this.confirmModal = true;
    if (containerID) {
      this.confirmType = new ConfirmType(type, containerID);
    } else {
      this.confirmType = new ConfirmType(type);
    }
  }

  boardControl(type: string, containerID?: string) {
    this.loadingFlag = true;
    this.disableApply = true;
    if (type === 'rb') {
      this.dashboardService.restartBoard(this.user).subscribe(
        () => {
          this.disableApply = false;
          this.confirmModal = false;
          this.user = new User();
          this.getMonitor();
          alert('Waiting for restart.');
        },
        (err: HttpErrorResponse) => {
          this.loadingFlag = false;
          this.disableApply = false;
          this.commonError(err);
        }
      );
    } else if (type === 'rc') {
      this.loadingFlag = false;
      this.disableApply = false;
      this.confirmModal = false;
      alert('Sorry, this feature is not yet supported. Restart container(' + containerID + ') fail.');
    } else if (type === 'sb') {
      this.dashboardService.shutdownBoard(this.user).subscribe(
        () => {
          this.disableApply = false;
          this.confirmModal = false;
          this.user = new User();
          this.getMonitor();
          alert('Waiting for STOP.');
        },
        (err: HttpErrorResponse) => {
          this.loadingFlag = false;
          this.disableApply = false;
          this.commonError(err);
        }
      );
    } else {
      this.loadingFlag = false;
      this.disableApply = false;
      this.confirmModal = false;
      alert('Wrong parameter!');
    }
  }

  commonError(err: HttpErrorResponse) {
    if (err.status === 401) {
      alert('User status error! Please login again!');
      this.router.navigateByUrl('account/login');
    } else {
      alert('Unknown Error!');
    }
  }
}

class ConfirmType {
  title: string;
  comment: string;
  button: string;
  type: string;
  containerId = '';

  constructor(type: string, containerID?: string, title?: string, comment?: string, button?: string, ) {
    this.type = type;
    const currentLang = (window.localStorage.getItem('currentLang') === 'zh-cn' || window.localStorage.getItem('currentLang') === 'zh');
    if (currentLang) {
      this.button = '重启';
      if (type === 'rb') {
        this.title = '重启Board';
        this.comment = '您确定要重新启动Board吗？如果是这样，请输入主机的帐户和密码：';
      } else if (type === 'rc') {
        this.title = '重启容器';
        this.comment = '请输入主机的帐户和密码以重新启动容器：';
        this.containerId = containerID;
      } else if (type === 'sb') {
        this.title = '停止Board';
        this.comment = '您确定要停止Board吗？如果是这样，请输入主机的帐户和密码：';
        this.button = '停止';
      } else {
        this.title = title ? title : 'Title';
        this.comment = comment ? comment : 'Comment';
        this.button = button ? button : 'Button';
      }
    } else {
      this.button = 'restart';
      if (type === 'rb') {
        this.title = 'Restart Board?';
        this.comment = 'Are you sure to RESTART the Board? If so, please enter the account and password of the host machine.';
      } else if (type === 'rc') {
        this.title = 'Restart Container?';
        this.comment = 'Please enter the account and password of the host machine to Restart the Container:';
        this.containerId = containerID;
      } else if (type === 'sb') {
        this.title = 'Stop Board?';
        this.comment = 'Are you sure to STOP the Board? If so, please enter the account and password of the host machine:';
        this.button = 'STOP';
      } else {
        this.title = title ? title : 'Title';
        this.comment = comment ? comment : 'Comment';
        this.button = button ? button : 'Button';
      }
    }
  }
}
