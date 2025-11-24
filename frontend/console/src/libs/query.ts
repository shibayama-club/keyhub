import { Code, ConnectError } from '@connectrpc/connect';
import { useMutation, useQuery } from '@connectrpc/connect-query';
import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';
import { loginWithOrgId, logout } from '../../../gen/src/keyhub/console/v1/console-ConsoleAuthService_connectquery';
import {
  createTenant,
  getAllTenants,
  getTenantById,
} from '../../../gen/src/keyhub/console/v1/console-ConsoleService_connectquery';
import {
  createRoom,
  assignRoomToTenant,
} from '../../../gen/src/keyhub/console/v1/room-ConsoleRoomService_connectquery';

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

export const useMutationLoginWithOrgId = () => {
  return useMutation(loginWithOrgId);
};

export const useMutationLogout = () => {
  return useMutation(logout);
};

export const useMutationCreateTenant = () => {
  return useMutation(createTenant);
};

export const useQueryGetAllTenants = () => {
  return useQuery(getAllTenants, {});
};

export const useQueryGetTenantById = (id: string) => {
  return useQuery(getTenantById, { id });
};

export const useMutationCreateRoom = () => {
  return useMutation(createRoom);
};

export const useMutationAssignRoomToTenant = () => {
  return useMutation(assignRoomToTenant);
};
