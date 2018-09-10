import { ComponentFactoryResolver, Injectable, ViewContainerRef } from "@angular/core"
import { CreateProjectComponent } from "./create-project/create-project/create-project.component";
import { Observable } from "rxjs/Observable";
import { Project } from "../project/project";
import "rxjs/add/operator/do"
import { MemberComponent } from "./create-project/member/member.component";

@Injectable()
export class SharedActionService {
  constructor(private resolver: ComponentFactoryResolver) {
  }

  createProjectComponent(container: ViewContainerRef): Observable<string> {
    let factory = this.resolver.resolveComponentFactory(CreateProjectComponent);
    let componentRef = container.createComponent(factory);
    return componentRef.instance.openModal().do(() => container.remove(container.indexOf(componentRef.hostView)));
  }

  createProjectMemberComponent(project: Project, container: ViewContainerRef): void {
    let factory = this.resolver.resolveComponentFactory(MemberComponent);
    let componentRef = container.createComponent(factory);
    componentRef.instance.openModal(project)
      .subscribe(() => container.remove(container.indexOf(componentRef.hostView)));
  }
}