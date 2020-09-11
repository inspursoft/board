import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { delay, map, timeout } from 'rxjs/operators';
import { HttpHeaders, HttpResponse } from '@angular/common/http';
import { NodeDetail, NodeControlStatus, NodeGroupStatus, NodeStatus, EdgeNode, NodeGroupDetail } from './node.types';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

@Injectable()
export class NodeService {
  constructor(private http: ModelHttpClient) {
  }

  getNodes(): Observable<Array<NodeStatus>> {
    return this.http.getArray(`/api/v1/nodes`, NodeStatus)
      .pipe(
        map((nodeStatusList: Array<NodeStatus>) => nodeStatusList.filter(value => value.nodeName !== ''))
      );
  }

  getNodeDetailByName(nodeName: string): Observable<NodeDetail> {
    return this.http.getJson(`/api/v1/node`, NodeDetail, {param: {node_name: nodeName}});
  }

  removeEdgeNode(nodeName: string): Observable<any> {
    return this.http.delete(`/api/v1/edgenodes/${nodeName}`);
  }

  toggleNodeStatus(nodeName: string, status: boolean): Observable<HttpResponse<object>> {
    return this.http
      .get(`/api/v1/node/toggle`, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {
          node_name: nodeName,
          node_status: status ? '1' : '0'
        }
      });
  }

  getNodeGroupsOfOneNode(nodeName: string): Observable<Array<string>> {
    return this.http.get<Array<string>>(`/api/v1/node/0/group`,
      {observe: 'response', params: {node_name: nodeName}})
      .pipe(map((res: HttpResponse<Array<string>>) => res.body || []));
  }

  addNodeToNodeGroup(nodeName: string, nodeGroupName: string): Observable<object> {
    return this.http.post<object>(`/api/v1/node/0/group`, null,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {node_name: nodeName, groupname: nodeGroupName}
      }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  deleteNodeToNodeGroup(nodeName: string, nodeGroupName: string): Observable<object> {
    return this.http.delete<object>(`/api/v1/node/0/group`,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {node_name: nodeName, groupname: nodeGroupName}
      })
      .pipe(map((res: HttpResponse<object>) => res.body));
  }

  getNodeGroups(): Observable<Array<NodeGroupStatus>> {
    return this.http.getArray(`/api/v1/nodegroup`, NodeGroupStatus);
  }

  addNodeGroup(group: NodeGroupStatus): Observable<HttpResponse<object>> {
    return this.http.post(`/api/v1/nodegroup`, group.postBody(), {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  addEdgeNode(edgeNode: EdgeNode): Observable<any> {
    return this.http.post(`/api/v1/edgenodes`, edgeNode.getPostBody()).pipe(timeout(20000));
  }

  deleteNodeGroup(groupId: number, nodeGroupName: string): Observable<HttpResponse<object>> {
    return this.http.delete(`/api/v1/nodegroup/${groupId}`,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {groupname: nodeGroupName}
      });
  }

  checkNodeGroupExist(groupName: string): Observable<HttpResponse<object>> {
    return this.http.get(`/api/v1/nodegroup/existing`,
      {observe: 'response', params: {nodegroup_name: groupName}});
  }

  getNodeControlStatus(nodeName: string): Observable<NodeControlStatus> {
    return this.http.getJson(`/api/v1/nodes/${nodeName}`, NodeControlStatus);
  }

  drainNodeService(nodeName: string, serviceInstanceCount: number): Observable<any> {
    return this.http.put(`/api/v1/nodes/${nodeName}/drain`, null)
      .pipe(delay(500 * serviceInstanceCount));
  }

  getNodeName(nodeIp, nodePassword: string): Observable<string> {
    return this.http.get(`/api/v1/edgenodes/checkedgename`, {
      responseType: 'text',
      params: {
        edge_ip: nodeIp,
        edge_password: nodePassword
      }
    });
  }

  getGroupMembers(groupId: number): Observable<Array<string>> {
    return this.http.getJson(`/api/v1/nodegroup/${groupId}`, NodeGroupDetail)
      .pipe(map((res: NodeGroupDetail) => res.nodeList));
  }

  updateGroup(nodeGroup: NodeGroupStatus): Observable<any> {
    return this.http.put(`/api/v1/nodegroup/${nodeGroup.id}`, nodeGroup.postBody(),
      {
        params: {
          id: nodeGroup.id.toString()
        }
      }
    );
  }
}
