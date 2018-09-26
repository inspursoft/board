import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Project } from './project';
import { Observable } from "rxjs/Observable";

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

  togglePublicity(projectId: number, projectPublic: number): Promise<any> {
    return this.http
      .put(`/api/v1/projects/${projectId}/publicity`, {
        'project_public': projectPublic
      }, {observe: "response"})
      .toPromise()
  }

  deleteProject(project: Project): Observable<any> {
    return this.http.delete(`/api/v1/projects/${project.project_id}`, {observe: "response"})
  }
}