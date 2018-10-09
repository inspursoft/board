import { Injectable } from "@angular/core"
import { HttpClient, HttpHeaders, HttpResponse } from "@angular/common/http";
import { Member, Project } from "../project/project";
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from "./shared.const";
import { Observable } from "rxjs/Observable";

@Injectable()
export class SharedService {
  constructor(private http: HttpClient) {

  }

  getProjectMembers(projectId: number): Observable<Array<Member>> {
    return this.http
      .get<Array<Member>>(`/api/v1/projects/${projectId}/members`, {observe: "response"})
      .map((res: HttpResponse<Array<Member>>) => res.body || [])
  }

  getOneProject(projectName: string): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects', {
        observe: "response",
        params: {'project_name': projectName}
      }).map(res => res.body)
  }

  getAvailableMembers(): Observable<Array<Member>> {
    return this.http.get('/api/v1/users', {observe: "response"})
      .map((res: HttpResponse<Object>) => {
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
      });
  }

  addOrUpdateProjectMember(projectId: number, userId: number, roleId: number): Observable<any> {
    return this.http.post(`/api/v1/projects/${projectId}/members`, {
      'project_member_role_id': roleId,
      'project_member_user_id': userId
    }, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    })
  }

  deleteProjectMember(projectId: number, userId: number): Observable<any> {
    return this.http.delete(`/api/v1/projects/${projectId}/members/${userId}`, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response"
      })
  }

  createProject(project: Project): Observable<any> {
    return this.http
      .post('/api/v1/projects', project, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response"
      })
  }
}