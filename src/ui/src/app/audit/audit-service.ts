import { Injectable } from '@angular/core';
import { HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AuditPagination, AuditQueryData } from './audit';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { User } from '../shared/shared.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

const BASE_URL = '/api/v1';

@Injectable()
export class OperationAuditService {
  constructor(private http: ModelHttpClient) {
  }

  getUserList(): Observable<Array<User>> {
    return this.http.getArray(`${BASE_URL}/users`, User);
  }

  getAuditList(querydata: AuditQueryData): Observable<AuditPagination> {
    return this.http
      .getPagination(`${BASE_URL}/operations`, AuditPagination, {
        header: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        param: {
          page_index: querydata.pageIndex.toString(),
          page_size: querydata.pageSize.toString(),
          order_field: querydata.sortBy,
          order_asc: querydata.isReverse ? '0' : '1',
          operation_fromdate: Math.ceil(querydata.beginTimestamp / 1000).toString(),
          operation_todate: Math.ceil(querydata.endTimestamp / 1000).toString(),
          operation_status: querydata.status,
          operation_user: querydata.userName,
          operation_action: querydata.action,
          operation_object: querydata.objectName
        }
      });
  }
}
