import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { HttpClient } from "@angular/common/http";

@Injectable()
export class KibanaService {
  constructor(private http: HttpClient) {
  }

  testKibana(kibanaUrl: string): Observable<any> {
    return this.http.get(kibanaUrl, {observe: "response"})
  }
}
