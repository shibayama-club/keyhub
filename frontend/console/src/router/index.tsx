import { createBrowserRouter, Navigate } from 'react-router-dom';
import * as Sentry from '@sentry/react';
import AuthGuard from '../components/AuthGuard';
import ErrorFallback from '../components/ErrorFallback';
import LoginPage from '../pages/LoginPage';
import DashboardPage from '../pages/DashboardPage';

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
      ],
    },
  ],
  {
    basename: '/console',
  },
);
