import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { Project } from './project';
import { Member } from './member/member';

@Injectable()
export class ProjectService {
  
  headers: Headers = new Headers({"Content-type": "application/json"});
  
  constructor(private http: Http){}
  
  getProjects(projectName?: string): Promise<Project[]> {
    return this.http
      .get('/api/v1/projects', {
        headers: this.headers,
        params: {
          'project_name': projectName
        }
      })
      .toPromise()
      .then(resp=><Project[]>resp.json())
      .catch(err=>Promise.reject(err));
  }

  createProject(project: Project): Promise<any> {    
    return this.http
      .post('/api/v1/projects', project)
      .toPromise()
      .then(resp=>resp)
      .catch(err=>Promise.reject(err));
  }

  togglePublicity(project: Project): Promise<any> {
    return this.http
      .put(`/api/v1/projects/${project.project_id}/publicity`, {
        project_public: project.project_public
      })
      .toPromise()
      .then(resp=>resp)
      .catch(err=>Promise.reject(err));
  }

  deleteProject(project: Project): Promise<any> {
    return this.http
      .delete(`/api/v1/projects/${project.project_id}`)
      .toPromise()
      .then(resp=>resp)
      .catch(err=>Promise.reject(err));
  }

  getProjectMembers(projectId: number): Promise<Member[]> {
    return this.http.get(`/api/v1/projects/${projectId}/members`)
      .toPromise()
      .then(resp=><Member[]>resp.json())
      .catch(err=>Promise.reject(err));
  }

  addOrUpdateProjectMember(projectId: number, userId: number, roleId: number): Promise<any> {
    return this.http.post(`/api/v1/projects/${projectId}/members`, {
      project_member_role_id: roleId,
      project_member_user_id: userId
    })
    .toPromise()
    .then(resp=>resp)
    .catch(err=>Promise.reject(err));
  }

  deleteProjectMember(projectId: number, userId: number): Promise<any> {
    return this.http.delete(`/api/v1/projects/${projectId}/members/${userId}`)
    .toPromise()
    .then(resp=>resp)
    .catch(err=>Promise.reject(err));
  }

  getAvailableMembers(): Promise<Member[]> {
    return this.http.get('/api/v1/users')
      .toPromise()
      .then(resp=>{
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