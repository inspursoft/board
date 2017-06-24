import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule, Http } from '@angular/http';

@NgModule({
  exports:[
    BrowserAnimationsModule,
    BrowserModule,
    FormsModule,
    HttpModule
  ]
})
export class CoreModule {}