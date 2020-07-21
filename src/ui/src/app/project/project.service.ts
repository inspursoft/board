import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { PaginationProject } from './project.types';
import { SharedProject } from '../shared/shared.types';

@Injectable()
export class ProjectService {
  constructor(private http: ModelHttpClient) {
  }

  getProjects(projectName: string, pageIndex: number, pageSize: number, sortBy: string, isReverse: boolean): Observable<any> {
    return this.http.getPagination('/api/v1/projects', PaginationProject, {
      param: {
        project_name: projectName,
        page_index: pageIndex.toString(),
        page_size: pageSize.toString(),
        order_field: sortBy,
        order_asc: isReverse ? '0' : '1'
      }
    });
  }

  togglePublicity(projectId: number, projectPublic: number): Observable<any> {
    return this.http.put(`/api/v1/projects/${projectId}/publicity`, {
      project_public: projectPublic
    });
  }

  deleteProject(project: SharedProject): Observable<any> {
    return this.http.delete(`/api/v1/projects/${project.projectId}`, {observe: 'response'});
  }
}
