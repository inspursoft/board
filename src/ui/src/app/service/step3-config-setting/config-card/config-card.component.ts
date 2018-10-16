import { Component, EventEmitter, Input, Output } from "@angular/core";
import { DragStatus } from "../../../shared/shared.types";
import { ConfigCardData, ConfigCardModel } from "../../service-step.component";

@Component({
  selector: 'config-card',
  styleUrls: ['./config-card.component.css'],
  templateUrl: './config-card.component.html'
})
export class ConfigCardComponent {
  @Input() data: ConfigCardData;
  @Input() minWidth = 100;
  @Input() minHeight = 60;
  @Input() model: ConfigCardModel = ConfigCardModel.cmDefault;
  @Input() disabled = false;
  @Output() onRemoveFromList: EventEmitter<ConfigCardData>;
  @Output() onSelected: EventEmitter<ConfigCardData>;
  @Output() onUnselected: EventEmitter<ConfigCardData>;
  public isSelected = false;

  constructor() {
    this.onRemoveFromList = new EventEmitter<ConfigCardData>();
    this.onSelected = new EventEmitter<ConfigCardData>();
    this.onUnselected = new EventEmitter<ConfigCardData>();
  }

  dragStartEvent(event: DragEvent) {
    this.data.status = DragStatus.dsStart;
    event.dataTransfer.setData("data", this.data.key);
  }

  dragEvent() {
    this.data.status = DragStatus.dsDragIng;
  }

  get containerStyle() {
    return {
      'min-width': `${this.minWidth}px`,
      'min-height': `${this.minHeight}px`,
      'cursor': this.data.status == DragStatus.dsEnd || this.disabled ? 'not-allowed' : 'pointer'
    }
  }

  removeFromList() {
    this.onRemoveFromList.emit(this.data);
  }

  onSelect() {
    this.isSelected = !this.isSelected;
    if (this.isSelected) {
      this.onSelected.emit(this.data);
    } else {
      this.onUnselected.emit(this.data);
    }
  }

}