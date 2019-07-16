import { Component, EventEmitter, Input, Output, QueryList, ViewChildren } from "@angular/core";
import { JobAffinityCardData, JobAffinityCardListView } from "../job.type";

@Component({
  selector: 'job-affinity-card-list',
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
  @Output() onDrop: EventEmitter<string>;
  @Output() onRemove: EventEmitter<JobAffinityCardData>;
  isDragActive = false;
  filterString = '';

  constructor() {
    this.onDrop = new EventEmitter<string>();
    this.onRemove = new EventEmitter<JobAffinityCardData>();
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
    }
  }

  dragLeave() {
    this.isDragActive = false;
  }

  drop(event: DragEvent) {
    if (this.acceptDrag && !this.disabled) {
      this.isDragActive = false;
      const dataKey = event.dataTransfer.getData('text');
      this.onDrop.emit(dataKey);
    }
  }

  removeAffinityCard(data: JobAffinityCardData) {
    this.onRemove.emit(data);
  }

  filterExecute($event: KeyboardEvent) {
    this.filterString = ($event.target as HTMLInputElement).value
  }
}
