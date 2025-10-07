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

export interface CreateExpenseForBill {
  amount: number;
  description: string;
  category: string;
  date: string;
}

export interface Expense {
  expenseId: string;
  amountPen: number;
  amountUsd: number;
  exchangeRate: number;
  currency: 'PEN' | 'USD';
  description: string;
  category: string;
  date: string;
  billId: string;
  userId: string;
  createdAt: string;
  updatedAt: string;
}

export interface Bill {
  billId: string;
  amountPen: number;
  amountUsd: number;
  description: string;
  category: string;
  currency: 'PEN' | 'USD';
  userId: string;
  date: string;
  expenses?: Expense[];
  createdAt: string;
  updatedAt: string;
}

export interface CreateBillRequest {
  description: string;
  category: string;
  userId: string;
  date: string;
  currency: 'PEN' | 'USD';
  exchangeRate: number;
  expenses: CreateExpenseForBill[];
}
