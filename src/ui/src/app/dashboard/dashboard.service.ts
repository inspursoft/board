import { Injectable } from '@angular/core';
import { Http } from "@angular/http"
import "rxjs/add/operator/map";
import 'rxjs/add/operator/toPromise';

export interface ServiceListModel {
	readonly id: number;
	readonly serviceName: string;
}

export interface ServiceDataModel {
	readonly date: Date;
	readonly value: number;
}

export interface NodeDataModel {
	readonly date: Date;
	readonly value: number;
}

export interface StorageDataModel {
	readonly date: Date;
	readonly value: number;
}

@Injectable()
export class DashboardService {
	static baseDate: Date = new Date();

	getOneStepTime(dateScaleId: number): number {
		switch (dateScaleId) {
			case 1:
				return 5 * 1000;
			case 2:
				return 60 * 1000;
			case 3:
				return 24 * 60 * 1000;
			case 4:
				return 12 * 24 * 60 * 1000;

			default:
				return 1;
		}
	}

	private getSimulateData(serviceID: number): number {
		switch (serviceID) {
			case 0:
				return 130 + Math.round(Math.random() * 50);
			case 1:
				return 30 + Math.round(Math.random() * 10);
			case 2:
				return 20 + Math.round(Math.random() * 10);
			case 3:
				return 50 + Math.round(Math.random() * 10);
			case 4:
				return 30 + Math.round(Math.random() * 20);
		}
	}

	private getSimulateDate(dateScaleId: number): Date {
		DashboardService.baseDate.setTime(DashboardService.baseDate.getTime()
			+ this.getOneStepTime(dateScaleId));
		return new Date(DashboardService.baseDate.getTime());
	}

	/**
	 * getServiceData
	 * @param serviceID
	 * @param dateScaleId
	 *node:dateScaleId=>1:1min;2:1hr;3:1day;4:1mth
	 */
	getServiceData(serviceID: number, dateScaleId: number): Map<number, ServiceDataModel[]> {
		if (!dateScaleId || dateScaleId < 1 || dateScaleId > 4) return null;
		let r: Map<number, ServiceDataModel[]> = new Map<number, ServiceDataModel[]>();
		r[0] = new Array<ServiceDataModel>(0);
		r[1] = new Array<ServiceDataModel>(0);
		for (let i = 0; i < 11; i++) {
			let date: Date = this.getSimulateDate(dateScaleId)
			let arrBuf1 = [date, this.getSimulateData(serviceID)];
			let arrBuf2 = [date, this.getSimulateData(serviceID)];
			r[0].push(arrBuf1);
			r[1].push(arrBuf2);
		}
		return r;
	}

	getServiceList(): ServiceListModel[] {
		return Array.from([
			{ "id": 0, "serviceName": "total" },
			{ "id": 1, "serviceName": "myService1" },
			{ "id": 2, "serviceName": "myService2" },
			{ "id": 3, "serviceName": "我的服务3" },
			{ "id": 4, "serviceName": "我的服务4" }
		]);
		// return this.http.get("./someData/ServiceList.JSON")
		//   .toPromise()
		//   .then(res=> Array.from(res.json()))
		//   .catch(this.handleError);
	}


	private handleError(error: Response | any) {
		let errMsg: string;
		if (error instanceof Response) {
			const body = error.json() || '';
			const err = body["error"] || JSON.stringify(body);
			errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
		} else {
			errMsg = error.message ? error.message : error.toString();
		}
		console.error(errMsg);
		return Promise.reject(errMsg);
	}

}
