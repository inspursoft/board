import { Injectable } from "@angular/core";
import { HttpClient, HttpHeaders, HttpResponse } from "@angular/common/http";
import { AuditQueryData } from "./audit";
import { Observable } from "rxjs/Observable";
import { User } from "../user-center/user";
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from "../shared/shared.const";

const BASE_URL = "/api/v1";

@Injectable()
export class OperationAuditService {
  constructor(private http: HttpClient) {
  }

  getUserList(): Observable<Array<User>> {
    return this.http.get(`${BASE_URL}/users`, {
      observe: "response",
    }).map((res: HttpResponse<Array<User>>) => res.body);
  }

  getAuditList(querydata: AuditQueryData): Observable<any> {
    return this.http
      .get(`${BASE_URL}/operations`, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {
          'page_index': querydata.pageIndex.toString(),
          'page_size': querydata.pageSize.toString(),
          "order_field": querydata.sortBy,
          "order_asc": querydata.isReverse ? "0" : "1",
          "operation_fromdate": Math.ceil(querydata.beginTimestamp / 1000).toString(),
          "operation_todate": Math.ceil(querydata.endTimestamp / 1000).toString(),
          "operation_status": querydata.status,
          "operation_user": querydata.user_name,
          "operation_action": querydata.action,
          "operation_object": querydata.object_name
        }
      })
      .map((res: HttpResponse<Object>) => res.body);
  }
}
