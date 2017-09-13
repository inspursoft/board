import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Service } from './service';
import { AppInitService } from "../app.init.service";
import { Http, Headers, RequestOptions } from "@angular/http";
import { Project } from "../project/project";
import { Image, ImageDetail } from "../image/image";
import { ServiceStep2Output } from "./service-step.component";

@Injectable()
export class K8sService {
  stepSource: Subject<number> = new Subject<number>();
  step$: Observable<number> = this.stepSource.asObservable();
  stepData: Map<number, Object>;

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  constructor(private http: Http,
              private appInitService: AppInitService) {
    this.stepData = new Map<number, Object>();
  }

  clearStepData() {
    this.stepData.clear();
  }

  setStepData(step: number, Data: Object) {
    this.stepData.set(step, Data);
  }

  buildImage(imageData: ServiceStep2Output): Promise<boolean> {
    return this.http.post(`/api/v1/images/building`, imageData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  getDockerFilePreview(imageData: ServiceStep2Output): Promise<string> {
    return this.http.post(`/api/v1/images/preview`, imageData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.text();
      })
      .catch(err => Promise.reject(err));
  }

  getFileList(formData: FormData): Promise<Array<{path: string, file_name: string, size: number}>> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/list`, formData, options).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err => Promise.reject(err));
  }

  uploadFile(formData: FormData): Promise<boolean> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/upload`, formData, options).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  getStepData(step: number): Object {
    return this.stepData.get(step);
  }

  getProjects(projectName?: string): Promise<Project[]> {
    return this.http.get('/api/v1/projects', {
      headers: this.defaultHeader,
      params: {'project_name': projectName}
    }).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err => Promise.reject(err));
  }

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Promise<Image[]> {
    let options = new RequestOptions({
      headers: this.defaultHeader,
      params: {
        'image_name': image_name,
        'image_list_page': image_list_page.toString(),
        'image_list_page_size': image_list_page_size.toString()
      }
    });
    return this.http.get("/api/v1/images", options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getImageDetailList(image_name: string): Promise<ImageDetail[]> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.get(`/api/v1/images/${image_name}`, options)
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getServices(): Promise<Service[]> {
    return new Promise((resolve, reject) => resolve([
      {
        service_name: 'portal_bu01',
        service_project_name: 'bu01',
        service_owner: 'aron',
        service_create_time: new Date('2017-08-04T09:54:32+08:00'),
        service_public: true,
        service_status: 0
      },
      {
        service_name: 'hr_bu01',
        service_project_name: 'bu01',
        service_owner: 'bron',
        service_create_time: new Date('2017-08-03T13:52:16+08:00'),
        service_public: false,
        service_status: 0
      },
      {
        service_name: 'bigdata_bu02',
        service_project_name: 'bu02',
        service_owner: 'mike',
        service_create_time: new Date('2017-07-31T14:20:44+08:00'),
        service_public: false,
        service_status: 1
      },
      {
        service_name: 'testenv',
        service_project_name: 'du01',
        service_owner: 'tim',
        service_create_time: new Date('2017-07-28T14:20:44+08:00'),
        service_public: true,
        service_status: 1
      }
    ]));
  }
}