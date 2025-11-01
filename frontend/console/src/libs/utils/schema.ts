import { z } from 'zod';
import { isBlankOrInvisible } from './string';
import { TenantType } from '../../../../gen/src/keyhub/console/v1/console_pb';

export const tenantnameValidation = z.preprocess(
  (val) => (typeof val === 'string' ? val.trim() : val),
  z
    .string({ message: 'テナント名を文字列で入力してください' })
    .nonempty({ message: 'テナント名を1文字以上入力してください' })
    .max(15, { message: 'テナント名は15文字以内で入力してください' })
    .refine((value: string) => !isBlankOrInvisible(value), {
      message: 'テナント名を1文字以上で入力してください',
    }),
);

export const descriptionValidation = z.preprocess(
  (val) => (typeof val == 'string' ? val.trim() : val),
  z
    .string({ message: 'テナント説明を文字列で入力してください' })
    .max(300, { message: 'テナント説明を300文字以内で入力してください' })
    .refine((value: string) => !isBlankOrInvisible(value), {
      message: 'テナント説明を1文字以上で入力してください',
    }),
);

export const tenanttypeValidation = z
  .number({ message: 'テナントタイプを選択してください' })
  .int({ message: 'テナントタイプは整数である必要があります' })
  .refine((val) => Object.values(TenantType).includes(val), {
    message: '有効なテナントタイプを選択してください',
  });
