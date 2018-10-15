import { Injectable } from "@angular/core";
import { HttpClient, HttpEvent, HttpHeaders, HttpRequest, HttpResponse } from "@angular/common/http";
import { BuildImageData, Image, ImageDetail } from "../image";
import { Observable } from "rxjs/Observable";
import { Project } from "../../project/project";
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from "../../shared/shared.const";

@Injectable()
export class ImageService {
  constructor(private http: HttpClient) {
  }

  uploadFile(formData: FormData): Observable<HttpEvent<Object>> {
    const req = new HttpRequest('POST', `/api/v1/files/upload`, formData, {
      reportProgress: true,
    });
    return this.http.request<Object>(req)
  }

  getProjects(projectName: string = ""): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects', {
      observe: "response",
      params: {'project_name': projectName, 'member_only': "1"}
    }).map((res: HttpResponse<Array<Project>>) => res.body || [])
  }

  uploadDockerFile(formData: FormData): Observable<any> {
    return this.http.post(`/api/v1/images/dockerfile/upload`, formData, {observe: "response"})
  }

  downloadDockerFile(fileInfo: {imageName: string, tagName: string, projectName: string}): Observable<any> {
    return this.http.get(`/api/v1/images/dockerfile/download`,
      {
        observe: "response",
        responseType:"text",
        params: {
          image_name: fileInfo.imageName,
          image_tag: fileInfo.tagName,
          project_name: fileInfo.projectName
        }
      });
  }


  removeFile(formData: FormData): Observable<any> {
    return this.http.post(`/api/v1/files/remove`, formData, {observe: "response"})
  }

  cancelConsole(jobName: string): Observable<any> {
    return this.http
      .get(`/api/v1/jenkins-job/stop`, {
        observe: "response",
        params: {
          "job_name": jobName
        }
      })
  }

  buildImageFromTemp(imageData: BuildImageData): Observable<any> {
    return this.http.post(`/api/v1/images/building`, imageData, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    })
  }

  buildImageFromDockerFile(fileInfo: {imageName: string, tagName: string, projectName: string}): Observable<any> {
    return this.http.post(`/api/v1/images/dockerfilebuilding`, fileInfo, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response",
      params: {
        image_name: fileInfo.imageName,
        image_tag: fileInfo.tagName,
        project_name: fileInfo.projectName
      }
    });
  }

  getFileList(formData: FormData): Observable<Array<{path: string, file_name: string, size: number}>> {
    return this.http.post(`/api/v1/files/list`, formData, {observe: "response"})
      .map((res: HttpResponse<Array<{path: string, file_name: string, size: number}>>) => res.body)
  }

  getDockerFilePreview(imageData: BuildImageData): Observable<string> {
    return this.http.post(`/api/v1/images/preview`, imageData, {observe: "response", responseType: 'text'})
      .map(res => res.body)
  }

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Observable<Array<Image>> {
    return this.http.get<Array<Image>>("/api/v1/images", {
      observe: "response",
      params: {
        'image_name': image_name,
        'image_list_page': image_list_page.toString(),
        'image_list_page_size': image_list_page_size.toString()
      }
    }).map((res: HttpResponse<Array<Image>>) => res.body || [])
  }

  getImageDetailList(image_name: string): Observable<ImageDetail[]> {
    return this.http.get<ImageDetail[]>(`/api/v1/images/${image_name}`, {observe: "response"})
      .map((res: HttpResponse<ImageDetail[]>) => res.body || [])
  }

  deleteImages(imageName: string): Observable<any> {
    return this.http.delete(`/api/v1/images`, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {image_name: imageName}
      })
  }

  deleteImageTag(imageName: string, imageTag: string): Observable<any> {
    return this.http.delete(`/api/v1/images/${imageName}`, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response",
        params: {image_tag: imageTag}
      })
  }

  checkImageExist(projectName: string, imageName: string, imageTag: string): Observable<any> {
    return this.http.get(`/api/v1/images/${imageName}/existing`, {
      observe: "response",
      params: {image_tag: imageTag, project_name: projectName}
    })
  }

  getBoardRegistry(): Observable<string> {
    return this.http.get(`/api/v1/images/registry`, {observe: "response", responseType: "text"})
      .map((obs: HttpResponse<string>) => obs.body)
  }

  deleteImageConfig(projectName: string): Observable<Object> {
    return this.http.delete(`/api/v1/images/configclean`, {
      observe: "response",
      params: {project_name: projectName}
    }).map((obs: HttpResponse<Object>) => obs.body)
  }
}