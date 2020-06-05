import { Component, EventEmitter, Input, Output } from '@angular/core';
import { JobAffinityCardData } from '../job.type';
import { DragStatus } from '../../shared/shared.types';

@Component({
  selector: 'app-job-affinity-card',
  styleUrls: ['./job-affinity-card.component.css'],
  templateUrl: './job-affinity-card.component.html'
})
export class JobAffinityCardComponent {
  @Input() data: JobAffinityCardData;
  @Input() disabled = false;
  @Input() width = 0;
  @Output() removeFromListEvent: EventEmitter<JobAffinityCardData>;

  constructor() {
    this.removeFromListEvent = new EventEmitter<JobAffinityCardData>();
  }

  dragStartEvent(event: DragEvent) {
    this.data.status = DragStatus.dsStart;
    event.dataTransfer.setData('text', this.data.key);
  }

  dragEvent() {
    // this.data.status = DragStatus.dsDragIng;
  }

  get containerStyle() {
    return {
      'border-left': `5px green solid`,
      width: this.width === 0 ? `100%` : `${this.width}px`,
      cursor: this.data.status === DragStatus.dsEnd || this.disabled ? 'not-allowed' : 'pointer'
    };
  }

  removeFromList() {
    this.removeFromListEvent.emit(this.data);
  }
}
