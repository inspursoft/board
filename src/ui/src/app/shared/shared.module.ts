import { NgModule } from '@angular/core';
import { CoreModule } from '../core/core.module';
import { NgxEchartsModule } from 'ngx-echarts';
import { HeaderComponent } from './header/header.component';
import { ConfirmationDialogComponent } from './confirmation-dialog/confirmation-dialog.component';
import { InlineAlertComponent } from './inline-alert/inline-alert.component';
import { GlobalMessageComponent } from './global-message/global-message.component';
import { CheckItemExistingDirective } from './directives/check-item-existing.directive';
import { CheckItemIdenticalDirective } from './directives/check-item-identical.directive';
import { MessageService } from './message-service/message.service';
import { AuthGuard, ServiceGuard } from './auth-guard.service';
import { CheckItemPatternDirective } from "./directives/check-item-pattern.directive";
import { ChangePasswordComponent } from "./change-password/change-password.component";
import { AccountSettingComponent } from "./account-setting/account-setting.component";
import { ValidateOnBlurDirective } from './directives/validate-onblur.directive';
import { CsDropdownComponent } from "./cs-components-library/cs-dropdown/cs-dropdown.component";
import { WebsocketService } from './websocket-service/websocket.service';
import { CsSearchInput } from "./cs-components-library/cs-search-input/cs-search-input.component";
import { CsInputComponent } from "./cs-components-library/cs-input/cs-input.component";
import { CsInputArrayComponent } from "./cs-components-library/cs-input-array/cs-input-array.component";
import { CreateImageComponent } from "../image/image-create/image-create.component";
import { EnvironmentValueComponent } from "./environment-value/environment-value.component";
import { SizePipe } from "./pipes/size-pipe";
import { CsGuideComponent } from "./cs-components-library/cs-guide/cs-guide.component";
import { CsProgressComponent } from "./cs-components-library/cs-progress/cs-progress.component";
import { SafePipe } from "./pipes/safe-pipe";
import {
  CsSyntaxHighlighterComponent,
  CsSyntaxHighlighterDirective
} from './cs-components-library/cs-syntax-highlighter/cs-syntax-highlighter.component'

@NgModule({
  imports: [
    CoreModule,
    NgxEchartsModule
  ],
  exports: [
    CoreModule,
    NgxEchartsModule,
    ConfirmationDialogComponent,
    HeaderComponent,
    InlineAlertComponent,
    GlobalMessageComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CsDropdownComponent,
    CheckItemPatternDirective,
    ValidateOnBlurDirective,
    CreateImageComponent,
    EnvironmentValueComponent,
    CsSearchInput,
    CsInputComponent,
    CsInputArrayComponent,
    CsGuideComponent,
    CsProgressComponent,
    SizePipe,
    SafePipe,
    CsSyntaxHighlighterComponent
  ],
  declarations: [
    ConfirmationDialogComponent,
    HeaderComponent,
    InlineAlertComponent,
    GlobalMessageComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CheckItemPatternDirective,
    ValidateOnBlurDirective,
    ChangePasswordComponent,
    CsDropdownComponent,
    AccountSettingComponent,
    CreateImageComponent,
    EnvironmentValueComponent,
    CsSearchInput,
    CsInputComponent,
    CsInputArrayComponent,
    CsProgressComponent,
    SizePipe,
    SafePipe,
    CsGuideComponent,
    CsSyntaxHighlighterDirective,
    CsSyntaxHighlighterComponent,
  ],
  providers: [
    AuthGuard,
    ServiceGuard,
    MessageService,
    WebsocketService
  ]
})
export class SharedModule {}
