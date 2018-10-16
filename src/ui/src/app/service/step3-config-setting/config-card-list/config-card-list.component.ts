import { AfterViewInit, ChangeDetectorRef, Component, EventEmitter, Input, Output, QueryList, ViewChildren } from "@angular/core";
import { DragStatus } from "../../../shared/shared.types";
import { ConfigCardData, ConfigCardModel } from "../../service-step.component";
import { ConfigCardComponent } from "../config-card/config-card.component";

@Component({
  selector: 'config-card-list',
  styleUrls: ['./config-card-list.component.css'],
  templateUrl: './config-card-list.component.html'
})
export class ConfigCardListComponent implements AfterViewInit {
  @Input() title = '';
  @Input() description = '';
  @Input() acceptDrag = true;
  @Input() sourceList: ConfigCardListComponent;
  @Input() ConfigCardDataList: Array<ConfigCardData>;
  @Input() model: ConfigCardModel = ConfigCardModel.cmDefault;
  @Input() cardMinHeight = 60;
  @Input() cardMinWidth = 100;
  @Input() listMinHeight = 100;
  @Input() listBorder = false;
  @Input() disabled = false;
  @Input() selected: ConfigCardData;
  @Output() selectedChange: EventEmitter<ConfigCardData>;
  @Output() onDrop: EventEmitter<ConfigCardData>;
  @ViewChildren(ConfigCardComponent) cardComponentList: QueryList<ConfigCardComponent>;
  isDragActive = false;

  constructor(private changeDetectorRef: ChangeDetectorRef) {
    this.ConfigCardDataList = Array<ConfigCardData>();
    this.onDrop = new EventEmitter<ConfigCardData>();
    this.selectedChange = new EventEmitter<ConfigCardData>();
  }

  ngAfterViewInit() {
    if (this.model == ConfigCardModel.cmSelect && this.selected) {
      this.cardComponentList.forEach(component => {
        if (component.data.key === this.selected.key) {
          component.isSelected = true;
          this.changeDetectorRef.detectChanges();
        }
      });
    }
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
      let dataKey = event.dataTransfer.getData('data');
      let configCardData = this.sourceList.ConfigCardDataList.find(value => value.key === dataKey);
      if (configCardData) {
        let index = this.sourceList.ConfigCardDataList.indexOf(configCardData);
        configCardData.status = DragStatus.dsEnd;
        this.ConfigCardDataList.push(configCardData);
        this.sourceList.ConfigCardDataList.splice(index, 1);
        this.onDrop.emit(configCardData);
      }
    }
  }

  removeContainerCard(data: ConfigCardData) {
    let configCardData = this.ConfigCardDataList.find(value => value.key === data.key);
    if (configCardData) {
      let index = this.ConfigCardDataList.indexOf(configCardData);
      configCardData.externalInfo = '';
      configCardData.status = DragStatus.dsReady;
      this.sourceList.ConfigCardDataList.push(configCardData);
      this.ConfigCardDataList.splice(index, 1);
    }
  }

  selectedCard(data: ConfigCardData) {
    this.cardComponentList.forEach(component => {
      if (component.data.key !== data.key) {
        component.isSelected = false;
      }
    });
    this.selected = data;
    this.selectedChange.emit(data);
  }

  unselectedCard(data: ConfigCardData) {
    this.selected = undefined;
    this.selectedChange.emit(undefined);
  }
}