import { createBrowserRouter, Navigate } from 'react-router-dom';
import * as Sentry from '@sentry/react';
import { AuthGuard } from '../components/AuthGuard';
import { ErrorFallback } from '../components/ErrorFallback';
import { LoginPage } from '../pages/LoginPage';
import { CallbackPage } from '../pages/CallbackPage';
import { HomePage } from '../pages/HomePage';
import { JoinTenantPage } from '../pages/JoinTenantPage';

const sentryCreateBrowserRouter = Sentry.wrapCreateBrowserRouterV6(createBrowserRouter);

export const router = sentryCreateBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
    errorElement: <ErrorFallback />,
  },
  {
    path: '/callback',
    element: <CallbackPage />,
    errorElement: <ErrorFallback />,
  },
  {
    path: '/',
    element: <AuthGuard />,
    errorElement: <ErrorFallback />,
    children: [
      {
        index: true,
        element: <Navigate to="/home" replace />,
      },
      {
        path: 'home',
        element: <HomePage />,
      },
      {
        path: 'join-tenant',
        element: <JoinTenantPage />,
      },
    ],
  },
]);
