import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { NodeDetails, NodeList, NodeLog, NodeLogs, NodePostData, NodePreparationData } from '../resource.types';
import { CustomHttpClient } from './custom-http.service';
import { map } from 'rxjs/operators';

@Injectable()
export class ResourceService {

  constructor(private http: CustomHttpClient) {
  }

  getNodeList(): Observable<NodeList> {
    return this.http.getArrayJson(`/v1/admin/node`, NodeList);
  }

  addNode(postData: NodePostData): Observable<NodeLog> {
    return this.http.postJson('/v1/admin/node', postData.getPostData(), NodeLog);
  }

  removeNode(paramsData: NodePostData): Observable<NodeLog> {
    return this.http.delete(`/v1/admin/node?`, {params: paramsData.getParamsData()})
      .pipe(map((res: object) => new NodeLog(res)));
  }

  getNodeLogs(pageIndex, pageSize: number): Observable<NodeLogs> {
    return this.http.getPagination('/v1/admin/node/logs', NodeLogs,
      {page_index: pageIndex.toString(), page_size: pageSize.toString()});
  }

  getNodePreparation(): Observable<NodePreparationData> {
    return this.http.getJson('/v1/admin/node/preparation', NodePreparationData);
  }

  getNodeLogDetail(ip: string, creationTime: number): Observable<NodeDetails> {
    return this.http.getArrayJson(`/v1/admin/node/log?`, NodeDetails, {
      node_ip: ip,
      creation_time: creationTime.toString()
    });
  }
}
