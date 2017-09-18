import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { AppInitService } from '../app.init.service';

@Injectable()
export class GlobalSearchService {

  defaultHeaders: Headers = new Headers({contentType: 'application/json'});

  constructor(
    private http: Http,
    private appInitService: AppInitService
  ) {}

  search(content: string): Promise<any>{
    return this.http.get("/api/v1/search", { 
        headers: this.defaultHeaders,
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