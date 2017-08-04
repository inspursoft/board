import { Injectable } from '@angular/core';

import { Subject} from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';

import { Service } from './service';


@Injectable()
export class K8sService {

  stepSource: Subject<number>= new Subject<number>();
  step$: Observable<number> = this.stepSource.asObservable();

  getServices(): Promise<Service[]> {
    return new Promise((resolve, reject)=>resolve([
        { service_name: 'portal_bu01', service_project_name: 'bu01', service_owner: 'aron', service_create_time: new Date('2017-08-04T09:54:32+08:00'), service_public: true, service_status: 0 },
        { service_name: 'hr_bu01', service_project_name: 'bu01', service_owner: 'bron', service_create_time: new Date('2017-08-03T13:52:16+08:00'), service_public: false, service_status: 0 },
        { service_name: 'bigdata_bu02', service_project_name: 'bu02', service_owner: 'mike', service_create_time: new Date('2017-07-31T14:20:44+08:00'), service_public: false, service_status: 1 },
        { service_name: 'testenv', service_project_name: 'du01', service_owner: 'tim', service_create_time: new Date('2017-07-28T14:20:44+08:00'), service_public: true, service_status: 1 }
      ]));
  }
}