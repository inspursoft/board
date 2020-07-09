import { Injectable } from '@angular/core';
import { HttpEvent, HttpHeaders, HttpRequest, HttpResponse } from '@angular/common/http';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { BuildImageData, Image, ImageDetail, ImageProject } from './image.types';

@Injectable()
export class ImageService {
  constructor(private http: ModelHttpClient) {
  }

  uploadFile(formData: FormData): Observable<HttpEvent<object>> {
    const req = new HttpRequest('POST', `/api/v1/files/upload`, formData, {
      reportProgress: true,
    });
    return this.http.request<object>(req);
  }

  getProjects(projectName: string = ''): Observable<Array<ImageProject>> {
    return this.http.getArray('/api/v1/projects', ImageProject, {
      param: {project_name: projectName, member_only: '1'}
    });
  }

  uploadDockerFile(formData: FormData): Observable<string> {
    return this.http.post(`/api/v1/images/dockerfile/upload`, formData, {observe: 'response', responseType: 'text'})
      .pipe(map((res: HttpResponse<string>) => res.body));
  }

  downloadDockerFile(fileInfo: { imageName: string, tagName: string, projectName: string }): Observable<any> {
    return this.http.get(`/api/v1/images/dockerfile/download`,
      {
        observe: 'response',
        responseType: 'text',
        params: {
          image_name: fileInfo.imageName,
          image_tag: fileInfo.tagName,
          project_name: fileInfo.projectName
        }
      });
  }


  removeFile(formData: FormData): Observable<any> {
    return this.http.post(`/api/v1/files/remove`, formData, {observe: 'response'});
  }

  cancelConsole(jobName: string): Observable<any> {
    return this.http
      .get(`/api/v1/jenkins-job/stop`, {
        observe: 'response',
        params: {
          job_name: jobName
        }
      });
  }

  buildImageFromTemp(imageData: BuildImageData): Observable<any> {
    return this.http.post(`/api/v1/images/building`, imageData.getPostBody(), {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  buildImageFromDockerFile(fileInfo: { imageName: string, tagName: string, projectName: string }): Observable<any> {
    return this.http.post(`/api/v1/images/dockerfilebuilding`, fileInfo, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response',
      params: {
        image_name: fileInfo.imageName,
        image_tag: fileInfo.tagName,
        project_name: fileInfo.projectName
      }
    });
  }

  buildImageFromImagePackage(params: {
    imageName: string,
    tagName: string,
    projectName: string,
    imagePackageName: string
  }): Observable<any> {
    return this.http.post(`/api/v1/images/imagepackage`, null, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response',
      params: {
        image_name: params.imageName,
        image_tag: params.tagName,
        project_name: params.projectName,
        image_package_name: params.imagePackageName
      }
    });
  }

  getFileList(formData: FormData): Observable<Array<{ path: string, file_name: string, size: number }>> {
    return this.http.post(`/api/v1/files/list`, formData, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<{ path: string, file_name: string, size: number }>>) => res.body));
  }

  getDockerFilePreview(imageData: BuildImageData): Observable<string> {
    return this.http.post(`/api/v1/images/preview`, imageData.getPostBody(), {observe: 'response', responseType: 'text'})
      .pipe(map(res => res.body));
  }

  getImages(imageName?: string, imageListPage?: number, imageListPageSize?: number): Observable<Array<Image>> {
    return this.http.getArray('/api/v1/images', Image, {
      param: {
        image_name: imageName,
        image_list_page: imageListPage.toString(),
        image_list_page_size: imageListPageSize.toString()
      }
    });
  }

  getImageDetailList(imageName: string): Observable<Array<ImageDetail>> {
    return this.http.getArray(`/api/v1/images/${imageName}`, ImageDetail);
  }

  deleteImages(imageName: string): Observable<any> {
    return this.http.delete(`/api/v1/images`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response',
      params: {image_name: imageName}
    });
  }

  deleteImageTag(imageName: string, imageTag: string): Observable<any> {
    return this.http.delete(`/api/v1/images/${imageName}`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response',
      params: {image_tag: imageTag}
    });
  }

  checkImageExist(projectName: string, imageName: string, imageTag: string): Observable<any> {
    return this.http.get(`/api/v1/images/${imageName}/existing`, {
      observe: 'response',
      params: {image_tag: imageTag, project_name: projectName}
    });
  }

  getBoardRegistry(): Observable<string> {
    return this.http.get(`/api/v1/images/registry`, {observe: 'response', responseType: 'text'})
      .pipe(map((obs: HttpResponse<string>) => obs.body));
  }

  deleteImageConfig(projectName: string): Observable<object> {
    return this.http.delete(`/api/v1/images/configclean`, {
      observe: 'response',
      params: {project_name: projectName}
    }).pipe(map((obs: HttpResponse<object>) => obs.body));
  }
}
