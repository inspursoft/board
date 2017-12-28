import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers, Response } from "@angular/http";
import { BuildImageData, Image, ImageDetail } from "../image";
import { AppInitService } from "../../app.init.service";
import "rxjs/operator/toPromise";
import { Project } from "app/project/project";

@Injectable()
export class ImageService {
  constructor(private http: Http,
              private appInitService: AppInitService) {
  }

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  uploadFile(formData: FormData): Promise<any> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/upload`, formData, options).toPromise()
      .then(resp => this.appInitService.chainResponse(resp))
      .catch(err => Promise.reject(err));
  }

  getProjects(projectName?: string): Promise<Project[]> {
    return this.http.get('/api/v1/projects', {
      headers: this.defaultHeader,
      params: {'project_name': projectName, 'member_only': 1}
    }).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err => Promise.reject(err));
  }

  uploadDockerFile(formData: FormData): Promise<any> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/images/dockerfile/upload`, formData, options).toPromise()
      .then(resp => this.appInitService.chainResponse(resp))
      .catch(err => Promise.reject(err));
  }

  downloadDockerFile(fileInfo: {imageName: string, tagName: string, projectName: string}): Promise<any> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({
      headers: headers,
      params: {
        image_name: fileInfo.imageName,
        tag_name: fileInfo.tagName,
        project_name: fileInfo.projectName
      }
    });
    return this.http.get(`/api/v1/images/dockerfile/download`, options).toPromise()
      .then(res => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }


  removeFile(formData: FormData): Promise<any> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/remove`, formData, options).toPromise()
      .then(resp => this.appInitService.chainResponse(resp))
      .catch(err => Promise.reject(err));
  }

  cancelConsole(jobName: string): Promise<any> {
    return this.http
      .get(`/api/v1/jenkins-job/stop`, {
        headers: this.defaultHeader,
        params: {
          "job_name": jobName
        }
      }).toPromise()
      .then((res: Response) => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  buildImageFromTemp(imageData: BuildImageData): Promise<any> {
    return this.http.post(`/api/v1/images/building`, imageData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(res => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  buildImageFromDockerFile(fileInfo: {imageName: string, tagName: string, projectName: string}): Promise<any> {
    return this.http.post(`/api/v1/images/dockerfilebuilding`, fileInfo, {
      headers: this.defaultHeader,
      params: {
        image_name: fileInfo.imageName,
        image_tag: fileInfo.tagName,
        project_name: fileInfo.projectName
      }
    }).toPromise()
      .then(res => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  getFileList(formData: FormData): Promise<Array<{path: string, file_name: string, size: number}>> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/list`, formData, options).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err => Promise.reject(err));
  }

  getDockerFilePreview(imageData: BuildImageData): Promise<string> {
    return this.http.post(`/api/v1/images/preview`, imageData, {
      headers: this.defaultHeader
    }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.text();
      })
      .catch(err => Promise.reject(err));
  }

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Promise<Image[]> {
    let options = new RequestOptions({
      headers: this.defaultHeader,
      params: {
        'image_name': image_name,
        'image_list_page': image_list_page.toString(),
        'image_list_page_size': image_list_page_size.toString()
      }
    });
    return this.http.get("/api/v1/images", options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json() || [];
      })
      .catch(err => Promise.reject(err));
  }

  getImageDetailList(image_name: string): Promise<ImageDetail[]> {
    return this.http.get(`/api/v1/images/${image_name}`, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  deleteImages(imageName: string): Promise<any> {
    return this.http
      .delete(`/api/v1/images`, {
        headers: this.defaultHeader,
        params: {image_name: imageName}
      })
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res;
      })
      .catch(err => Promise.reject(err));
  }

  deleteImageTag(imageName: string, imageTag: string): Promise<any> {
    return this.http
      .delete(`/api/v1/images/${imageName}`, {
        headers: this.defaultHeader,
        params: {image_tag: imageTag}
      })
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res;
      })
      .catch(err => Promise.reject(err));
  }

  checkImageExist(projectName: string, imageName: string, imageTag: string): Promise<any> {
    return this.http.get(`/api/v1/images/${imageName}/existing`, {
      headers: this.defaultHeader,
      params: {image_tag: imageTag, project_name: projectName}
    }).toPromise()
      .then((res: Response) => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }
}