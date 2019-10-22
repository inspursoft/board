import { Component, EventEmitter, Input, Output, QueryList, ViewChildren } from "@angular/core";
import { DragStatus } from "../../../shared/shared.types";
import { AffinityCardData, AffinityCardListView } from "../../service-step.component";
import { AffinityCardComponent } from "../affinity-card/affinity-card.component";

@Component({
  selector: 'affinity-card-list',
  styleUrls: ['./affinity-card-list.component.css'],
  templateUrl: './affinity-card-list.component.html'
})
export class AffinityCardListComponent {
  @Input() title = '';
  @Input() description = '';
  @Input() acceptDrag = true;
  @Input() sourceList: AffinityCardListComponent;
  @Input() affinityCardDataList: Array<AffinityCardData>;
  @Input() listMinHeight = 100;
  @Input() listBorder = false;
  @Input() disabled = false;
  @Input() cardWidth = 0;
  @Input() viewModel: AffinityCardListView = AffinityCardListView.aclvColumn;
  @Output() selectedChange: EventEmitter<AffinityCardData>;
  @Output() onDrop: EventEmitter<AffinityCardData>;
  @ViewChildren(AffinityCardComponent) cardComponentList: QueryList<AffinityCardComponent>;
  isDragActive = false;
  filterString = '';

  constructor() {
    this.affinityCardDataList = Array<AffinityCardData>();
    this.onDrop = new EventEmitter<AffinityCardData>();
    this.selectedChange = new EventEmitter<AffinityCardData>();
  }

  get filteredCardDataList(): Array<AffinityCardData> {
    return this.affinityCardDataList.filter(value => this.filterString === '' ? true : value.serviceName.includes(this.filterString));
  }

  dragOver(event: DragEvent) {
    if (this.acceptDrag && !this.disabled) {
      this.isDragActive = true;
      event.preventDefault();
      event.stopPropagation();
    }
  }

  dragLeave() {
    this.isDragActive = false;
  }

  drop(event: DragEvent) {
    if (this.acceptDrag && !this.disabled) {
      this.isDragActive = false;
      let dataKey = event.dataTransfer.getData('text');
      let card = this.sourceList.affinityCardDataList.find(value => value.key === dataKey);
      if (card) {
        let index = this.sourceList.affinityCardDataList.indexOf(card);
        card.status = DragStatus.dsEnd;
        this.affinityCardDataList.push(card);
        this.sourceList.affinityCardDataList.splice(index, 1);
        this.onDrop.emit(card);
      }
    }
    event.preventDefault();
    event.stopPropagation();
  }

  removeAffinityCard(data: AffinityCardData) {
    let card = this.affinityCardDataList.find(value => value.key === data.key);
    if (card) {
      let index = this.affinityCardDataList.indexOf(card);
      card.status = DragStatus.dsReady;
      this.sourceList.affinityCardDataList.push(card);
      this.affinityCardDataList.splice(index, 1);
    }
  }

  filterExecute($event: KeyboardEvent) {
    this.filterString = ($event.target as HTMLInputElement).value
  }
}
