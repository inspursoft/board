import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import { Member, Project } from '../project/project';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { PersistentVolume, PersistentVolumeClaim } from '../shared/shared.types';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { AppInitService } from "./app-init.service";

@Injectable()
export class SharedService {
  public showMaxGrafanaWindow = false;

  constructor(private http: HttpClient,
              private appInitService: AppInitService) {

  }

  signOut(username: string): Observable<any> {
    return this.http.get('/api/v1/log-out', {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        params: {
          'username': username
        }
      }
    );
  }

  search(content: string): Observable<any> {
    return this.http.get("/api/v1/search", {
      observe: "response",
      params: {
        q: content,
        token: this.appInitService.token
      }
    }).pipe(map(res => res.body));
  }

  getAssignedMembers(projectId: number): Observable<Array<Member>> {
    return this.http
      .get<Array<Member>>(`/api/v1/projects/${projectId}/members`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<Member>>) => res.body || []));
  }

  getAvailableMembers(projectId: number): Observable<Array<Member>> {
    return this.http.get(`/api/v1/projects/${projectId}/members`, {
      params: {
        type: 'available'
      }
    }).pipe(map((res: Array<any>) => {
        const members = Array<Member>();
        res.forEach(u => {
          if (u.user_deleted === 0) {
            const m = new Member();
            m.project_member_username = u.user_name;
            m.project_member_user_id = u.user_id;
            m.project_member_role_id = 1;
            m.isMember = false;
            members.push(m);
          }
        });
        return members;
      }));
  }

  getOneProject(projectName: string): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects', {
      observe: 'response',
      params: {project_name: projectName}
    }).pipe(map(res => res.body));
  }

  getAllProjects(): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects', {observe: 'response'}).pipe(map(res => res.body));
  }

  getAllPvList(): Observable<Array<PersistentVolume>> {
    return this.http.get(`/api/v1/pvolumes`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<object>>) => {
        const result: Array<PersistentVolume> = Array<PersistentVolume>();
        res.body.forEach(resObject => {
          const persistentVolume = new PersistentVolume();
          persistentVolume.initFromRes(resObject);
          result.push(persistentVolume);
        });
        return result;
      }));
  }

  createNewPvc(pvc: PersistentVolumeClaim): Observable<any> {
    return this.http.post(`/api/v1/pvclaims`, pvc.postObject(), {observe: 'response'});
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

  createProject(project: Project): Observable<any> {
    return this.http
      .post('/api/v1/projects', project, {
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
}
