import axios from 'axios';
import type { AuthResponse, LoginCredentials, SignUpCredentials, Bill, CreateBillRequest } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authService = {
  login: async (credentials: LoginCredentials): Promise<AuthResponse> => {
    const { data } = await api.post<AuthResponse>('/auth/login', credentials);
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    return data;
  },

  signup: async (credentials: SignUpCredentials): Promise<AuthResponse> => {
    const { data } = await api.post<AuthResponse>('/auth/signup', credentials);
    localStorage.setItem('token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
    return data;
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getCurrentUser: () => {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token');
  },
};

export const billService = {
  getAll: async (): Promise<Bill[]> => {
    const { data } = await api.get<Bill[]>('/bills');
    return data;
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
