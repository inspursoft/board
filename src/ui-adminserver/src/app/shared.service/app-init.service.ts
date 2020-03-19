import { Injectable } from '@angular/core';

@Injectable()
export class AppInitService {
  isInited = false;
  currentLang: string;

  constructor() { }
}
