import { Directive, Input } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { AppInitService } from '../../shared.service/app-init.service';
import { ValidationErrors } from "@angular/forms/src/directives/validators";
import { UsernameInUseKey } from "../shared.const";
import { Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";
import { MessageService } from "../../shared.service/message.service";
import { InputExComponent } from "board-components-library";
import { TranslateService } from "@ngx-translate/core";

@Directive({
  selector: '[appLibCheckItemExistingEx]'
})
export class LibCheckExistingExDirective {
  @Input() appLibCheckItemExistingEx = '';
  @Input() userID = 0;
  usernameIsKey = '';
  usernameExists = '';
  emailExists = '';
  projectNameExists = '';

  constructor(private inputComponent: InputExComponent,
              private http: HttpClient,
              private appInitService: AppInitService,
              private messageService: MessageService,
              private translateService: TranslateService) {
    this.inputComponent.validatorAsyncFn = this.validateAction.bind(this);
    this.translateService.get('ACCOUNT.USERNAME_IS_KEY').subscribe(res => this.usernameIsKey = res);
    this.translateService.get('PROJECT.PROJECT_NAME_ALREADY_EXISTS').subscribe(res => this.projectNameExists = res);
    this.translateService.get('ACCOUNT.USERNAME_ALREADY_EXISTS').subscribe(res => this.usernameExists = res);
    this.translateService.get('ACCOUNT.EMAIL_ALREADY_EXISTS').subscribe(res => this.emailExists = res);
  }

  checkUserExists(value: string, errorMsg: string): Observable<ValidationErrors | null> {
    if (this.appLibCheckItemExistingEx === 'username' && UsernameInUseKey.indexOf(value) > 0) {
      return of({checkItemExistingEx: this.usernameIsKey});
    }
    return this.http.get('/api/v1/user-exists', {
      observe: 'response',
      params: {
        target: this.appLibCheckItemExistingEx,
        value,
        user_id: this.userID.toString()
      }
    }).pipe(
      map(() => null),
      catchError(err => {
        this.messageService.cleanNotification();
        if (err && err.status === 409) {
          return of({checkItemExistingEx: errorMsg});
        }
        return null;
      }));
  }

  checkProjectExists(projectName: string): Observable<ValidationErrors | null> {
    return this.http.head('/api/v1/projects', {
      observe: 'response',
      params: {
        project_name: projectName
      }
    }).pipe(
      map(() => null),
      catchError(err => {
        this.messageService.cleanNotification();
        if (err && err.status === 409) {
          return of({checkItemExistingEx: this.projectNameExists});
        }
        return null;
      }));
  }

  validateAction(control: AbstractControl): Observable<ValidationErrors | null> {
    switch (this.appLibCheckItemExistingEx) {
      case 'username':
        return this.checkUserExists(control.value, this.usernameExists);
      case 'email':
        return this.checkUserExists(control.value, this.emailExists);
      case 'project':
        return this.checkProjectExists(control.value);
      default:
        return of(null);
    }
  }
}
