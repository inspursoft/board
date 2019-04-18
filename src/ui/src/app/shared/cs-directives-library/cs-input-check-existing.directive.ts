import { Directive, Input } from '@angular/core';
import { AbstractControl } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { AppInitService } from '../../shared.service/app-init.service';
import { CsInputComponent } from "../cs-components-library/cs-input/cs-input.component";
import { ValidationErrors } from "@angular/forms/src/directives/validators";
import { UsernameInUseKey } from "../shared.const";
import { Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";
import { MessageService } from "../../shared.service/message.service";

@Directive({
  selector: "[checkItemExistingEx]"
})
export class CsInputCheckExistingDirective {
  @Input() checkItemExistingEx = "";
  @Input() userID: number = 0;

  constructor(private csInputComponent: CsInputComponent,
              private http: HttpClient,
              private appInitService: AppInitService,
              private messageService: MessageService) {
    this.csInputComponent.customerValidatorAsyncFunc = this.validateAction.bind(this);
  }

  checkUserExists(value: string, errorMsg: string): Observable<ValidationErrors | null> {
    if (this.checkItemExistingEx === 'username' && UsernameInUseKey.indexOf(value) > 0) {
      return of({'checkItemExistingEx': 'ACCOUNT.USERNAME_IS_KEY'})
    }
    return this.http.get("/api/v1/user-exists", {
      observe: "response",
      params: {
        'target': this.checkItemExistingEx,
        'value': value,
        'user_id': this.userID.toString()
      }
    }).pipe(
      map(() => null),
      catchError(err => {
        this.messageService.cleanNotification();
        if (err && err.status === 409) {
          return of({'checkItemExistingEx': errorMsg});
        }
        return null;
      }));
  }

  checkProjectExists(projectName: string): Observable<ValidationErrors | null> {
    return this.http.head('/api/v1/projects', {
      observe: "response",
      params: {
        'project_name': projectName
      }
    }).pipe(
      map(() => null),
      catchError(err => {
        this.messageService.cleanNotification();
        if (err && err.status === 409) {
          return of({'checkItemExistingEx': "PROJECT.PROJECT_NAME_ALREADY_EXISTS"});
        }
        return null;
      }))
  }

  validateAction(control: AbstractControl): Observable<ValidationErrors | null> {
    switch (this.checkItemExistingEx) {
      case 'username':
        return this.checkUserExists(control.value, "ACCOUNT.USERNAME_ALREADY_EXISTS");
      case 'email':
        return this.checkUserExists(control.value, "ACCOUNT.EMAIL_ALREADY_EXISTS");
      case 'project':
        return this.checkProjectExists(control.value);
      default:
        return of(null);
    }
  }
}
