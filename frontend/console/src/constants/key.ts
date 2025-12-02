import { KeyStatus } from '../../../gen/src/keyhub/console/v1/common_pb';

export const KEY_STATUS_OPTIONS = [
  { value: KeyStatus.AVAILABLE, label: '利用可能' },
  { value: KeyStatus.IN_USE, label: '使用中' },
  { value: KeyStatus.LOST, label: '紛失' },
  { value: KeyStatus.DAMAGED, label: '破損' },
] as const;
