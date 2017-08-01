import { Directive, Input} from '@angular/core';
import { NG_ASYNC_VALIDATORS, AsyncValidator, 
         Validators, ValidatorFn,
         AbstractControl } from '@angular/forms';

import { Observable } from 'rxjs/Observable';

import 'rxjs/add/operator/debounceTime';
import 'rxjs/add/operator/distinctUntilChanged';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/first';

import { Http } from '@angular/http';

import { AppInitService } from '../../app.init.service';

@Directive({
  selector: '[checkItemExisting]',
  providers: [
    {provide: NG_ASYNC_VALIDATORS, useExisting: CheckItemExistingDirective, multi: true }
  ]
})
export class CheckItemExistingDirective implements AsyncValidator {
  @Input() checkItemExisting;
    
  valFn: ValidatorFn = Validators.nullValidator;

  constructor(
    private http: Http,
    private appInitService: AppInitService
  ){}

  checkUserExists(target: string, value: string): Observable<{[key: string]: any}> {
    return this.http.get("/api/v1/user-exists", {
      params: {
        'target': target,
        'value': value
      }
    })
    .map(()=>this.valFn)
    .catch(()=>Observable.of({ 'checkItemExisting': {value} }));
  }
  
  checkProjectExists(token: string, projectName: string): Observable<{[key: string]: any}>{
    return this.http.head('/api/v1/projects', {
      params: {
        'token': token,
        'project_name': projectName
      }
    })
    .map(()=>this.valFn)
    .catch(()=>Observable.of({ 'checkItemExisting': projectName }));
  } 

  validate(control: AbstractControl): Observable<{[key: string]: any}> {
    return control.valueChanges
      .debounceTime(200)
      .distinctUntilChanged()
      .switchMap(value=> {
        switch(this.checkItemExisting) {
        case 'username':
        case 'email':
          return this.checkUserExists(this.checkItemExisting, value);
        case 'project':
          return this.checkProjectExists(this.appInitService.token, value);
        }
      })
      .first();
  }
}