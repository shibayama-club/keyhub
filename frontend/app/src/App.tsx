import { RouterProvider } from 'react-router-dom';
import { QueryClientProvider } from '@tanstack/react-query';
import { TransportProvider } from '@connectrpc/connect-query';
import { Toaster } from 'react-hot-toast';
import { router } from './router';
import { queryClient } from './lib/query';
import { transport } from './lib/connect';

function App() {
  return (
    <TransportProvider transport={transport}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
        <Toaster position="top-right" />
      </QueryClientProvider>
    </TransportProvider>
  );
}

export default App;
