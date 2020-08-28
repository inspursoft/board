import { ResponsePaginationBase } from '../shared/ui-model/model-types';
import { User } from '../shared/shared.types';

export class UserPagination extends ResponsePaginationBase<User> {
  CreateOneItem(res: object): User {
    return new User(res);
  }

  ListKeyName(): string {
    return 'user_list';
  }
}
