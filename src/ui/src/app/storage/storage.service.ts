import { Injectable } from "@angular/core";
import { Observable } from "rxjs/Observable";
import { HttpClient, HttpResponse } from "@angular/common/http";
import { NFSPersistentVolume, PersistentVolume, RBDPersistentVolume } from "../shared/shared.types";

@Injectable()
export class StorageService {

  constructor(private http: HttpClient) {

  }

  getPvList(pvName: string, pvListPage: number, pvListPageSize: number): Observable<Array<PersistentVolume>> {
    return this.http.get(`/api/v1/pvolumes`, {
      observe: "response", params: {
        pv_name: pvName,
        pv_listPage: pvListPage.toString(),
        pv_list_page_size: pvListPageSize.toString()
      }
    }).map((res: HttpResponse<Array<Object>>) => {
      let result: Array<PersistentVolume> = Array<PersistentVolume>();
      res.body.forEach(resObject => {
        let persistentVolume = new PersistentVolume();
        persistentVolume.id = Reflect.get(resObject, 'pv_id');
        persistentVolume.name = Reflect.get(resObject, 'pv_name');
        persistentVolume.type = Reflect.get(resObject, 'pv_type');
        persistentVolume.state = Reflect.get(resObject, 'pv_state');
        persistentVolume.capacity = Reflect.get(resObject, 'pv_capacity');
        persistentVolume.accessMode = Reflect.get(resObject, 'pv_accessmode');
        persistentVolume.reclaim = Reflect.get(resObject, 'pv_reclaim');
        result.push(persistentVolume);
      });
      return result;
    })
  }

  createNewPv(body: PersistentVolume): Observable<Object> {
    return this.http.post(`/api/v1/pvolumes`, body.postObject(), {observe: "response"})
  }

  deletePv(pvId: number): Observable<Object> {
    return this.http.delete(`/api/v1/pvolumes/${pvId}`, {observe: "response"})
  }

  getPvDetailInfo(id: number): Observable<PersistentVolume> {
    return this.http.get(`/api/v1/pvolumes/${id}`, {observe: "response"})
      .map((res: HttpResponse<Object>) => {
        let result: PersistentVolume;
        let options = Reflect.get(res.body, 'pv_options');
        if (Reflect.get(res.body, 'pv_type') == 1) {
          result = new NFSPersistentVolume();
          if (options) {
            (result as NFSPersistentVolume).options.server = Reflect.get(options, 'server');
            (result as NFSPersistentVolume).options.path = Reflect.get(options, 'path');
          }
        } else if (Reflect.get(res.body, 'pv_type') == 2) {
          result = new RBDPersistentVolume();
          if (options) {
            (result as RBDPersistentVolume).options.user = Reflect.get(options, 'user');
            (result as RBDPersistentVolume).options.keyring = Reflect.get(options, 'keyring');
            (result as RBDPersistentVolume).options.pool = Reflect.get(options, 'pool');
            (result as RBDPersistentVolume).options.image = Reflect.get(options, 'image');
            (result as RBDPersistentVolume).options.fstype = Reflect.get(options, 'fstype');
            (result as RBDPersistentVolume).options.monitors = Reflect.get(options, 'monitors');
            (result as RBDPersistentVolume).options.secretname = Reflect.get(options, 'secretname');
            (result as RBDPersistentVolume).options.secretnamespace = Reflect.get(options, 'secretnamespace');
          }
        }
        if (result) {
          result.name = Reflect.get(res.body, 'pv_name');
          result.type = Reflect.get(res.body, 'pv_type');
          result.state = Reflect.get(res.body, 'pv_state');
          result.capacity = Reflect.get(res.body, 'pv_capacity');
          result.accessMode = Reflect.get(res.body, 'pv_accessmode');
          result.reclaim = Reflect.get(res.body, 'pv_reclaim');
        }
        return result;
      })
  }


}