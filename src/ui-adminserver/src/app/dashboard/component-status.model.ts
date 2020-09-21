import { ResponseBase, HttpBind } from '../shared/shared.type';

export class ComponentStatus extends ResponseBase {
  @HttpBind('id') id: string;
  @HttpBind('image') image: string;
  @HttpBind('created_at') created_at: string;
  @HttpBind('status') status: string;
  @HttpBind('ports') ports: string;
  @HttpBind('name') name: string;

  @HttpBind('cpu_perc') cpu_perc: string;
  @HttpBind('mem_usage') mem_usage: string;
  @HttpBind('mem_perc') mem_perc: string;
  @HttpBind('net_io') net_io: string;
  @HttpBind('block_io') block_io: string;
  @HttpBind('pids') pids: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.id = '';
      this.image = 'test:dev';
      this.created_at = '8 days ago';
      this.status = 'Up 8 days';
      this.ports = '0.0.0.0:80->80/tcp';
      this.name = 'board_module_' + this.id;

      this.cpu_perc = '0.00%';
      this.mem_usage = '0.1MiB';
      this.mem_perc = String((Math.random() * 100).toFixed(2));
      this.net_io = '1.0MB/1.0MB';
      this.block_io = '0B/1B';
      this.pids = '5';
    }
  }
}
