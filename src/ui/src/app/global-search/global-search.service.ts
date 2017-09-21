import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { AppInitService } from '../app.init.service';

@Injectable()
export class GlobalSearchService {

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type','application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  constructor(
    private http: Http,
    private appInitService: AppInitService
  ) {}

  search(content: string): Promise<any>{
    return this.http.get("/api/v1/search", { 
        headers: this.defaultHeader,
        params: {
          q: content,
          token: this.appInitService.token
        }
      })
      .toPromise()
      .then(res=>{
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err=>Promise.reject(err));
  }

}