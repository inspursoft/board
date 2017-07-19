export enum MESSAGE_TARGET { TOGGLE_PROJECT, DELETE_PROJECT };

export const DISMISS_INLINE_ALERT_INTERVAL: number = 2 * 1000;

export enum MESSAGE_TYPE {
  INVALID_USER, INTERNAL_ERROR
}

export const ROLES: {[key: number]: string} = { 1: 'Project Admin', 2: 'Developer', 3: 'Visitor' };