import { ResponsePaginationBase } from '../shared/ui-model/model-types';
import { SharedProject } from '../shared/shared.types';

export class PaginationProject extends ResponsePaginationBase<SharedProject> {
  ListKeyName(): string {
    return 'project_list';
  }

  CreateOneItem(res: object): SharedProject {
    return new SharedProject(res);
  }

}
