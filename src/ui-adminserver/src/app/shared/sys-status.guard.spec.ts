import { TestBed, async, inject } from '@angular/core/testing';

import { SysStatusGuard } from './sys-status.guard';

describe('SysStatusGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [SysStatusGuard]
    });
  });

  it('should ...', inject([SysStatusGuard], (guard: SysStatusGuard) => {
    expect(guard).toBeTruthy();
  }));
});
