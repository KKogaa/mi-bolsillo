import axios from 'axios';
import type { Bill, CreateBillRequest } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
    'ngrok-skip-browser-warning': 'true',
  },
});

// Token will be set by the setAuthToken function from Clerk
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Only redirect to login if we get a 401 and we're not already on login/signup pages
    if (error.response?.status === 401 &&
        !window.location.pathname.includes('/login') &&
        !window.location.pathname.includes('/signup')) {
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const setAuthToken = (token: string | null) => {
  if (token) {
    api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  } else {
    delete api.defaults.headers.common['Authorization'];
  }
};

export const billService = {
  getAll: async (): Promise<Bill[]> => {
    const { data } = await api.get<Bill[]>('/bills');
    return data || [];
  },

  getById: async (id: string): Promise<Bill> => {
    const { data } = await api.get<Bill>(`/bills/${id}`);
    return data;
  },

  create: async (bill: CreateBillRequest): Promise<Bill> => {
    const { data } = await api.post<Bill>('/bills', bill);
    return data;
  },

  update: async (id: string, bill: Partial<CreateBillRequest>): Promise<Bill> => {
    const { data } = await api.put<Bill>(`/bills/${id}`, bill);
    return data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/bills/${id}`);
  },
};

export default api;
