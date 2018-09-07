import { Directive, Input } from '@angular/core';
import { AbstractControl, AsyncValidator, NG_ASYNC_VALIDATORS, ValidatorFn, Validators } from '@angular/forms';

import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/first';
import 'rxjs/add/operator/map'
import { HttpClient } from '@angular/common/http';
import { AppInitService } from '../../app.init.service';
import { MessageService } from '../message-service/message.service';

@Directive({
  selector: '[checkItemExisting]',
  providers: [
    {provide: NG_ASYNC_VALIDATORS, useExisting: CheckItemExistingDirective, multi: true}
  ]
})
export class CheckItemExistingDirective implements AsyncValidator {
  @Input() checkItemExisting;
  @Input() userID: number = 0;

  valFn: ValidatorFn = Validators.nullValidator;

  constructor(private http: HttpClient,
              private appInitService: AppInitService,
              private messageService: MessageService) {
  }

  checkUserExists(target: string, value: string, userID: number): Observable<{[key: string]: any}> {
    return this.http.get("/api/v1/user-exists", {
      observe: "response",
      params: {
        'target': target,
        'value': value,
        'user_id': userID.toString()
      }
    })
    .map(()=>this.valFn)
    .catch(err=>{
      this.messageService.cleanNotification();
      if(err && err.status === 409) {
        return Observable.of({ 'checkItemExisting': { value } });
      }
    });
  }
  
  checkProjectExists(token: string, projectName: string): Observable<{[key: string]: any}>{
    return this.http.head('/api/v1/projects', {
      observe: "response",
      params: {
        'project_name': projectName
      }
    })
    .map(()=>this.valFn)
    .catch(err=>{
      this.messageService.cleanNotification();
      if(err && err.status === 409) {
        return Observable.of({ 'checkItemExisting': { projectName } });
      }
    });
  } 

  validate(control: AbstractControl): Observable<{[key: string]: any}> {
    return control.valueChanges
      .debounceTime(200)
      .distinctUntilChanged()
      .switchMap(value=> {
        switch(this.checkItemExisting) {
        case 'username':
        case 'email':
          return this.checkUserExists(this.checkItemExisting, value, this.userID);
        case 'project':
          return this.checkProjectExists(this.appInitService.token, value);
        }
      })
      .first();
  }
}