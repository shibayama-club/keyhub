import { createBrowserRouter, Navigate } from 'react-router-dom';
import * as Sentry from '@sentry/react';
import { AuthGuard } from '../components/AuthGuard';
import { ErrorFallback } from '../components/ErrorFallback';
import { LoginPage } from '../pages/LoginPage';
import { DashboardPage } from '../pages/DashboardPage';
import { TenantsPage } from '../pages/TenantsPage';
import { CreateTenantPage } from '../pages/CreateTenantPage';
import { CreateRoomPage } from '../pages/CreateRoomPage';
import { AssignRoomPage } from '../pages/AssignRoomPage';

// SentryでラップされたcreateBrowserRouter
const sentryCreateBrowserRouter = Sentry.wrapCreateBrowserRouterV6(createBrowserRouter);

export const router = sentryCreateBrowserRouter(
  [
    {
      path: '/login',
      element: <LoginPage />,
      errorElement: <ErrorFallback />,
    },
    {
      path: '/',
      element: <AuthGuard />,
      errorElement: <ErrorFallback />,
      children: [
        {
          index: true,
          element: <Navigate to="dashboard" replace />,
        },
        {
          path: 'dashboard',
          element: <DashboardPage />,
        },
        {
          path: 'tenants',
          element: <TenantsPage />,
        },
        {
          path: 'tenants/create',
          element: <CreateTenantPage />,
        },
        {
          path: 'tenants/:tenantId/assign-room',
          element: <AssignRoomPage />,
        },
        {
          path: 'rooms/create',
          element: <CreateRoomPage />,
        },
      ],
    },
  ],
  {
    basename: '/console',
  },
);
