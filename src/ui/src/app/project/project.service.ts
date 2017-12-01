import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Project } from './project';
import { Member } from './member/member';

import { AppInitService } from '../app.init.service';

@Injectable()
export class ProjectService {

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type','application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  constructor(
    private http: Http,
    private appInitService: AppInitService
  ){}
  
  getProjects(projectName?: string, pageIndex?: number, pageSize?: number): Promise<any> {
    
    return this.http
      .get('/api/v1/projects', {
        headers: this.defaultHeader,
        params: {
          'project_name': projectName,
          'page_index': pageIndex,
          'page_size': pageSize
        }
      })
      .toPromise()
      .then(resp=>{
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err=>Promise.reject(err));
  }

  createProject(project: Project): Promise<any> {    
    return this.http
      .post('/api/v1/projects', project, {
        headers: this.defaultHeader,
      })
      .toPromise()
      .then(resp=>this.appInitService.chainResponse(resp))
      .catch(err=>Promise.reject(err));
  }

  togglePublicity(projectId:number,projectPublic:number): Promise<any> {
    return this.http
      .put(`/api/v1/projects/${projectId}/publicity`, {
        'project_public': projectPublic
      }, {
        headers: this.defaultHeader
      })
      .toPromise()
      .then(resp=>this.appInitService.chainResponse(resp))
      .catch(err=>Promise.reject(err));
  }

  deleteProject(project: Project): Promise<any> {
    return this.http
      .delete(`/api/v1/projects/${project.project_id}`, {
        headers: this.defaultHeader
      })
      .toPromise()
      .then(resp=>this.appInitService.chainResponse(resp))
      .catch(err=>Promise.reject(err));
  }

  getProjectMembers(projectId: number): Promise<Member[]> {
    return this.http
      .get(`/api/v1/projects/${projectId}/members`, {
        headers: this.defaultHeader
      })
      .toPromise()
      .then(resp=>{
        this.appInitService.chainResponse(resp);
        return <Member[]>resp.json();
      })
      .catch(err=>Promise.reject(err));
  }

  addOrUpdateProjectMember(projectId: number, userId: number, roleId: number): Promise<any> {
    return this.http.post(`/api/v1/projects/${projectId}/members`, {
        'project_member_role_id': roleId,
        'project_member_user_id': userId
      },{
        headers: this.defaultHeader
      })
      .toPromise()
      .then(resp=>this.appInitService.chainResponse(resp))
      .catch(err=>Promise.reject(err));
  }

  deleteProjectMember(projectId: number, userId: number): Promise<any> {
    return this.http
      .delete(`/api/v1/projects/${projectId}/members/${userId}`, {
        headers: this.defaultHeader
      })
      .toPromise()
      .then(resp=>this.appInitService.chainResponse(resp))
      .catch(err=>Promise.reject(err));
  }

  getAvailableMembers(): Promise<Member[]> {
    return this.http
      .get('/api/v1/users', {
        headers: this.defaultHeader
      })
      .toPromise()
      .then(resp=>{
        this.appInitService.chainResponse(resp);
        let members = Array<Member>();
        let users = <any[]>resp.json();
        users.forEach(u => {
          if(u.user_deleted === 0) {
            let m = new Member();
            m.project_member_username = u.user_name;
            m.project_member_user_id = u.user_id;
            m.project_member_role_id = 1;
            members.push(m);
          }
        });
        return members;
      })
      .catch(err=>Promise.reject(err));
  }
}