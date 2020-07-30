import { Injectable } from '@angular/core';
import { HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { AppInitService } from './app-init.service';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import {
  PersistentVolume,
  PersistentVolumeClaim,
  SharedConfigMap,
  SharedConfigMapDetail,
  SharedCreateProject,
  SharedMember,
  SharedProject, User
} from '../shared/shared.types';

@Injectable()
export class SharedService {
  public showMaxGrafanaWindow = false;

  constructor(private http: ModelHttpClient,
              private appInitService: AppInitService) {

  }

  signOut(username: string): Observable<any> {
    return this.http.get('/api/v1/log-out', {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        params: {username}
      }
    );
  }

  getAssignedMembers(projectId: number): Observable<Array<SharedMember>> {
    return this.http.getArray(`/api/v1/projects/${projectId}/members`, SharedMember);
  }

  getAvailableMembers(projectId: number): Observable<Array<SharedMember>> {
    return this.http.getArray(`/api/v1/projects/${projectId}/members`, User, {
      param: {type: 'available'}
    }).pipe(map((res: Array<User>) => {
      const members = Array<SharedMember>();
      if (res && res.length > 0) {
        res.forEach(u => {
          if (u.userDeleted === 0) {
            const m = new SharedMember();
            m.userName = u.userName;
            m.userId = u.userId;
            m.roleId = 1;
            m.isMember = false;
            members.push(m);
          }
        });
      }
      return members;
    }));
  }

  getAllProjects(): Observable<Array<SharedProject>> {
    return this.http.getArray('/api/v1/projects', SharedProject);
  }

  getPVList(): Observable<Array<PersistentVolume>> {
    return this.http.getArray(`/api/v1/pvolumes`, PersistentVolume);
  }

  createNewPvc(pvc: PersistentVolumeClaim): Observable<any> {
    return this.http.post(`/api/v1/pvclaims`, pvc.getPostBody());
  }

  addOrUpdateProjectMember(projectId: number, userId: number, roleId: number): Observable<any> {
    return this.http.post(`/api/v1/projects/${projectId}/members`, {
      project_member_role_id: roleId,
      project_member_user_id: userId
    }, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  deleteProjectMember(projectId: number, userId: number): Observable<any> {
    return this.http.delete(`/api/v1/projects/${projectId}/members/${userId}`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  createProject(project: SharedCreateProject): Observable<any> {
    return this.http
      .post('/api/v1/projects', project.getPostBody(), {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response'
      });
  }

  checkPvcNameExist(projectName: string, pvcName: string): Observable<any> {
    return this.http.get(`/api/v1/pvclaims/existing`, {
      observe: 'response',
      params: {project_name: projectName, pvc_name: pvcName}
    });
  }

  getConfigMapList(projectName: string, pageIndex, pageSize: number): Observable<Array<SharedConfigMap>> {
    return this.http.getArray(`/api/v1/configmaps`, SharedConfigMap, {
      param: {
        project_name: projectName,
        configmap_list_page: pageIndex.toString(),
        configmap_list_page_size: pageSize.toString()
      }
    });
  }

  getConfigMapDetail(configMapName, projectName: string): Observable<SharedConfigMapDetail> {
    return this.http.getJson(`/api/v1/configmaps`, SharedConfigMapDetail, {
      param: {
        project_name: projectName,
        configmap_name: configMapName,
      }
    });
  }
}
