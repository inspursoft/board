import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { NodeDetails, NodeList, NodeLog, NodeLogs } from '../resource.types';
import { CustomHttpClient } from './custom-http.service';
import { map } from "rxjs/operators";

@Injectable()
export class ResourceService {

  constructor(private http: CustomHttpClient) {
  }

  getNodeList(): Observable<NodeList> {
    return this.http.getArrayJson(`/v1/admin/node/list`, NodeList);
  }

  addNode(nodeIp: string): Observable<NodeLog> {
    return this.http.postJson('/v1/admin/node/add', {node_ip: nodeIp}, NodeLog);
  }

  removeNode(nodeIp: string): Observable<NodeLog> {
    return this.http.delete(`/v1/admin/node/remove?node_ip=${nodeIp}`)
      .pipe(map((res: object) => new NodeLog(res)));
  }

  getNodeLogs(): Observable<NodeLogs> {
    return this.http.getArrayJson('/v1/admin/node/logs', NodeLogs);
  }

  getNodeLog(logFileName: string): Observable<NodeDetails> {
    return this.http.getArrayJson(`/v1/admin/node/log?file_name=${logFileName}`, NodeDetails);
  }
}
