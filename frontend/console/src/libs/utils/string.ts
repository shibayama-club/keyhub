// 空白文字や不可視文字のみで構成されているかをチェック
export const isBlankOrInvisible = (value: string): boolean => {
  // 空白文字、改行、タブなどの不可視文字のみで構成されている場合はtrue
  return /^\s*$/.test(value);
};
