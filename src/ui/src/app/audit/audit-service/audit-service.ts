import {Injectable} from "@angular/core";
import {HttpClient, HttpResponse} from "@angular/common/http";
import {Audit, Query} from "../audit";

const BASE_URL = "/api/v1";

@Injectable()
export class AuditService {
  constructor(private http: HttpClient) {
  }

  getUserList(): Promise<Object> {
    return this.http.get(`${BASE_URL}/users`, {
      observe: "response",
    }).toPromise().then((res: HttpResponse<Object>) => res.body);
  }

  getAuditList(querydata: Query): Promise<any> {
    return this.http
      .get(`${BASE_URL}/operations`, {
        observe: "response",
        params: {
          'page_index': querydata.pageIndex.toString(),
          'page_size': querydata.pageSize.toString(),
          "order_field": querydata.sortBy,
          "order_asc": querydata.isReverse ? "0" : "1",
          "operation_fromdate": querydata.beginDate,
          "operation_todate": querydata.endDate,
          "operation_status": status,
          "operation_user": querydata.user_name,
          "operation_action": querydata.action,
          "operation_object": querydata.object_name
        }
      })
      .toPromise()
      .then((res: HttpResponse<Object>) => res.body);
  }


}
