import { timestampDate } from '@bufbuild/protobuf/wkt';
import type { Timestamp } from '@bufbuild/protobuf/wkt';

/**
 * protobufのTimestampを日本語フォーマットの日付文字列に変換
 * @param timestamp - protobuf Timestamp
 * @returns フォーマット済み日付文字列 (例: "2024年11月18日")
 */
export const formatTimestampToJapaneseDate = (timestamp?: Timestamp): string => {
  if (!timestamp) {
    return '-';
  }

  return timestampDate(timestamp).toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};
