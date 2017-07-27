import { Directive, Input} from '@angular/core';
import { NG_ASYNC_VALIDATORS, AsyncValidator, 
         Validators, ValidatorFn,
         AbstractControl } from '@angular/forms';

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

  checkUserExists(target: string, value: string): Promise<any> {
    return this.http.get("/api/v1/user-exists", {
      params: {
        'target': target,
        'value': value
      }
    }).toPromise();
  }
  
  checkProjectExists(token: string, projectName: string): Promise<any>{
    return this.http.head('/api/v1/projects', {
      params: {
        'token': token,
        'project_name': projectName
      }
    }).toPromise();
  } 

  validate(control: AbstractControl): Promise<{[key: string]: any}> {
    const value = control.value; 
    switch(this.checkItemExisting) {
    case 'username':
    case 'email':
      return this.checkUserExists(this.checkItemExisting, value)
        .then(()=>this.valFn)
        .catch(()=>Promise.resolve({ 'checkItemExisting': value }));
    case 'project':
      return this.checkProjectExists(this.appInitService.token, value)
        .then(()=>this.valFn)
        .catch(()=>Promise.resolve({ 'checkItemExisting': value }));
    default:
      return Promise.resolve(this.valFn);
    }
  }
}