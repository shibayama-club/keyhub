import { useRouteError, isRouteErrorResponse, useNavigate } from 'react-router-dom';

export const ErrorFallback = () => {
  const error = useRouteError();
  const navigate = useNavigate();

  let errorMessage: string;

  if (isRouteErrorResponse(error)) {
    errorMessage = error.statusText || error.data?.message || 'Unknown error occurred';
  } else if (error instanceof Error) {
    errorMessage = error.message;
  } else {
    errorMessage = 'An unexpected error occurred';
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-gray-50 px-6 py-12">
      <div className="text-center">
        <h1 className="text-6xl font-bold text-gray-900">Oops!</h1>
        <p className="mt-4 text-xl text-gray-600">Something went wrong</p>
        <p className="mt-2 text-sm text-gray-500">{errorMessage}</p>
        <div className="mt-8 flex justify-center gap-4">
          <button
            onClick={() => navigate(-1)}
            className="rounded-md bg-white px-4 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-gray-300 ring-inset hover:bg-gray-50"
          >
            Go Back
          </button>
          <button
            onClick={() => navigate('/')}
            className="rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
          >
            Go Home
          </button>
        </div>
      </div>
    </div>
  );
};
