import { OnInit, AfterViewInit, Component, OnDestroy } from '@angular/core';
import { Assist } from "./dashboard-assist"
import { scaleOption } from "app/dashboard/time-range-scale.component/time-range-scale.component";
import { DashboardService, ServiceListModel, LineDataModel } from "app/dashboard/dashboard.service";
import { TranslateService } from "@ngx-translate/core";
import { Subscription } from "rxjs/Subscription";

@Component( {
  selector: 'dashboard',
  templateUrl: 'dashboard.component.html',
  styleUrls: [ 'dashboard.component.css' ]
} )

export class DashboardComponent implements OnInit, AfterViewInit, OnDestroy {
  _intervalRead: any;
  _onLangChangeSubscription: Subscription;
  scaleOptions: Array<scaleOption> = [
    {"id": 1, "description": "DASHBOARD.MIN", "value": "second"},
    {"id": 2, "description": "DASHBOARD.HR", "value": "minute"},
    {"id": 3, "description": "DASHBOARD.DAY", "value": "hour"},
    {"id": 4, "description": "DASHBOARD.MTH", "value": "day"} ];

  podCount: number = 0;
  containerCount: number = 0;

  memoryPercent: string = '70%';
  cpuPercent: string = '40%';
  usageVolume: string = '3T';
  totalVolume: string = '10T';

  serviceBtnValue: string;
  serviceList: Array<ServiceListModel>;
  serviceQuery: { model: ServiceListModel, scale: scaleOption, count: number };
  serviceOptions: object = {};
  serviceAlready: boolean = false;
  serviceNoData: boolean = false;
  serviceData: Map<string, LineDataModel[]>;

  nodeBtnValue: string;
  nodeOptions: object = {};

  storageBtnValue: string;
  storageOptions: object = {};

  constructor( private service: DashboardService,
               private translateService: TranslateService ) {
  }

  serviceScaleChange( data: scaleOption ) {
    this.serviceQuery.scale = data;
    this.serviceAlready = false;
    this.refreshServiceData();
  }

  serviceDropDownChange( service: ServiceListModel ) {
    this.serviceBtnValue = service.service_name;
    this.serviceAlready = false;
    this.serviceQuery.model = service;
    this.refreshServiceData();
  }

  translateForService( isResolve: boolean ): void {
    this.translateService.get( [ "DASHBOARD.CONTAINERS", "DASHBOARD.PODS" ] )
      .subscribe( res => {
        let podsTranslate: string = res[ "DASHBOARD.PODS" ];
        let containersTranslate: string = res[ "DASHBOARD.CONTAINERS" ];
        this.serviceOptions = Assist.getServiceOptions();
        this.serviceOptions[ "tooltip" ] = Assist.getTooltip( podsTranslate, containersTranslate );
        this.serviceOptions[ "series" ] = [ Assist.getBaseSeries(), Assist.getBaseSeries() ];
        this.serviceOptions[ "series" ][ 0 ][ "name" ] = podsTranslate;
        this.serviceOptions[ "series" ][ 1 ][ "name" ] = containersTranslate;
        this.serviceOptions[ "legend" ] = {data: [ podsTranslate, containersTranslate ], x: "left"};
        if (isResolve && this.serviceData && this.serviceData[ "pods" ].length > 0) {
          this.serviceOptions[ "series" ][ 0 ][ "data" ] = this.serviceData[ "pods" ];
          this.serviceOptions[ "series" ][ 1 ][ "data" ] = this.serviceData[ "container" ];
          this.podCount = this.serviceData[ "pods" ][ this.serviceData[ "pods" ].length - 1 ][ 1 ] | 0;
          this.containerCount = this.serviceData[ "container" ][ this.serviceData[ "container" ].length - 1 ][ 1 ] | 0;
          this.serviceNoData = false;
        } else {
          this.serviceOptions[ "series" ][ 0 ][ "data" ] = [ [ new Date(), 0 ] ];
          this.serviceOptions[ "series" ][ 1 ][ "data" ] = [ [ new Date(), 0 ] ];
          this.podCount = 0;
          this.containerCount = 0;
          this.serviceNoData = true;
        }
        this.serviceAlready = true;
      } );
  }

  refreshServiceData() {
    this.service.getServiceData(
      this.serviceQuery.count,
      this.serviceQuery.scale.value,
      this.serviceQuery.model.service_name ).then( res => {
      this.serviceData = res;
      this.translateForService( true );
    } ).catch( err => {
      if (this.serviceData) {
        this.serviceData.clear();
      }
      this.translateForService( false );
    } );
  }

  ngOnInit() {
    this._intervalRead = setInterval( () => {
      // this.refreshServiceData();
    }, 10000 );
    this._onLangChangeSubscription = this.translateService.onLangChange.subscribe( () => {
      this.refreshServiceData();
    } )
  }

  ngOnDestroy() {
    clearInterval( this._intervalRead );
    if (this._onLangChangeSubscription) {
      this._onLangChangeSubscription.unsubscribe();
    }
  }

  ngAfterViewInit() {
    this.service.getServiceList().then( res => {
      this.serviceList = res;
      this.serviceBtnValue = this.serviceList[ 0 ].service_name;
      this.nodeBtnValue = this.serviceList[ 0 ].service_name;
      this.storageBtnValue = this.serviceList[ 0 ].service_name;

      this.serviceQuery = Object.create( {
        count: 300,
        model: this.serviceList[ 0 ],
        scale: this.scaleOptions[ 0 ]
      } );
      this.refreshServiceData();
    } );

    let serviceSimulateData = this.service.getBySimulateData( 0, 1 );

    this.nodeOptions = Assist.getBaseOptions();
    this.nodeOptions[ "tooltip" ] = Assist.getTooltip( "CPU", "Memory" );
    this.nodeOptions[ "series" ] = [ Assist.getBaseSeries(), Assist.getBaseSeries() ];
    this.nodeOptions[ "series" ][ 0 ][ "data" ] = serviceSimulateData[ 0 ];
    this.nodeOptions[ "series" ][ 1 ][ "data" ] = serviceSimulateData[ 1 ];

    this.storageOptions = Assist.getBaseOptions();
    this.storageOptions[ "tooltip" ] = Assist.getTooltip( "", "Total" );
    this.storageOptions[ "series" ] = [ Assist.getBaseSeries(), Assist.getBaseSeries() ];
    this.storageOptions[ "series" ][ 0 ][ "data" ] = serviceSimulateData[ 0 ];
    this.storageOptions[ "series" ][ 1 ][ "data" ] = serviceSimulateData[ 1 ];
  }
}