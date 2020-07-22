import { Observable } from 'rxjs';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { Injectable } from '@angular/core';
import { GlobalSearchResult } from './global-search.types';

@Injectable()
export class GlobalSearchService {
  constructor(private http: ModelHttpClient) {
  }

  search(q, token: string): Observable<GlobalSearchResult> {
    return this.http.getJson('/api/v1/search', GlobalSearchResult, {param: {q, token}});
  }
}

