import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';

import { AppInitService } from '../app.init.service';


@Injectable()
export class NodeService {
  get defaultHeader(): Headers {
    let header = new Headers();
    header.append('content-type', 'application/json');
    header.append('token', this.appInitService.token);
    return header; 
  }
  
  constructor(
    private http: Http,
    private appInitService: AppInitService
  ){}

  getNodes(): Promise<any> {
    return this.http
      .get(`/api/v1/nodes`, { headers: this.defaultHeader })
      .toPromise()
      .then(res=>{
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err=>Promise.reject(err));
  }

  getNodeByName(nodeName: string): Promise<any> {
    return this.http
      .get(`/api/v1/node`, { 
         headers: this.defaultHeader,
         params: {
           'node_name': nodeName
         }
      })
      .toPromise()
      .then(res=>{
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err=>Promise.reject(err));
  }

  toggleNodeStatus(nodeName: string, status: boolean) {
    return this.http
      .get(`/api/v1/node/toggle`, {
        headers: this.defaultHeader,
        params: {
          'node_name': nodeName,
          'node_status': status
        }
      })
      .toPromise()
      .then(res=>this.appInitService.chainResponse(res))
      .catch(err=>Promise.reject(err));
  }

}