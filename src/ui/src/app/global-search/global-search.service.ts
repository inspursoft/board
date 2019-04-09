import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { AppInitService } from '../app.init.service';
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

@Injectable()
export class GlobalSearchService {

  constructor(
    private http: HttpClient,
    private appInitService: AppInitService
  ) {}

  search(content: string): Observable<any>{
    return this.http.get("/api/v1/search", {
        observe:"response",
        params: {
          q: content,
          token: this.appInitService.token
        }
      }).pipe(map(res=> res.body));
  }

}
