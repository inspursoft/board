import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AppInitService } from './app-init.service';
import { BoardService } from './board.service';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
  ],
  providers: [
    AppInitService,
    BoardService,
  ],
})
export class SharedServiceModule { }
