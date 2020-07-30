import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import {
  NFSPersistentVolume,
  PersistentVolume,
  PersistentVolumeClaim,
  PersistentVolumeClaimDetail,
  RBDPersistentVolume
} from './sotrage.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

@Injectable()
export class StorageService {

  constructor(private http: ModelHttpClient) {

  }

  getPvcList(pvcName: string, pvcListPage: number, pvcListPageSize: number): Observable<Array<PersistentVolumeClaim>> {
    return this.http.getArray(`/api/v1/pvclaims`, PersistentVolumeClaim, {
      param: {
        pvc_name: pvcName,
        pvc_list_page: pvcListPage.toString(),
        pvc_list_page_size: pvcListPageSize.toString()
      }
    });
  }

  deletePvc(pvcId: number): Observable<any> {
    return this.http.delete(`/api/v1/pvclaims/${pvcId}`);
  }

  getPvcDetailInfo(pvcId: number): Observable<PersistentVolumeClaimDetail> {
    return this.http.getJson(`/api/v1/pvclaims/${pvcId}`, PersistentVolumeClaimDetail);
  }

  getPvList(pvName: string, pvListPage: number, pvListPageSize: number): Observable<Array<PersistentVolume>> {
    return this.http.getArray(`/api/v1/pvolumes`, PersistentVolume, {
      param: {
        pv_name: pvName,
        pv_listPage: pvListPage.toString(),
        pv_list_page_size: pvListPageSize.toString()
      }
    });
  }

  createNewPv(body: PersistentVolume): Observable<any> {
    return this.http.postJson(`/api/v1/pvolumes`, PersistentVolume, body.getPostBody());
  }

  deletePv(pvId: number): Observable<Object> {
    return this.http.delete(`/api/v1/pvolumes/${pvId}`, {observe: 'response'});
  }

  checkPvNameExist(pvName: string): Observable<any> {
    return this.http.get(`/api/v1/pvolumes/existing`, {observe: 'response', params: {pv_name: pvName}});
  }

  getPvDetailInfo(id: number): Observable<PersistentVolume> {
    return this.http.getJson(`/api/v1/pvolumes/${id}`, PersistentVolume).pipe(
      map((pv: PersistentVolume) => {
        if (pv.type === 1) {
          const result = new NFSPersistentVolume(pv.res);
          result.initFromRes();
          return result;
        } else if (pv.type === 2) {
          const result = new RBDPersistentVolume(pv.res);
          result.initFromRes();
          return result;
        }
      })
    );
  }
}
