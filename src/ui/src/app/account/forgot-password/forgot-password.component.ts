import { Component, OnInit } from '@angular/core';
import { MessageService } from '../../shared.service/message.service';
import { ActivatedRoute, Router } from '@angular/router';
import { AccountService } from '../account.service';
import { RouteSignIn } from '../../shared/shared.const';
import { AppInitService } from '../../shared.service/app-init.service';
import { HttpErrorResponse } from '@angular/common/http';
import { CsComponentBase } from '../../shared/cs-components-library/cs-component-base';

@Component({
  selector: 'forgot-password',
  templateUrl: './forgot-password.component.html',
  styleUrls: ['./forgot-password.component.css']
})
export class ForgotPasswordComponent extends CsComponentBase implements OnInit {
  credential = '';
  sendRequestWIP = false;

  constructor(private accountService: AccountService,
              private messageService: MessageService,
              private activatedRoute: ActivatedRoute,
              private appInitService: AppInitService,
              private router: Router) {
    super();
  }

  ngOnInit() {
    if (this.appInitService.systemInfo.auth_mode !== 'db_auth') {
      this.router.navigate([RouteSignIn]).then();
    }
  }

  goBack(): void {
    this.router.navigate([RouteSignIn]).then();
  }

  sendRequest(): void {
    if (this.verifyInputExValid()) {
      this.sendRequestWIP = true;
      this.accountService.postEmail(this.credential).subscribe(
        () => this.messageService.showOnlyOkDialogObservable('ACCOUNT.SEND_REQUEST_SUCCESS_MSG', 'ACCOUNT.SEND_REQUEST_SUCCESS')
          .subscribe(() => this.router.navigate([RouteSignIn]).then()),
        (err: HttpErrorResponse) => {
          this.sendRequestWIP = false;
          if (err.status === 404) {
            this.messageService.showOnlyOkDialog('ACCOUNT.USER_NOT_EXISTS', 'ACCOUNT.SEND_REQUEST_ERR');
          } else {
            this.messageService.showOnlyOkDialog('ACCOUNT.SEND_REQUEST_ERR_MSG', 'ACCOUNT.SEND_REQUEST_ERR');
          }
        });
    }
  }
}
