import { Component, EventEmitter, Input, Output, QueryList, ViewChildren } from '@angular/core';
import { DragStatus } from '../../../shared/shared.types';
import { AffinityCardData, AffinityCardListView } from '../../service-step.component';
import { AffinityCardComponent } from '../affinity-card/affinity-card.component';

@Component({
  selector: 'app-affinity-card-list',
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
  @ViewChildren(AffinityCardComponent) cardComponentList: QueryList<AffinityCardComponent>;
  isDragActive = false;
  filterString = '';

  constructor() {
    this.affinityCardDataList = Array<AffinityCardData>();
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
      const dataKey = event.dataTransfer.getData('text');
      const card = this.sourceList.affinityCardDataList.find(value => value.key === dataKey);
      if (card) {
        const index = this.sourceList.affinityCardDataList.indexOf(card);
        card.status = DragStatus.dsEnd;
        this.affinityCardDataList.push(card);
        this.sourceList.affinityCardDataList.splice(index, 1);
      }
    }
    event.preventDefault();
    event.stopPropagation();
  }

  removeAffinityCard(data: AffinityCardData) {
    const card = this.affinityCardDataList.find(value => value.key === data.key);
    if (card) {
      const index = this.affinityCardDataList.indexOf(card);
      card.status = DragStatus.dsReady;
      this.sourceList.affinityCardDataList.push(card);
      this.affinityCardDataList.splice(index, 1);
    }
  }

  filterExecute($event: KeyboardEvent) {
    this.filterString = ($event.target as HTMLInputElement).value;
  }
}
