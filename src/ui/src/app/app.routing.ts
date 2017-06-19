/*
 * Copyright (c) 2016 VMware, Inc. All Rights Reserved.
 * This software is released under MIT license.
 * The full license information can be found in LICENSE in the root directory of this project.
 */
import { ModuleWithProviders } from '@angular/core/src/metadata/ng_module';
import { Routes, RouterModule } from '@angular/router';

import { DashboardComponent } from './dashboard/dashboard.component';
import { NodeComponent } from './node/node.component';
import { ProjectComponent } from './project/project.component';
import { MemberComponent } from './project/member/member.component';
import { ImageComponent } from './image/image.component';
import { ServiceComponent } from './service/service.component';
import { AdminOptionComponent } from './admin-option/admin-option.component';
import { ProfileComponent } from './profile/profile.component';

export const ROUTES: Routes = [
    { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
    { path: 'dashboard', component: DashboardComponent },
    { path: 'nodes', component: NodeComponent },
    { path: 'projects',  
          children: [
            { path: '', component: ProjectComponent },
            { path: 'members', component: MemberComponent }
        ] 
    },
    { path: 'images', component: ImageComponent },
    { path: 'services', component: ServiceComponent },
    { path: 'profiles', component: ProfileComponent },
    { path: 'admin-options', component: AdminOptionComponent }
];

export const ROUTING: ModuleWithProviders = RouterModule.forRoot(ROUTES);
