import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import { Member, Project } from '../project/project';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from './shared.const';
import { PersistentVolume, PersistentVolumeClaim } from './shared.types';
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

@Injectable()
export class SharedService {
  public showMaxGrafanaWindow = false;

  constructor(private http: HttpClient) {

  }

  getProjectMembers(projectId: number): Observable<Array<Member>> {
    return this.http
      .get<Array<Member>>(`/api/v1/projects/${projectId}/members`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<Member>>) => res.body || []));
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
      .pipe(map((res: HttpResponse<Array<Object>>) => {
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

  getAvailableMembers(): Observable<Array<Member>> {
    return this.http.get('/api/v1/users', {observe: 'response'})
      .pipe(map((res: HttpResponse<Object>) => {
        const members = Array<Member>();
        const users = res.body as Array<any>;
        users.forEach(u => {
          if (u.user_deleted === 0) {
            const m = new Member();
            m.project_member_username = u.user_name;
            m.project_member_user_id = u.user_id;
            m.project_member_role_id = 1;
            members.push(m);
          }
        });
        return members;
      }));
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
    return this.http.get(`/api/v1/pvclaims/existing`, {observe: 'response', params: {project_name: projectName, pvc_name: pvcName}});
  }
}
