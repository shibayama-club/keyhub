import { Code, ConnectError } from '@connectrpc/connect';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';

const retry = (failureCount: number, err: unknown) => {
  if (err instanceof ConnectError) {
    if (err.code === Code.PermissionDenied || err.code === Code.Unauthenticated) {
      return false;
    }
  }
  return failureCount < 3;
};

const onError = (err: unknown) => {
  console.error(err);
};

export const queryClient = new QueryClient({
  queryCache: new QueryCache({ onError }),
  mutationCache: new MutationCache({ onError }),
  defaultOptions: {
    queries: {
      retry,
      staleTime: 60 * 1000,
    },
    mutations: {
      retry: false,
    },
  },
});
