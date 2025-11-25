import type { Key } from '../../../gen/src/keyhub/console/v1/common_pb';
import { getKeyStatusLabel, getKeyStatusTextColor } from '../utils/key';

type KeyBadgeListProps = {
  keys: Key[];
};

export const KeyBadgeList = ({ keys }: KeyBadgeListProps) => {
  if (!keys || keys.length === 0) {
    return null;
  }

  return (
    <div className="mt-2">
      <p className="text-xs font-medium text-gray-500 mb-1">鍵情報:</p>
      <div className="flex flex-wrap gap-2">
        {keys.map((key) => (
          <span
            key={key.id}
            className="inline-flex items-center space-x-1 rounded-md bg-gray-100 px-2 py-1 text-xs"
          >
            <span className="font-medium text-gray-900">{key.keyNumber}</span>
            <span className="text-gray-400">-</span>
            <span className={getKeyStatusTextColor(key.status)}>{getKeyStatusLabel(key.status)}</span>
          </span>
        ))}
      </div>
    </div>
  );
};
