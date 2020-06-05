import { Component, EventEmitter, Input, Output } from '@angular/core';
import { JobAffinityCardData, JobAffinityCardListView } from '../job.type';

@Component({
  selector: 'app-job-affinity-card-list',
  styleUrls: ['./job-affinity-card-list.component.css'],
  templateUrl: './job-affinity-card-list.component.html'
})
export class JobAffinityCardListComponent {
  @Input() title = '';
  @Input() description = '';
  @Input() acceptDrag = true;
  @Input() sourceList: Array<JobAffinityCardData>;
  @Input() listMinHeight = 100;
  @Input() listBorder = false;
  @Input() disabled = false;
  @Input() viewModel: JobAffinityCardListView = JobAffinityCardListView.aclvColumn;
  @Output() dropEvent: EventEmitter<string>;
  @Output() removeEvent: EventEmitter<JobAffinityCardData>;
  isDragActive = false;
  filterString = '';

  constructor() {
    this.dropEvent = new EventEmitter<string>();
    this.removeEvent = new EventEmitter<JobAffinityCardData>();
  }

  get filteredCardDataList(): Array<JobAffinityCardData> {
    return this.sourceList.filter((jobCard: JobAffinityCardData) =>
      this.filterString === '' ? true : jobCard.jobName.includes(this.filterString)
    );
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
      this.dropEvent.emit(dataKey);
    }
    event.preventDefault();
    event.stopPropagation();
  }

  removeAffinityCard(data: JobAffinityCardData) {
    this.removeEvent.emit(data);
  }

  filterExecute($event: KeyboardEvent) {
    this.filterString = ($event.target as HTMLInputElement).value
  }
}
