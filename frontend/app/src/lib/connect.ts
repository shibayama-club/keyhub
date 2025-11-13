import { createConnectTransport } from '@connectrpc/connect-web';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export const transport = createConnectTransport({
  baseUrl: API_BASE_URL,
  fetch: (input, init) => {
    return fetch(input, {
      ...init,
      credentials: 'include',
    });
  },
});
