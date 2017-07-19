import {Injectable, OnInit, OnDestroy} from '@angular/core';
import {Subject} from 'rxjs/Subject';
import {Observable} from 'rxjs/Observable';

type Global_Alert_Type = "alert-success" | "alert-info" | "alert-warning" | "alert-danger";

@Injectable()
export class MessageService implements OnInit, OnDestroy {
	_Global_Message: string;
	_Alert_Type: Global_Alert_Type;
	_Global_Interval: any;
	_Global_Interval_Seed: number = 0;
	Show_Global_Alert: boolean = false;

	get get_Global_Message() {
		return this._Global_Message;
	}

	get get_AlertType() {
		return this._Alert_Type;
	}

	constructor() {
		this._Global_Interval = setInterval(() => {
			if (this.Show_Global_Alert && this._Global_Interval_Seed == 0) {
				this.Show_Global_Alert = false;
				this._Global_Interval_Seed = 2;
			}
			if (this._Global_Interval_Seed > 0) {
				this._Global_Interval_Seed--;
			}
		}, 1000);
	}

	ngOnInit() {
	}

	ngOnDestroy() {
		clearInterval(this._Global_Interval);
	}

	set globalInfo(message) {
		this._Alert_Type = "alert-info";
		this._Global_Message = message;
		this._Global_Interval_Seed = 2;
		this.Show_Global_Alert = true;
	}

	set globalWarning(message) {
		this._Alert_Type = "alert-warning";
		this._Global_Message = message;
		this._Global_Interval_Seed = 2;
		this.Show_Global_Alert = true;
	}

	messageAnnouncedSource: Subject<any> = new Subject<any>();
	messageAnnounced$: Observable<any> = this.messageAnnouncedSource.asObservable();

	messageConfirmedSource: Subject<any> = new Subject<any>();
	messageConfirmed$: Observable<any> = this.messageConfirmedSource.asObservable();

	announceMessage(message: any) {
		this.messageAnnouncedSource.next(message);
	}

	confirmMessage(message: any) {
		this.messageConfirmedSource.next(message);
	}

}