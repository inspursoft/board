import { Injectable } from "@angular/core";
import { Observable } from "rxjs/Observable";
import { HttpClient } from "@angular/common/http";

@Injectable()
export class GrafanaService {
  constructor(private http: HttpClient) {
  }

  testGrafana(grafanaUrl: string): Observable<any> {
    return this.http.get(grafanaUrl, {observe: "response"})
  }
}
