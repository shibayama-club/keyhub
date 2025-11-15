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

export const joinCodeValidation = z.preprocess(
  (val) => (typeof val === 'string' ? val.trim() : val),
  z
    .string({ message: '参加コードを文字列で入力してください' })
    .min(6, { message: '参加コードは6文字以上で入力してください' })
    .max(20, { message: '参加コードは20文字以内で入力してください' })
    .regex(/^[a-zA-Z0-9]+$/, { message: '参加コードは英数字のみで入力してください' }),
);

export const joinCodeExpiryValidation = z.date({ message: '有効期限を日時で指定してください' }).optional();

export const joinCodeMaxUseValidation = z
  .number({ message: '最大使用回数を数値で入力してください' })
  .int({ message: '最大使用回数は整数である必要があります' })
  .min(0, { message: '最大使用回数は0以上である必要があります' })
  .optional();

// テナント作成・編集フォームのスキーマ
export const tenantSchema = z.object({
  name: tenantnameValidation,
  description: descriptionValidation.optional(),
  tenantType: tenanttypeValidation,
  joinCode: joinCodeValidation,
  joinCodeExpiry: joinCodeExpiryValidation,
  joinCodeMaxUse: joinCodeMaxUseValidation,
});

export type TenantFormData = z.infer<typeof tenantSchema>;
