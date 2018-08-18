import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Project } from './project';
import { Member } from './member/member';

@Injectable()
export class ProjectService {
  constructor(private http: HttpClient) {
  }

  getProjects(projectName: string, pageIndex: number, pageSize: number, sortBy: string, isReverse: boolean): Promise<any> {
    return this.http
      .get('/api/v1/projects', {
        observe: "response",
        params: {
          'project_name': projectName,
          'page_index': pageIndex.toString(),
          'page_size': pageSize.toString(),
          "order_field": sortBy,
          "order_asc": isReverse ? "0" : "1"
        }
      })
      .toPromise()
      .then(res => res.body)
  }

  createProject(project: Project): Promise<any> {
    return this.http
      .post('/api/v1/projects', project, {observe: "response"})
      .toPromise()
  }

  togglePublicity(projectId: number, projectPublic: number): Promise<any> {
    return this.http
      .put(`/api/v1/projects/${projectId}/publicity`, {
        'project_public': projectPublic
      }, {observe: "response"})
      .toPromise()
  }

  deleteProject(project: Project): Promise<any> {
    return this.http
      .delete(`/api/v1/projects/${project.project_id}`, {observe: "response"})
      .toPromise()
  }

  getProjectMembers(projectId: number): Promise<Member[]> {
    return this.http
      .get<Member[]>(`/api/v1/projects/${projectId}/members`, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<Member[]>) => res.body || [])
  }

  addOrUpdateProjectMember(projectId: number, userId: number, roleId: number): Promise<any> {
    return this.http.post(`/api/v1/projects/${projectId}/members`, {
      'project_member_role_id': roleId,
      'project_member_user_id': userId
    }, {observe: "response"})
      .toPromise()
  }

  deleteProjectMember(projectId: number, userId: number): Promise<any> {
    return this.http
      .delete(`/api/v1/projects/${projectId}/members/${userId}`, {observe: "response"})
      .toPromise()
  }

  getAvailableMembers(): Promise<Member[]> {
    return this.http
      .get('/api/v1/users', {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<Object>) => {
        let members = Array<Member>();
        let users = res.body as Array<any>;
        users.forEach(u => {
          if (u.user_deleted === 0) {
            let m = new Member();
            m.project_member_username = u.user_name;
            m.project_member_user_id = u.user_id;
            m.project_member_role_id = 1;
            members.push(m);
          }
        });
        return members;
      })
      .catch(err => Promise.reject(err));
  }
}