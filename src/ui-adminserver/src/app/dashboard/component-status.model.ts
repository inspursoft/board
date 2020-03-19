import { ResponseBase, HttpBind } from '../shared/shared.type';

export class ComponentStatus extends ResponseBase {
  @HttpBind('id') containerId: string;
  @HttpBind('image') image: string;
  @HttpBind('created_at') created: string;
  @HttpBind('status') status: string;
  @HttpBind('ports') ports: string;
  @HttpBind('name') name: string;

  @HttpBind('cpu_perc') cpuRate: string;
  @HttpBind('mem_usage') memUsage: string;
  @HttpBind('mem_perc') memRate: string;
  @HttpBind('net_io') netIO: string;
  @HttpBind('block_io') blockIO: string;
  @HttpBind('pids') pid: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      const id = String(Math.round(Math.random() * 1000000000000));
      this.containerId = id.length > 11 ? id : (Array(12).join('0') + id).slice(-12);
      this.image = 'test:dev';
      this.created = '8 days ago';
      this.status = 'Up 8 days';
      this.ports = '0.0.0.0:80->80/tcp';
      this.name = 'board_module_' + this.containerId;

      this.cpuRate = '0.00%';
      this.memUsage = '0.1MiB';
      this.memRate = String((Math.random() * 100).toFixed(2));
      this.netIO = '1.0MB/1.0MB';
      this.blockIO = '0B/1B';
      this.pid = '5';
    }
  }
}
