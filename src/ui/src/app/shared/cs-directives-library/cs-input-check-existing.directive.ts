import { Directive, Input } from '@angular/core';
import { AbstractControl, Validators } from '@angular/forms';
import { Observable } from 'rxjs/Observable';
import { HttpClient } from '@angular/common/http';
import { AppInitService } from '../../app.init.service';
import { MessageService } from '../message-service/message.service';
import { CsInputComponent } from "../cs-components-library/cs-input/cs-input.component";
import { ValidationErrors } from "@angular/forms/src/directives/validators";
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map'
import 'rxjs/add/observable/of'

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
    this.csInputComponent.inputControl.setAsyncValidators(this.validateAction.bind(this));
  }

  checkUserExists(value: string, errorMsg: string): Observable<ValidationErrors | null> {
    return this.http.get("/api/v1/user-exists", {
      observe: "response",
      params: {
        'target': this.checkItemExistingEx,
        'value': value,
        'user_id': this.userID.toString()
      }
    }).map(() => Validators.nullValidator)
      .catch(err => {
        this.messageService.cleanNotification();
        if (err && err.status === 409) {
          return Observable.of({'checkItemExistingEx': errorMsg});
        }
      });
  }

  checkProjectExists(projectName: string): Observable<ValidationErrors | null> {
    return this.http.head('/api/v1/projects', {
      observe: "response",
      params: {
        'project_name': projectName
      }
    }).map(() => Validators.nullValidator)
      .catch(err => {
        this.messageService.cleanNotification();
        if (err && err.status === 409) {
          return Observable.of({'checkItemExistingEx': "PROJECT.PROJECT_NAME_ALREADY_EXISTS"});
        }
      });
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
        return Observable.of(null);
    }
  }
}