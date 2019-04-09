import { ComponentFactoryResolver, Injectable, ViewContainerRef } from "@angular/core"
import { CreateProjectComponent } from "./create-project/create-project/create-project.component";
import { Project } from "../project/project";
import { MemberComponent } from "./create-project/member/member.component";
import { Observable } from "rxjs";
import { tap } from "rxjs/operators";

@Injectable()
export class SharedActionService {
  constructor(private resolver: ComponentFactoryResolver) {
  }

  createProjectComponent(container: ViewContainerRef): Observable<string> {
    let factory = this.resolver.resolveComponentFactory(CreateProjectComponent);
    let componentRef = container.createComponent(factory);
    return componentRef.instance.openCreateProjectModal()
      .pipe(tap(() => container.remove(container.indexOf(componentRef.hostView))));
  }

  createProjectMemberComponent(project: Project, container: ViewContainerRef): void {
    let factory = this.resolver.resolveComponentFactory(MemberComponent);
    let componentRef = container.createComponent(factory);
    componentRef.instance.openMemberModal(project)
      .subscribe(() => container.remove(container.indexOf(componentRef.hostView)));
  }
}
