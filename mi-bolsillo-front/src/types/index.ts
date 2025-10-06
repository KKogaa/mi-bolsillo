export interface User {
  id: string;
  email: string;
  name: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface SignUpCredentials {
  name: string;
  email: string;
  password: string;
}

export interface Expense {
  id?: string;
  name: string;
  amount: number;
  category?: string;
}

export interface Bill {
  id: string;
  name: string;
  description?: string;
  totalAmount: number;
  expenses: Expense[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateBillRequest {
  name: string;
  description?: string;
  expenses: Expense[];
}
