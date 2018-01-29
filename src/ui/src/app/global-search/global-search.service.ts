import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { AppInitService } from '../app.init.service';

@Injectable()
export class GlobalSearchService {

  constructor(
    private http: HttpClient,
    private appInitService: AppInitService
  ) {}

  search(content: string): Promise<any>{
    return this.http.get("/api/v1/search", {
        observe:"response",
        params: {
          q: content,
          token: this.appInitService.token
        }
      })
      .toPromise()
      .then(res=> res.body)
      .catch(err=>Promise.reject(err));
  }

}