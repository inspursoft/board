import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import { Observable } from "rxjs/Observable";
import "rxjs/add/operator/map"
import "rxjs/add/operator/repeat"
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from "../shared/shared.const";
import { INode, INodeGroup } from "../shared/shared.types";

export interface INodeDetail {
  node_name: string,
  node_ip: string,
  create_time: number,
  cpu_usage: number,
  memory_usage: number,
  memory_size: string,
  storage_total: string,
  storage_use: string
}

@Injectable()
export class NodeService {
  constructor(private http: HttpClient) {
  }

  getNodes(): Observable<Array<INode>> {
    return this.http
      .get(`/api/v1/nodes`, {observe: "response"})
      .map((res: HttpResponse<Array<INode>>) => res.body)
  }

  getNodeByName(nodeName: string): Observable<INodeDetail> {
    return this.http
      .get(`/api/v1/node`, {
        observe: "response",
        params: {
          'node_name': nodeName
        }
      })
      .map((res: HttpResponse<INodeDetail>) => res.body)
  }

  toggleNodeStatus(nodeName: string, status: boolean): Observable<HttpResponse<Object>> {
    return this.http
      .get(`/api/v1/node/toggle`, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {
          'node_name': nodeName,
          'node_status': status ? "1" : "0"
        }
      });
  }

  getNodeGroupsOfOneNode(nodeName: string): Observable<Array<string>> {
    return this.http.get<Array<string>>(`/api/v1/node/0/group`,
      {observe: "response", params: {node_name: nodeName}})
      .map((res: HttpResponse<Array<string>>) => res.body || [])
  }

  addNodeToNodeGroup(nodeName:string,nodeGroupName:string): Observable<Object> {
    return this.http.post<Object>(`/api/v1/node/0/group`, null,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {node_name: nodeName, groupname: nodeGroupName}
      })
      .map((res: HttpResponse<Object>) => res.body)
  }

  deleteNodeToNodeGroup(nodeName:string,nodeGroupName:string): Observable<Object> {
    return this.http.delete<Object>(`/api/v1/node/0/group`,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {node_name: nodeName, groupname: nodeGroupName}
      })
      .map((res: HttpResponse<Object>) => res.body)
  }

  getNodeGroups(): Observable<Array<INodeGroup>> {
    return this.http.get<Array<INodeGroup>>(`/api/v1/nodegroup`, {observe: "response"})
      .map((res: HttpResponse<Array<INodeGroup>>) => res.body || [])
  }

  addNodeGroup(group: INodeGroup): Observable<HttpResponse<Object>> {
    return this.http.post(`/api/v1/nodegroup`, group, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    })
  }

  deleteNodeGroup(groupId: number, nodeGroupName: string): Observable<HttpResponse<Object>> {
    return this.http.delete(`/api/v1/nodegroup/${groupId}`,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {groupname: nodeGroupName}
      })
  }

  checkNodeGroupExist(groupName: string): Observable<HttpResponse<Object>> {
    return this.http.get(`/api/v1/nodegroup/existing`,
      {observe: "response", params: {nodegroup_name: groupName}})
  }

}