import { Component, OnInit, ViewChild, ViewChildren, QueryList } from '@angular/core';
import { AccountService } from '../account.service';
import { HttpErrorResponse } from '@angular/common/http';
import { MessageService } from 'src/app/shared/message/message.service';
import { VariableInputComponent } from 'src/app/shared/variable-input/variable-input.component';
import { ClrLoadingState } from '@clr/angular';
import { AppInitService } from 'src/app/shared.service/app-init.service';
import { InitStatus, InitStatusCode } from 'src/app/shared.service/app-init.model';
import { User } from '../account.model';
import { BoardService } from 'src/app/shared.service/board.service';
import { ConfigurationService } from 'src/app/shared.service/configuration.service';
import { Configuration } from 'src/app/shared.service/cfg.model';
import { Router } from '@angular/router';
import { Message, ReturnStatus } from 'src/app/shared/message/message.types';

@Component({
  selector: 'app-installation',
  templateUrl: './installation.component.html',
  styleUrls: ['./installation.component.css']
})
export class InstallationComponent implements OnInit {
  debugMode = false;
  status = 'status1';

  newDate = new Date('2016-01-01 09:00:00');

  passwordPattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)([A-Za-z\d#?!@$%^&*-]){8,20}$/;
  uuidPattern = /^\w{8}(-\w{4}){3}-\w{12}$/;

  uuid = '';
  config: Configuration;
  showBaselineHelper = false;

  installStep = 0;
  ignoreStep1 = false;
  ignoreStep2 = false;
  installProgress = 0;
  enableBtn = true;
  refresh = false;

  simpleMode = false;
  loadingFlag = true;
  disconnect = false;
  enableInitialization = false;
  openSSH = false;
  uninstallConfirm = false;
  clearDate = false;
  responsibility = false;
  isEditable = true;
  submitBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  editBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  startBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  uninstallBtnState: ClrLoadingState = ClrLoadingState.DEFAULT;
  isEdit = false;
  isStart = false;
  isUninstall = false;
  user: User;
  startLog = '';
  modalSize = '';
  isKeepUninstall = false;
  disableUninstall = false;

  @ViewChild('UUID') uuidInput: VariableInputComponent;
  @ViewChildren(VariableInputComponent) myInputTemplateComponents: QueryList<VariableInputComponent>;

  constructor(private accountService: AccountService,
              private appInitService: AppInitService,
              private boardService: BoardService,
              private configurationService: ConfigurationService,
              private messageService: MessageService,
              private router: Router) {
    this.user = new User();
  }

  ngOnInit() {
    if (this.debugMode) {
      setTimeout(() => {
        this.loadingFlag = false;
        this.enableInitialization = true;
      }, 100);
      return;
    }

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
    if (this.debugMode) {
      if ('status1' === this.status) {
        this.config = new Configuration();
        setTimeout(() => {
          this.ignoreStep1 = true;
          this.installStep = 2;
          this.installProgress = 50;
        }, 1000);
      } else if ('status2' === this.status) {
        this.isKeepUninstall = true;
        this.installStep++;
        this.installProgress += 33;
      } else if ('status3' === this.status) {
        this.ignoreStep1 = true;
        this.ignoreStep2 = true;
        this.installStep = 3;
        this.installProgress = 100;
        this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.ALREADY_START');
      } else {
        alert('error status!');
      }
      return;
    }

    this.uuidInput.checkSelf();
    if (this.uuidInput.isValid) {
      this.submitBtnState = ClrLoadingState.LOADING;
      this.user.username = 'boardadmin';
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
                      this.commonError(err, {}, 'INITIALIZATION.ALERTS.GET_CFG_FAILED');
                    }
                  );
                  break;
                }
                // 未起Board但更改过cfg
                case InitStatusCode.InitStatusSecond: {
                  this.isKeepUninstall = true;
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
    if (this.debugMode) {
      this.editBtnState = ClrLoadingState.LOADING;
      this.disableUninstall = true;
      this.isEdit = true;
      setTimeout(() => {
        this.installStep++;
        this.installProgress += 33;
        this.config = new Configuration();
        this.editBtnState = ClrLoadingState.DEFAULT;
        this.disableUninstall = false;
      }, 5000);
      return;
    }

    this.editBtnState = ClrLoadingState.LOADING;
    this.disableUninstall = true;
    this.isEdit = true;
    this.configurationService.getConfig().subscribe(
      (res: Configuration) => {
        this.config = new Configuration(res);
        if (this.config.tmpExist) {
          this.configurationService.getConfig('tmp').subscribe(
            (resTmp: Configuration) => {
              this.config = new Configuration(resTmp);
              this.newDate = new Date(this.config.k8s.imageBaselineTime);
              this.isEditable = this.config.isInit;
              this.installStep++;
              this.installProgress += 33;
              this.editBtnState = ClrLoadingState.DEFAULT;
            },
            (err: HttpErrorResponse) => {
              if (err.status === 401) {
                this.tokenError();
              } else {
                console.log('Can not read tmp file: ' + err.message);
                this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.GET_TMP_FAILED');
                this.newDate = new Date(this.config.k8s.imageBaselineTime);
                this.isEditable = this.config.isInit;
                this.installStep++;
                this.installProgress += 33;
                this.editBtnState = ClrLoadingState.DEFAULT;
              }
            }
          );
        } else {
          this.newDate = new Date(this.config.k8s.imageBaselineTime);
          this.isEditable = this.config.isInit;
          this.installStep++;
          this.installProgress += 33;
          this.editBtnState = ClrLoadingState.DEFAULT;
        }
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.isEdit = false;
        this.commonError(err, {}, 'INITIALIZATION.ALERTS.GET_CFG_FAILED');
      },
      () => {
        this.disableUninstall = false;
      }
    );
  }

  onStartBoard() {
    // for test
    if (this.debugMode) {
      this.startBtnState = ClrLoadingState.LOADING;
      this.submitBtnState = ClrLoadingState.LOADING;
      this.isStart = true;
      setTimeout(() => {
        this.openSSH = false;
        this.installStep += 2;
        this.ignoreStep2 = true;
        this.installProgress = 100;
        this.startBtnState = ClrLoadingState.DEFAULT;
        this.submitBtnState = ClrLoadingState.DEFAULT;
      }, 20000);
      return;
    }

    this.startBtnState = ClrLoadingState.LOADING;
    this.submitBtnState = ClrLoadingState.LOADING;
    this.isStart = true;
    this.boardService.start(this.user).subscribe(
      () => {
        this.openSSH = false;
        this.installStep += 2;
        this.ignoreStep2 = true;
        this.installProgress = 100;
        this.startBtnState = ClrLoadingState.DEFAULT;
        this.submitBtnState = ClrLoadingState.DEFAULT;
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.openSSH = false;
        this.isStart = false;
        this.commonError(err, {}, 'INITIALIZATION.ALERTS.START_BOARD_FAILED');
      },
    );
  }

  onUninstallBoard() {
    // for test
    if (this.debugMode) {
      this.uninstallBtnState = ClrLoadingState.LOADING;
      this.submitBtnState = ClrLoadingState.LOADING;
      this.isUninstall = true;
      setTimeout(() => {
        this.openSSH = false;
        this.installStep = 4;
        this.ignoreStep1 = true;
        this.ignoreStep2 = true;
        this.installProgress = 100;
        this.uninstallBtnState = ClrLoadingState.DEFAULT;
        this.submitBtnState = ClrLoadingState.DEFAULT;
      }, 20000);
      return;
    }

    this.uninstallBtnState = ClrLoadingState.LOADING;
    this.submitBtnState = ClrLoadingState.LOADING;
    this.isUninstall = true;
    this.boardService.shutdown(this.user, this.clearDate).subscribe(
      () => {
        this.openSSH = false;
        this.installStep = 4;
        this.ignoreStep1 = true;
        this.ignoreStep2 = true;
        this.installProgress = 100;
        this.uninstallBtnState = ClrLoadingState.DEFAULT;
        this.submitBtnState = ClrLoadingState.DEFAULT;
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.openSSH = false;
        this.isUninstall = false;
        this.commonError(err, { 503: 'INITIALIZATION.ALERTS.ALREADY_UNINSTALL' }, 'INITIALIZATION.ALERTS.UNINSTALL_BOARD_FAILED');
      },
    );
  }

  onApplyAndStartBoard() {
    // for test
    if (this.debugMode) {
      this.submitBtnState = ClrLoadingState.LOADING;
      let helloTime = 0;
      const outPutLog = setInterval(() => {
        this.startLog += `hello world! ${helloTime++}\n`;
        if (!this.modalSize) {
          this.modalSize = 'xl';
        }
      }, 1000);
      setTimeout(() => {
        this.openSSH = false;
        this.installStep++;
        this.installProgress = 100;
        this.submitBtnState = ClrLoadingState.DEFAULT;
        this.modalSize = '';
        clearInterval(outPutLog);
      }, 60 * 1000);
      return;
    }


    this.submitBtnState = ClrLoadingState.LOADING;
    let maxTry = 100;
    const installStepNow = this.installStep;
    this.configurationService.putConfig(this.config).subscribe(
      () => {
        this.boardService.applyCfg(this.user).subscribe(
          () => {
            const initProcess = setInterval(() => {
              this.appInitService.getSystemStatus().subscribe(
                (res: InitStatus) => {
                  if (InitStatusCode.InitStatusThird === res.status) {
                    if (this.installStep === installStepNow) {
                      this.installStep++;
                    }
                    this.installProgress = 100;
                    this.openSSH = false;
                    this.submitBtnState = ClrLoadingState.DEFAULT;
                    this.modalSize = '';
                    this.startLog = '';
                    clearInterval(initProcess);
                  } else {
                    this.startLog = res.log;
                    if (!this.modalSize) {
                      this.modalSize = 'xl';
                    }
                  }
                },
                (err: HttpErrorResponse) => {
                  console.error(err.message);
                  if (maxTry-- < 0) {
                    this.getSysStatusFailed();
                    clearInterval(initProcess);
                  }
                },
              );
            }, 5 * 1000);
          },
          (err: HttpErrorResponse) => {
            // COMMON
            this.commonError(err, {}, 'INITIALIZATION.ALERTS.START_BOARD_FAILED');
          },
        );
      },
      (err: HttpErrorResponse) => {
        // COMMON
        this.commonError(err, {}, 'INITIALIZATION.ALERTS.POST_CFG_FAILED');
      },
    );
  }

  goToBoard() {
    if (this.config) {
      const protocol = this.config.board.mode === 'normal' ? 'http' : 'https';
      window.open(`${protocol}://${this.config.board.hostname}`);
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
    this.config.k8s.imageBaselineTime = '' + year + '-' + month + '-' + day + ' 00:00:00';
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
    for (const item of this.myInputTemplateComponents.toArray()) {
      item.checkSelf();
      if (!item.disabled && !item.isValid) {
        item.element.nativeElement.scrollIntoView();
        result = false;
        break;
      }
    }
    return result;
  }

  commonError(err: HttpErrorResponse, errorList: object, finnalError: string) {
    console.error(err.message);
    this.refresh = true;
    this.submitBtnState = ClrLoadingState.DEFAULT;
    this.editBtnState = ClrLoadingState.DEFAULT;
    this.startBtnState = ClrLoadingState.DEFAULT;
    this.uninstallBtnState = ClrLoadingState.DEFAULT;
    if (err.status === 401) {
      this.tokenError();
      return;
    }
    for (const key in errorList) {
      if (err.status === Number(key)) {
        this.messageService.showOnlyOkDialog(String(errorList[key]), 'ACCOUNT.ERROR');
        return;
      }
    }
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
    this.messageService.showOnlyOkDialogObservable('ACCOUNT.TOKEN_ERROR_TO_REFRESH', 'ACCOUNT.ERROR').subscribe(
      (msg: Message) => {
        if (msg.returnStatus === ReturnStatus.rsConfirm) {
          location.reload();
        }
      }
    );
  }

}
