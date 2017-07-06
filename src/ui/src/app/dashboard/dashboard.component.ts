import {AfterViewInit, OnInit, Component, OnDestroy} from '@angular/core';
import {DatePipe} from "@angular/common";
import {Assist} from "./dashboard-assist"
import {scaleOption} from "app/dashboard/time-range-scale.component/time-range-scale.component";
import {DashboardService, ServiceListModel} from "app/dashboard/dashboard.service";

@Component({
	selector: 'dashboard',
	templateUrl: 'dashboard.component.html',
	styleUrls: ['dashboard.component.css']
})
	
export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy {
	scaleOptions: Array<scaleOption> = [
		{"id": 1, "description": "DASHBOARD.MIN"},
		{"id": 2, "description": "DASHBOARD.HR"},
		{"id": 3, "description": "DASHBOARD.DAY"},
		{ "id": 4, "description": "DASHBOARD.MTH" }];
	
	podCount: number = 7;
	containerCount: number = 32;
	memoryPercent: string = '70%';
	cpuPercent: string = '40%';
	usageVolume: string = '3T';
	totalVolume: string = '10T';

	serviceBtnValue: string;
	nodeBtnValue: string;
	storageBtnValue: string;
	serviceList: Array<ServiceListModel>;
	serviceOptions: object = {};
	nodeOptions: object = {};
	storageOptions: object = {};

	constructor(private service: DashboardService) {
		let s: Set<object> = new Set<object>();
		const a = { name: "123" };
		const b = { name: "123" };
		s.add(a).add(b).add(a);
    s.forEach(value => console.log(value));
	}

	ngOnInit() {
		this.serviceList = this.service.getServiceList();
		this.serviceBtnValue = this.serviceList[0].serviceName;
		this.nodeBtnValue = this.serviceList[0].serviceName;
		this.storageBtnValue = this.serviceList[0].serviceName;
	}

	ngOnDestroy() {

	}

	scaleChange(data: scaleOption) {

	}

	ngAfterViewInit() {
		this.serviceOptions = Assist.getBaseOptions();
		this.nodeOptions = Assist.getBaseOptions();
		this.storageOptions = Assist.getBaseOptions();

		this.serviceOptions["tooltip"] = Assist.getTooltip("pods", "containers");
		this.nodeOptions["tooltip"] = Assist.getTooltip("CPU", "Memory");
		this.storageOptions["tooltip"] = Assist.getTooltip("", "Total");

		this.serviceOptions["series"] = [Assist.getBaseSeries(), Assist.getBaseSeries()];
		let serviceData = this.service.getServiceData(0, 1);
		this.serviceOptions["series"][0]["data"] = serviceData[0];
		this.serviceOptions["series"][1]["data"] = serviceData[1];
	}

	get serviceIcon(): string {
		return '../../images/service_icon.png';
	}

	get nodeIcon(): string {
		return '../../images/node_icon.png';
	}

	get storageIcon(): string {
		return '../../images/storage_icon.png';
	}
}