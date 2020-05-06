import { Component, OnInit, ViewChild, ViewChildren, QueryList } from '@angular/core';
import { AccountService } from '../account.service';
import { HttpErrorResponse } from '@angular/common/http';
import { MessageService } from 'src/app/shared/message/message.service';
import { MyInputTemplateComponent } from 'src/app/shared/my-input-template/my-input-template.component';
import { ClrLoadingState } from '@clr/angular';
import { AppInitService } from 'src/app/shared.service/app-init.service';
import { InitStatus, InitStatusCode } from 'src/app/shared.service/app-init.model';
import { User } from '../account.model';
import { BoardService } from 'src/app/shared.service/board.service';
import { ConfigurationService } from 'src/app/shared.service/configuration.service';
import { Configuration } from 'src/app/shared.service/configuration.model';
import { Router } from '@angular/router';
import { Message, ReturnStatus } from 'src/app/shared/message/message.types';

@Component({
  selector: 'app-installation',
  templateUrl: './installation.component.html',
  styleUrls: ['./installation.component.css']
})
export class InstallationComponent implements OnInit {
  newDate = new Date('2016-01-01 09:00:00');

  passwordPattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)([A-Za-z\d#?!@$%^&*-]){8,20}$/;
  uuidPattern = /^\w{8}(-\w{4}){3}-\w{12}$/;

  uuid = '';
  // TODO email identity不为空时，无法映射默认值
  config: Configuration;
  showBaselineHelper = false;

  installStep = 0;
  ignoreStep1 = false;
  ignoreStep2 = false;
  installProgress = 0;
  enableBtn = true;
  refresh = false;

  // TODO put some config into default.
  simpleMode = false;
  loadingFlag = true;
  disconnect = false;
  enableInitialization = false;
  openSSH = false;
  uninstallConfirm = false;
  clearDate = false;
  responsibility = false;
  isEditable = false;
  submitBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  user: User;

  @ViewChild('UUID') uuidInput: MyInputTemplateComponent;
  @ViewChildren(MyInputTemplateComponent) myInputTemplateComponents: QueryList<MyInputTemplateComponent>;

  constructor(private accountService: AccountService,
              private appInitService: AppInitService,
              private boardService: BoardService,
              private configurationService: ConfigurationService,
              private messageService: MessageService,
              private router: Router) {
    this.user = new User();
  }

  ngOnInit() {
    this.appInitService.getSystemStatus().subscribe(
      (res: InitStatus) => {
        if (InitStatusCode.InitStatusThird === res.status) {
          this.messageService.showOnlyOkDialogObservable('INITIALIZATION.ALERTS.ALREADY_START').subscribe(
            (msg: Message) => {
              if (msg.returnStatus === ReturnStatus.rsConfirm) {
                this.router.navigateByUrl('account');
              }
            }
          );
        } else {
          this.accountService.createUUID().subscribe(
            () => {
              this.loadingFlag = false;
              this.enableInitialization = true;
            },
            (err: HttpErrorResponse) => {
              this.loadingFlag = false;
              this.disconnect = true;
              console.error(err.message);
              this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.INITIALIZATION', 'ACCOUNT.ERROR');
              this.refresh = true;
            });
        }
      },
      (err: HttpErrorResponse) => {
        console.error(err.message);
        this.getSysStatusFailed();
      }
    );

  }

  onNext() {
    // for test
    // test status1
    // this.config = new Configuration();
    // setTimeout(() => {
    //   this.ignoreStep1 = true;
    //   this.installStep = 2;
    //   this.installProgress = 50;
    // }, 1000)

    // test status2
    // this.installStep++;
    // this.installProgress += 33;

    // test status3
    // this.ignoreStep1 = true;
    // this.ignoreStep2 = true;
    // this.installStep = 3;
    // this.installProgress = 100;
    // this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.ALREADY_START');

    this.uuidInput.checkSelf();
    if (this.uuidInput.isValid) {
      this.submitBtnState = ClrLoadingState.LOADING;
      this.user.username = 'admin';
      this.user.password = this.uuid;
      this.accountService.signIn(this.user).subscribe(
        () => {
          this.user = new User();
          sessionStorage.setItem('token', this.uuid);
          this.appInitService.getSystemStatus().subscribe(
            (res: InitStatus) => {
              switch (res.status) {
                // 未起Board且未更改cfg
                case InitStatusCode.InitStatusFirst: {
                  this.configurationService.getConfig().subscribe(
                    (resTmp: Configuration) => {
                      this.config = new Configuration(resTmp);
                      this.isEditable = true;
                      this.ignoreStep1 = true;
                      this.installStep = 2;
                      this.installProgress = 50;
                      this.submitBtnState = ClrLoadingState.DEFAULT;
                    },
                    (err: HttpErrorResponse) => {
                      // COMMON
                      this.commonError(err, new Map(), 'INITIALIZATION.ALERTS.GET_CFG_FAILED');
                    }
                  );
                  break;
                }
                // 未起Board但更改过cfg
                case InitStatusCode.InitStatusSecond: {
                  this.installStep++;
                  this.installProgress += 33;
                  this.submitBtnState = ClrLoadingState.DEFAULT;
                  break;
                }
                // Board已经运行
                case InitStatusCode.InitStatusThird: {
                  this.ignoreStep1 = true;
                  this.ignoreStep2 = true;
                  this.installStep = 3;
                  this.installProgress = 100;
                  this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.ALREADY_START');
                  this.submitBtnState = ClrLoadingState.DEFAULT;
                  break;
                }
              }
            },
            (err: HttpErrorResponse) => {
              console.error(err.message);
              this.getSysStatusFailed();
            }
          );
        },
        (err: HttpErrorResponse) => {
          console.error(err.message);
          this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.VALIDATE_UUID_FAILED', 'ACCOUNT.ERROR');
          this.refresh = true;
          this.submitBtnState = ClrLoadingState.DEFAULT;
        }
      );
    }

  }

  onEditCfg() {
    // for test
    // this.installStep++;
    // this.installProgress += 33;


    this.submitBtnState = ClrLoadingState.LOADING;
    this.configurationService.getConfig().subscribe(
      (res: Configuration) => {
        this.config = new Configuration(res);
        if (this.config.tmpExist) {
          this.configurationService.getConfig('tmp').subscribe(
            (resTmp: Configuration) => {
              this.config = new Configuration(resTmp);
              this.newDate = new Date(this.config.apiserver.imageBaselineTime);
              this.isEditable = this.config.isInit;
              this.installStep++;
              this.installProgress += 33;
              this.submitBtnState = ClrLoadingState.DEFAULT;
            },
            (err: HttpErrorResponse) => {
              if (err.status === 401) {
                this.tokenError();
              } else {
                console.log('Can not read tmp file: ' + err.message);
                this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.GET_TMP_FAILED');
                this.newDate = new Date(this.config.apiserver.imageBaselineTime);
                this.isEditable = this.config.isInit;
                this.installStep++;
                this.installProgress += 33;
                this.submitBtnState = ClrLoadingState.DEFAULT;
              }
            }
          );
        } else {
          this.newDate = new Date(this.config.apiserver.imageBaselineTime);
          this.isEditable = this.config.isInit;
          this.installStep++;
          this.installProgress += 33;
          this.submitBtnState = ClrLoadingState.DEFAULT;
        }
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.commonError(err, new Map(), 'INITIALIZATION.ALERTS.GET_CFG_FAILED');
      },
    );
  }

  onStartBoard() {
    // for test
    // this.openSSH = false;
    // this.installStep += 2;
    // this.ignoreStep2 = true;
    // this.installProgress = 100;

    this.submitBtnState = ClrLoadingState.LOADING;
    this.openSSH = false;
    this.boardService.start(this.user).subscribe(
      () => {
        this.installStep += 2;
        this.ignoreStep2 = true;
        this.installProgress = 100;
        this.submitBtnState = ClrLoadingState.DEFAULT;
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.commonError(err, new Map(), 'INITIALIZATION.ALERTS.START_BOARD_FAILED');
      },
    );
  }

  onUninstallBoard() {
    // for test
    // this.openSSH = false;
    // this.installStep = 4;
    // this.ignoreStep1 = true;
    // this.ignoreStep2 = true;
    // this.installProgress = 100;


    this.submitBtnState = ClrLoadingState.LOADING;
    this.openSSH = false;
    this.boardService.shutdown(this.user, this.clearDate).subscribe(
      () => {
        this.installStep = 4;
        this.ignoreStep1 = true;
        this.ignoreStep2 = true;
        this.installProgress = 100;
        this.submitBtnState = ClrLoadingState.DEFAULT;
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.commonError(err, new Map(), 'INITIALIZATION.ALERTS.UNINSTALL_BOARD_FAILED');
      },
    );
  }

  onApplyAndStartBoard() {
    // for test
    // this.openSSH = false;
    // this.installStep++;
    // this.installProgress = 100;



    this.submitBtnState = ClrLoadingState.LOADING;
    this.openSSH = false;
    this.configurationService.putConfig(this.config).subscribe(
      () => {
        this.boardService.applyCfg(this.user).subscribe(
          () => {
            this.installStep++;
            this.installProgress = 100;
            this.submitBtnState = ClrLoadingState.DEFAULT;
          },
          (err: HttpErrorResponse) => {
            // COMMON
            this.commonError(err, new Map(), 'INITIALIZATION.ALERTS.START_BOARD_FAILED');
          },
        );
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.commonError(err, new Map(), 'INITIALIZATION.ALERTS.POST_CFG_FAILED');
      },
    );
  }

  goToBoard() {
    if (this.config) {
      window.open('http://' + this.config.apiserver.hostname);
    } else {
      const boardURL = window.location.hostname;
      window.open('http://' + boardURL);
    }
  }

  goToAdminserver() {
    window.sessionStorage.removeItem('token');
    this.router.navigateByUrl('account/login');
  }

  onFocusBaselineHelper() {
    this.showBaselineHelper = true;
  }

  onBlurBaselineHelper() {
    this.showBaselineHelper = false;
    const year = this.newDate.getFullYear();
    const month = this.newDate.getMonth() + 1;
    const day = this.newDate.getDate();
    this.config.apiserver.imageBaselineTime = '' + year + '-' + month + '-' + day + ' 00:00:00';
  }

  onCheckInput() {
    if (this.checkInput()) {
      this.openSSH = true;
      this.uninstallConfirm = false;
      this.user.password = '';
    }
  }

  checkInput(): boolean {
    let result = true;
    for (let item of this.myInputTemplateComponents.toArray()) {
      item.checkSelf();
      if (!item.disabled && !item.isValid) {
        item.element.nativeElement.scrollIntoView();
        result = false;
        break;
      }
    }
    return result;
  }

  commonError(err: HttpErrorResponse, errorList: Map<number, string>, finnalError: string) {
    console.error(err.message);
    this.refresh = true;
    this.submitBtnState = ClrLoadingState.DEFAULT;
    if (err.status === 401) {
      this.tokenError();
      return;
    }
    errorList.forEach((msg, e) => {
      if (err.status === e) {
        this.messageService.showOnlyOkDialog(msg, 'ACCOUNT.ERROR');
        return;
      }
    });
    this.messageService.showOnlyOkDialog(finnalError, 'ACCOUNT.ERROR');
  }

  getSysStatusFailed() {
    this.messageService.showOnlyOkDialogObservable('INITIALIZATION.ALERTS.GET_SYS_STATUS_FAILED', 'ACCOUNT.ERROR').subscribe(
      (msg: Message) => {
        if (msg.returnStatus === ReturnStatus.rsConfirm) {
          location.reload();
        }
      }
    );
    this.submitBtnState = ClrLoadingState.DEFAULT;
  }

  tokenError() {
    this.messageService.showOnlyOkDialogObservable('ACCOUNT.TOKEN_ERROR', 'ACCOUNT.ERROR').subscribe(
      (msg: Message) => {
        if (msg.returnStatus === ReturnStatus.rsConfirm) {
          location.reload();
        }
      }
    );
  }

}
