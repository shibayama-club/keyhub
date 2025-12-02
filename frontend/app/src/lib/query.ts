import { Code, ConnectError } from '@connectrpc/connect';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';
import { getMe, logout } from '../../../gen/src/keyhub/app/v1/app-AuthService_connectquery';
import { getTenantByJoinCode, joinTenant } from '../../../gen/src/keyhub/app/v1/app-TenantService_connectquery';
import { getRoomsByTenant } from '../../../gen/src/keyhub/app/v1/room-RoomService_connectquery';
import { transport } from './connect';

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

export const useQueryGetMe = () => {
  return useQuery(getMe, {}, { enabled: false });
};

export const useMutationLogout = () => {
  return useMutation(logout);
};

export const useQueryGetTenantByJoinCode = (joinCode: string) => {
  return useQuery(getTenantByJoinCode, { joinCode }, { enabled: !!joinCode });
};

export const useMutationJoinTenant = () => {
  return useMutation(joinTenant);
};

export const useQueryGetRoomsByTenant = (tenantId: string) => {
  return useQuery(getRoomsByTenant, { tenantId }, { transport, enabled: !!tenantId });
};
