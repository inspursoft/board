/**
 * Created by liyanq on 8/28/17.
 */

import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-service-wizard',
  templateUrl: './service-wizard.component.html',
  styleUrls: ['./service-wizard.component.css']
})
export class ServiceWizardComponent {
  @Input() curStep = 0;
}

