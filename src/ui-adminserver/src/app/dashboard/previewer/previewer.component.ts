import { Component, OnInit, OnDestroy } from '@angular/core';
import { ComponentStatus } from '../component-status.model';
import { DashboardService } from '../dashboard.service';
import { map } from 'rxjs/operators';
import { HttpErrorResponse } from '@angular/common/http';

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
  token: string;
  timer;

  constructor(private dashboardService: DashboardService) {
    this.modal = new ComponentStatus();
    this.confirmType = new ConfirmType('rb');
    this.token = sessionStorage.getItem('token');
  }

  ngOnInit() {
    // 10s 刷新一次
    this.timer = setInterval(
      () => {
        this.dashboardService.monitorContainer(this.token).subscribe(
          (res: Array<ComponentStatus>) => {
            this.componentList = res;
          },
          (err: HttpErrorResponse) => {
            alert("no resource");
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

  confirm(type: string, containerID?: string) {
    this.confirmModal = true;
    if (containerID) {
      this.confirmType = new ConfirmType(type, containerID);
    } else {
      this.confirmType = new ConfirmType(type);
    }
  }

  boardControl(type: string, containerID?: string) {
    if (type === 'rb') {
      this.dashboardService.restartBoard(encodeURIComponent(this.token)).subscribe(
        () => {
          this.confirmModal = false;
          alert('Waiting for restart.');
        },
        () => {
          this.confirmModal = false;
          alert('Unknown Error');
        }
      );
    } else if (type === 'rc') {
      this.confirmModal = false;
      alert('Sorry, this feature is not yet supported. Restart container(' + containerID + ') fail.');
    } else if (type === 'sb') {
      this.dashboardService.shutdownBoard(encodeURIComponent(this.token)).subscribe(
        () => {
          this.confirmModal = false;
          alert('Waiting for STOP.');
        },
        () => {
          this.confirmModal = false;
          alert('Unknown Error');
        }
      );
    } else {
      this.confirmModal = false;
      alert('Wrong parameter!');
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
    if (type === 'rb') {
      this.title = 'Restart Board?';
      this.comment = 'Are you sure to RESTART the Board?';
      this.button = 'restart';
    } else if (type === 'rc') {
      this.title = 'Restart Container?';
      this.comment = 'Are you sure to Restart the Container:';
      this.button = 'restart';
      this.containerId = containerID;
    } else if (type === 'sb') {
      this.title = 'Stop Board?';
      this.comment = 'Are you sure to STOP the Board?';
      this.button = 'STOP';
    } else {
      this.title = title ? title : 'Title';
      this.comment = comment ? comment : 'Comment';
      this.button = button ? button : 'Button';
    }
  }
}
