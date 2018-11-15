import { NgModule } from "@angular/core";
import { StorageComponent } from "./storage.component";
import { RouterModule, Routes } from "@angular/router";
import { RoutePV, RoutePvc } from "../shared/shared.const";
import { PvcComponent } from "./pvc/pvc.component";
import { SharedModule } from "../shared/shared.module";
import { CreatePvComponent } from "./pv/create-pv.component/create-pv.component";
import { StorageService } from "./storage.service";
import { AppTokenService } from "../app.init.service";
import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { HttpClientInterceptor } from "../shared/http-interceptor/http-client-interceptor";
import { MessageService } from "../shared/message-service/message.service";
import { MonitorsComponent } from "./pv/create-pv.component/monitors/monitors.component";
import { PvDetailComponent } from "./pv/pv-detail.compoent/pv-detail.component";
import { PvListComponent } from "./pv/pv-list.compoent/pv-list.component";

const routes: Routes = [
  {
    path: '', component: StorageComponent, children: [
      {path: RoutePV, component: PvListComponent},
      {path: RoutePvc, component: PvcComponent}
    ]
  },
];

@NgModule({
  imports: [
    SharedModule,
    RouterModule.forChild(routes),
  ],
  entryComponents: [
    CreatePvComponent,
    PvDetailComponent
  ],
  providers: [
    StorageService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpClientInterceptor,
      deps: [AppTokenService, MessageService],
      multi: true
    }
  ],
  declarations: [
    StorageComponent,
    PvListComponent,
    CreatePvComponent,
    PvcComponent,
    PvDetailComponent,
    MonitorsComponent
  ]
})
export class StorageModule {

}