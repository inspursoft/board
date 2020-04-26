import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AppInitService } from './app-init.service';
import { BoardService } from './board.service';
import { ConfigurationService } from './configuration.service';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
  ],
  providers: [
    AppInitService,
    BoardService,
    ConfigurationService,
  ],
})
export class SharedServiceModule { }
