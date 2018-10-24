import { Component, EventEmitter, Input, Output } from "@angular/core";
import { DragStatus } from "../../../shared/shared.types";
import { AffinityCardData } from "../../service-step.component";
import { SERVICE_STATUS } from "../../../shared/shared.const";

@Component({
  selector: 'affinity-card',
  styleUrls: ['./affinity-card.component.css'],
  templateUrl: './affinity-card.component.html'
})
export class AffinityCardComponent {
  @Input() data: AffinityCardData;
  @Input() disabled = false;
  @Input() width = 0;
  @Output() onRemoveFromList: EventEmitter<AffinityCardData>;
  @Output() onUnselected: EventEmitter<AffinityCardData>;

  constructor() {
    this.onRemoveFromList = new EventEmitter<AffinityCardData>();
    this.onUnselected = new EventEmitter<AffinityCardData>();
  }

  dragStartEvent(event: DragEvent) {
    this.data.status = DragStatus.dsStart;
    event.dataTransfer.setData("text", this.data.key);
  }

  dragEvent() {
    this.data.status = DragStatus.dsDragIng;
  }

  get containerStyle() {
    let getColor: () => string = () => {
      switch (this.data.serviceStatus) {
        case SERVICE_STATUS.PREPARING:
          return 'darkorange';
        case SERVICE_STATUS.RUNNING:
          return 'green';
        case SERVICE_STATUS.STOPPED:
          return 'gray';
        case SERVICE_STATUS.WARNING:
          return 'darkorange';
        case SERVICE_STATUS.DELETED:
          return 'red';
      }
    };
    return {
      'border-left': `5px ${getColor()} solid`,
      'width': this.width == 0 ? `100%` : `${this.width}px`,
      'cursor': this.data.status == DragStatus.dsEnd || this.disabled ? 'not-allowed' : 'pointer'
    }
  }

  get serviceStatus(): string {
    switch (this.data.serviceStatus) {
      case SERVICE_STATUS.PREPARING:
        return 'SERVICE.STATUS_PREPARING';
      case SERVICE_STATUS.RUNNING:
        return 'SERVICE.STATUS_RUNNING';
      case SERVICE_STATUS.STOPPED:
        return 'SERVICE.STATUS_STOPPED';
      case SERVICE_STATUS.WARNING:
        return 'SERVICE.STATUS_UNCOMPLETED';
      case SERVICE_STATUS.DELETED:
        return 'SERVICE.STATUS_DELETED';
    }
  }

  removeFromList() {
    this.onRemoveFromList.emit(this.data);
  }
}