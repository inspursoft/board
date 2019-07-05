import { NgModule } from '@angular/core';
import { HeaderComponent } from './header/header.component';
import { CheckItemExistingDirective } from './directives/check-item-existing.directive';
import { CheckItemIdenticalDirective } from './directives/check-item-identical.directive';
import { CheckItemPatternDirective } from './directives/check-item-pattern.directive';
import { ChangePasswordComponent } from './change-password/change-password.component';
import { AccountSettingComponent } from './account-setting/account-setting.component';
import { ValidateOnBlurDirective } from './directives/validate-onblur.directive';
import { CsDropdownComponent } from './cs-components-library/cs-dropdown/cs-dropdown.component';
import { CsSearchInput } from './cs-components-library/cs-search-input/cs-search-input.component';
import { CsInputComponent } from './cs-components-library/cs-input/cs-input.component';
import { CsInputArrayComponent } from './cs-components-library/cs-input-array/cs-input-array.component';
import { EnvironmentValueComponent } from './environment-value/environment-value.component';
import { SizePipe } from './pipes/size-pipe';
import { CsGuideComponent } from './cs-components-library/cs-guide/cs-guide.component';
import { CsProgressComponent } from './cs-components-library/cs-progress/cs-progress.component';
import { SafePipe } from './pipes/safe-pipe';
import { CreateProjectComponent } from './create-project/create-project/create-project.component';
import { MemberComponent } from './create-project/member/member.component';
import { CsHighlightComponent } from './cs-components-library/cs-highlight/cs-highlight.component';
import { CsInputCheckExistingDirective } from './cs-directives-library/cs-input-check-existing.directive';
import { CsInputCheckPatternDirective } from './cs-directives-library/cs-input-check-pattern.directive';
import { CsMenuItemUrlDirective, CsVerticalNavComponent } from './cs-components-library/cs-vertical-nav/cs-vertical-nav.component';
import 'inspurprism';
import { CsModalChildBaseSelector } from './cs-modal-base/cs-modal-child-base';
import { CreatePvcComponent } from './create-pvc/create-pvc.component';
import { CsInputDropdownComponent } from './cs-components-library/cs-input-dropdown/cs-input-dropdown.component';
import { CoreModule } from "../core/core.module";
import { LibCheckPatternExDirective } from "./lib-directives/input-check-pattern.directive";
import { LibCheckExistingExDirective } from "./lib-directives/input-check-existing.directive";

@NgModule({
  imports: [CoreModule],
  exports: [
    HeaderComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    LibCheckPatternExDirective,
    LibCheckExistingExDirective,
    CsDropdownComponent,
    CheckItemPatternDirective,
    ValidateOnBlurDirective,
    EnvironmentValueComponent,
    CsSearchInput,
    CsInputComponent,
    CsInputCheckPatternDirective,
    CsInputCheckExistingDirective,
    CsInputArrayComponent,
    CsModalChildBaseSelector,
    CsGuideComponent,
    CsProgressComponent,
    CsHighlightComponent,
    CsVerticalNavComponent,
    CsMenuItemUrlDirective,
    CsInputDropdownComponent,
    SizePipe,
    SafePipe
  ],
  declarations: [
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CheckItemPatternDirective,
    LibCheckPatternExDirective,
    LibCheckExistingExDirective,
    ValidateOnBlurDirective,
    CsDropdownComponent,
    CsSearchInput,
    CsInputComponent,
    CsInputCheckPatternDirective,
    CsInputCheckExistingDirective,
    CsInputArrayComponent,
    CsProgressComponent,
    CsModalChildBaseSelector,
    SizePipe,
    SafePipe,
    CsGuideComponent,
    CreateProjectComponent,
    CsInputDropdownComponent,
    CsHighlightComponent,
    CsVerticalNavComponent,
    CsMenuItemUrlDirective,
    AccountSettingComponent,
    EnvironmentValueComponent,
    MemberComponent,
    CreatePvcComponent,
    ChangePasswordComponent,
    HeaderComponent,
  ],
  entryComponents: [
    CreateProjectComponent,
    MemberComponent,
    CreatePvcComponent
  ]
})
export class SharedModule {
}
