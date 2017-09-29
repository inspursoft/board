/*
 * Copyright (c) 2016 VMware, Inc. All Rights Reserved.
 * This software is released under MIT license.
 * The full license information can be found in LICENSE in the root directory of this project.
 */
import { ModuleWithProviders } from '@angular/core/src/metadata/ng_module';
import { Routes, RouterModule } from '@angular/router';

import { GlobalSearchComponent } from './global-search/global-search.component';
import { SignInComponent } from './account/sign-in/sign-in.component';
import { SignUpComponent } from './account/sign-up/sign-up.component';
import { MainContentComponent } from './main-content/main-content.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodeComponent } from './node/node.component';
import { ProjectComponent } from './project/project.component';
import { MemberComponent } from './project/member/member.component';
import { ImageListComponent } from './image/image-list/image-list.component';
import { ServiceComponent } from './service/service.component';
import { UserCenterComponent } from './user-center/user-center.component';
import { AuthGuard } from './shared/auth-guard.service';

export const ROUTES: Routes = [   
    { path: 'sign-in', component: SignInComponent },
    { path: 'sign-up', component: SignUpComponent },
    { path: '', component: MainContentComponent, 
        canActivate: [ AuthGuard ],
        children: [
        { path: 'search', component: GlobalSearchComponent },
        { path: 'dashboard', component: DashboardComponent },
        { path: 'nodes', component: NodeComponent },
        { path: 'projects',  
            children: [
                { path: '', component: ProjectComponent },
                { path: 'members', component: MemberComponent }
            ]
        },
        { path: 'images', component: ImageListComponent },
        { path: 'services', component: ServiceComponent },
        { path: 'user-center', component: UserCenterComponent }
    ]},
    { path: '', redirectTo: '/sign-in', pathMatch: 'full' },
    { path: '**', component: SignInComponent }
];

export const ROUTING: ModuleWithProviders = RouterModule.forRoot(ROUTES);
