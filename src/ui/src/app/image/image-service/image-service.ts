import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers } from "@angular/http";
import { Image, ImageDetail } from "../image";
import { AppInitService } from "../../app.init.service";
import "rxjs/operator/toPromise";

@Injectable()
export class ImageService {
  constructor(private http: Http,
              private appInitService: AppInitService) {
  }

  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Promise<Image[]> {
    let params: Map<string, string> = new Map<string, string>();
    params["image_name"] = image_name;
    params["image_list_page"] = image_list_page.toString();
    params["image_list_page_size"] = image_list_page_size.toString();
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {'token': this.appInitService.token},
      search: params
    });
    return this.http.get("/api/v1/images", options).toPromise()
      .then(res => res.json())
      .catch(err => Promise.reject(err));
  }

  getImageDetailList(image_name: string): Promise<ImageDetail[]> {
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {'token': this.appInitService.token}
    });
    return this.http.get(`/api/v1/images/${image_name}`, options)
      .timeout(3000)
      .toPromise()
      .then(res => {
        let s = res.json();
        let result: ImageDetail[] = Array<ImageDetail>();
        s.forEach(item => {
          let image_creationtime = JSON.parse(item["image_creationtime"]);
          let image_author = JSON.parse(item["image_author"]);
          result.push({
            image_name: item["image_name"],
            image_tag: item["image_tag"],
            image_author: image_author["author"],
            image_id: (item["image_id"] as string).replace(/sha256:/g, ""),
            image_creationtime: image_creationtime["created"],
            image_size_number: (item["image_size_number"] / (1024 * 1024)).toFixed(2),
            image_size_unit: "MB"
          })
        });
        return result;
      })
      .catch(err => Promise.reject(err));
  }
}