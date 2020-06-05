import { Component, Input, OnInit } from '@angular/core';
import { PaginationJob } from '../job.type';

@Component({
  selector: 'app-job-create-guide',
  templateUrl: './job-create-guide.component.html',
  styleUrls: ['./job-create-guide.component.css']
})
export class JobCreateGuideComponent implements OnInit {
  @Input() paginationJobs: PaginationJob;

  constructor() {
  }

  ngOnInit() {
  }

}
