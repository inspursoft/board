import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { NodeActionsType, NodeLogResponse, ResponseArrayNode } from '../resource.types';
import { CustomHttpClient } from './custom-http.service';
import { WebsocketService } from "./websocket.service";
import { map } from "rxjs/operators";

@Injectable()
export class ResourceService {

  constructor(private http: CustomHttpClient,
              private wsService: WebsocketService) {
  }

  getNodeList(): Observable<ResponseArrayNode> {
    return this.http.getArrayJson(`/v1/admin/node/list`, ResponseArrayNode);
  }

  addRemoveNode(type: NodeActionsType, nodeIp: string): Observable<NodeLogResponse> {
    const url = type === NodeActionsType.Add ?
      'ws://10.110.25.227:8080/v1/admin/node/add?node_ip=' :
      'ws://10.110.25.227:8080/v1/admin/node/delete?node_ip=';
    return this.wsService.connect(`${url}${nodeIp}`)
      .pipe(map((msg: MessageEvent) => new NodeLogResponse(JSON.parse(msg.data))));
  }
}
