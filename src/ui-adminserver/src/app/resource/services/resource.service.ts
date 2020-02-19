import { Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { webSocket } from 'rxjs/webSocket';
import { NodeActionsType, NodeLogResponse, ResponseArrayNode } from '../resource.types';
import { CustomHttpClient } from './custom-http.service';

@Injectable()
export class ResourceService {

  constructor(private http: CustomHttpClient) {
  }

  getNodeList(): Observable<ResponseArrayNode> {
    return this.http.getArrayJson(`/v1/admin/node/list`, ResponseArrayNode);
  }

  addRemoveNode(type: NodeActionsType, nodeIp: string): Observable<NodeLogResponse> {
    const url = type === NodeActionsType.Add ?
      `ws://127.0.0.1:8080/v1/admin/node/add?node_ip=${nodeIp}` :
      `ws://127.0.0.1:8080/v1/admin/node/delete?node_ip=${nodeIp}`;
    return webSocket<NodeLogResponse>(url).asObservable();
  }
}
