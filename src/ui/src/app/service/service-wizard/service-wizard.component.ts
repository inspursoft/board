/**
 * Created by liyanq on 8/28/17.
 */

import { Component, Input } from "@angular/core"

@Component({
  selector:"service-wizard",
  templateUrl:"./service-wizard.component.html",
  styleUrls:["./service-wizard.component.css"]
})
export class ServiceWizardComponent{
  @Input("step") curStep:number = 0;
}