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

  // TODO put some config into default.
  simpleMode = false;
  loadingFlag = true;
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
    this.accountService.createUUID().subscribe(
      () => {
        this.loadingFlag = false;
        this.enableInitialization = true;
      },
      (err: HttpErrorResponse) => {
        console.error(err);
        this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.INITIALIZATION', 'ACCOUNT.ERROR');
      });
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
    // this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.ALREADY_START', 'GLOBAL_ALERT.HINT');

    this.uuidInput.checkSelf();
    if (this.uuidInput.isValid) {
      this.submitBtnState = ClrLoadingState.LOADING;
      this.accountService.validateUUID(this.uuid).subscribe(
        () => {
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
                      console.error(err);
                      this.submitBtnState = ClrLoadingState.DEFAULT;
                      this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.GET_CFG_FAILED', 'ACCOUNT.ERROR');
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
                  this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.ALREADY_START', 'GLOBAL_ALERT.HINT');
                  this.submitBtnState = ClrLoadingState.DEFAULT;
                  break;
                }
              }
            },
            (err: HttpErrorResponse) => {
              console.error(err);
              this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.GET_SYS_STATUS_FAILED', 'ACCOUNT.ERROR');
              this.submitBtnState = ClrLoadingState.DEFAULT;
            }
          );
        },
        (err: HttpErrorResponse) => {
          console.error(err);
          this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.VALIDATE_UUID_FAILED', 'ACCOUNT.ERROR');
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
              console.log('Can not read tmp file: ' + err.message);
              this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.GET_TMP_FAILED', 'GLOBAL_ALERT.HINT');
              this.newDate = new Date(this.config.apiserver.imageBaselineTime);
              this.isEditable = this.config.isInit;
              this.installStep++;
              this.installProgress += 33;
              this.submitBtnState = ClrLoadingState.DEFAULT;
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
        console.error(err);
        this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.GET_CFG_FAILED', 'ACCOUNT.ERROR');
        this.submitBtnState = ClrLoadingState.DEFAULT;
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
        console.error(err);
        this.submitBtnState = ClrLoadingState.DEFAULT;
        this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.START_BOARD_FAILED', 'ACCOUNT.ERROR');
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
        console.error(err);
        this.submitBtnState = ClrLoadingState.DEFAULT;
        this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.UNINSTALL_BOARD_FAILED', 'ACCOUNT.ERROR');
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
    this.configurationService.postConfig(this.config).subscribe(
      () => {
        this.boardService.applyCfg(this.user).subscribe(
          () => {
            this.installStep++;
            this.installProgress = 100;
            this.submitBtnState = ClrLoadingState.DEFAULT;
          },
          (err: HttpErrorResponse) => {
            console.error(err);
            this.submitBtnState = ClrLoadingState.DEFAULT;
            this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.START_BOARD_FAILED', 'ACCOUNT.ERROR');
          },
        );
      },
      (err: HttpErrorResponse) => {
        console.error(err);
        this.submitBtnState = ClrLoadingState.DEFAULT;
        this.messageService.showOnlyOkDialog('INITIALIZATION.ALERTS.POST_CFG_FAILED', 'ACCOUNT.ERROR');
      },
    );
  }

  goToBoard() {
    if (this.config) {
      window.open(this.config.apiserver.hostname);
    } else {
      const boardURL = window.location.hostname;
      window.open(boardURL + ':80');
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

}
