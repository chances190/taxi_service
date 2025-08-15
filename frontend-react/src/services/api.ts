import axios, { AxiosResponse } from 'axios';

const api = axios.create({
  baseURL: (import.meta as ImportMeta).env.VITE_API_URL || 'http://localhost:3000',
  timeout: 10000
});

api.interceptors.response.use(
  (r: AxiosResponse) => r,
  (err: unknown) => {
    // centralizar logging futuramente
    return Promise.reject(err);
  }
);

export default api;
