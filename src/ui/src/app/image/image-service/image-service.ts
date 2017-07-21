import { Injectable } from "@angular/core"
import { Http, RequestOptions, Headers, Response } from "@angular/http"
import { Image,ImageDetail } from "../image"
import { MessageService } from "../../shared/message-service/message.service";
import "rxjs/operator/toPromise"

@Injectable()
export class ImageService {
  constructor(private http: Http,
              private messageService: MessageService) {
  }

  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  static getErrorMsg(reason: Response | Error, statusArr: Array<number>, errorKey: string): string {
    if (reason instanceof Response) {
      return statusArr.indexOf(reason.status) > -1 ?
        `USER_CENTER.${errorKey}_ERR_${reason.status}` :
        `${reason.status}:${reason.statusText}`;
    }
    else {
      return `${reason.name}:${reason.message}`;
    }
  }

  handleGetError(reason:Response | Error){
    let errMsg: string = ImageService.getErrorMsg(reason, Array.from([400, 404]), "GET");
    this.messageService.dispatchError(reason, errMsg);
    return Promise.reject(errMsg);
  }

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Promise<Image[]> {
    let params: Map<string, string> = new Map<string, string>();
    params["image_name"] = image_name;
    params["image_list_page"] = image_list_page.toString();
    params["image_list_page_size"] = image_list_page_size.toString();
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      search: params
    });
    return this.http.get("/api/v1/images", options).toPromise()
      .then(res => res.json())
      .catch(this.handleGetError);
  }

  getImageDetailList(image_name: string): Promise<ImageDetail[]> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.get(`/api/v1/images/${image_name}`, options).toPromise()
      .then(res => res.json())
      .catch(this.handleGetError);
  }

}