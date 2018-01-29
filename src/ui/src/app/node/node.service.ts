import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable()
export class NodeService {
  constructor(private http: HttpClient) {
  }

  getNodes(): Promise<any> {
    return this.http
      .get(`/api/v1/nodes`, {observe: "response"})
      .toPromise()
      .then(res => res.body)
  }

  getNodeByName(nodeName: string): Promise<any> {
    return this.http
      .get(`/api/v1/node`, {
        observe: "response",
        params: {
          'node_name': nodeName
        }
      })
      .toPromise()
      .then(res => res.body)
  }

  toggleNodeStatus(nodeName: string, status: boolean): Promise<any> {
    return this.http
      .get(`/api/v1/node/toggle`, {
        observe: "response",
        params: {
          'node_name': nodeName,
          'node_status': status ? "1" : "0"
        }
      })
      .toPromise()
  }

}