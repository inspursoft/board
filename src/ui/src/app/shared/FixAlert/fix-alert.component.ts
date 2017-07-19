import {Component} from "@angular/core"
import {MessageService} from "app/shared/service/message.service";

@Component({
	selector: "fix-alert",
	templateUrl: "./fix-alert.component.html",
	styleUrls: ["./fix-alert.component.css"]
})

export class FixAlert {
	constructor(private messageService: MessageService) {
	};

	get Message() {
		return this.messageService.get_Global_Message;
	}

	get AlertType() {
		return this.messageService.get_AlertType;
	}

	get IsShow() {
		return this.messageService.Show_Global_Alert;
	}

	set IsShow(value:boolean){
		this.messageService.Show_Global_Alert = value;
	}
}