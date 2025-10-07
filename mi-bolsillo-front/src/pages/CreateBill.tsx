import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUser } from '@clerk/clerk-react';
import { MainLayout } from '../layouts/MainLayout';
import { billService } from '../services/api';
import type { CreateExpenseForBill } from '../types';

export const CreateBill = () => {
  const navigate = useNavigate();
  const { user } = useUser();
  const [description, setDescription] = useState('');
  const [category, setCategory] = useState('');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [currency, setCurrency] = useState<'PEN' | 'USD'>('PEN');
  const [exchangeRate, setExchangeRate] = useState<number>(3.74);
  const [expenses, setExpenses] = useState<CreateExpenseForBill[]>([
    {
      amount: 0,
      description: '',
      category: '',
      date: new Date().toISOString().split('T')[0]
    }
  ]);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const addExpense = () => {
    setExpenses([...expenses, {
      amount: 0,
      description: '',
      category: '',
      date: new Date().toISOString().split('T')[0]
    }]);
  };

  const removeExpense = (index: number) => {
    setExpenses(expenses.filter((_, i) => i !== index));
  };

  const updateExpense = (index: number, field: keyof CreateExpenseForBill, value: string | number) => {
    const newExpenses = [...expenses];
    newExpenses[index] = { ...newExpenses[index], [field]: value };
    setExpenses(newExpenses);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');

    const validExpenses = expenses.filter(exp => exp.description && exp.amount > 0);

    if (validExpenses.length === 0) {
      setError('Please add at least one expense with a description and amount');
      return;
    }

    if (!description) {
      setError('Please provide a bill description');
      return;
    }

    if (!category) {
      setError('Please provide a bill category');
      return;
    }

    if (!user?.id) {
      setError('User not authenticated');
      return;
    }

    setIsLoading(true);

    try {
      await billService.create({
        description,
        category,
        userId: user.id,
        date: new Date(date).toISOString(),
        currency,
        exchangeRate,
        expenses: validExpenses,
      });
      navigate('/');
    } catch (err) {
      setError('Failed to create bill');
    } finally {
      setIsLoading(false);
    }
  };

  const totalAmount = expenses.reduce((sum, exp) => sum + (Number(exp.amount) || 0), 0);

  return (
    <MainLayout>
      <div className="max-w-3xl mx-auto px-4 sm:px-0">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900">Create New Bill</h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div className="bg-white shadow rounded-lg p-6">
            <div className="space-y-4">
              <div>
                <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                  Description
                </label>
                <textarea
                  id="description"
                  required
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  rows={2}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="e.g., Office supplies and equipment purchase"
                />
              </div>

              <div>
                <label htmlFor="category" className="block text-sm font-medium text-gray-700">
                  Category
                </label>
                <input
                  id="category"
                  type="text"
                  required
                  value={category}
                  onChange={(e) => setCategory(e.target.value)}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="e.g., Office Expenses"
                />
              </div>

              <div>
                <label htmlFor="date" className="block text-sm font-medium text-gray-700">
                  Date
                </label>
                <input
                  id="date"
                  type="date"
                  required
                  value={date}
                  onChange={(e) => setDate(e.target.value)}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label htmlFor="currency" className="block text-sm font-medium text-gray-700">
                    Currency
                  </label>
                  <select
                    id="currency"
                    required
                    value={currency}
                    onChange={(e) => setCurrency(e.target.value as 'PEN' | 'USD')}
                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="PEN">PEN</option>
                    <option value="USD">USD</option>
                  </select>
                </div>
                <div>
                  <label htmlFor="exchangeRate" className="block text-sm font-medium text-gray-700">
                    Exchange Rate
                  </label>
                  <input
                    id="exchangeRate"
                    type="number"
                    step="0.0001"
                    min="0"
                    required
                    value={exchangeRate || ''}
                    onChange={(e) => setExchangeRate(parseFloat(e.target.value) || 0)}
                    className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="3.74"
                  />
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white shadow rounded-lg p-6">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-medium text-gray-900">Expenses</h3>
              <button
                type="button"
                onClick={addExpense}
                className="px-3 py-1 text-sm text-blue-600 hover:text-blue-700 font-medium"
              >
                + Add Expense
              </button>
            </div>

            <div className="space-y-4">
              {expenses.map((expense, index) => (
                <div key={index} className="border border-gray-200 rounded-lg p-4 space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium text-gray-700">Expense {index + 1}</span>
                    {expenses.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeExpense(index)}
                        className="text-red-600 hover:text-red-800 text-xl"
                      >
                        ×
                      </button>
                    )}
                  </div>

                  <div>
                    <input
                      type="text"
                      value={expense.description}
                      onChange={(e) => updateExpense(index, 'description', e.target.value)}
                      placeholder="Expense description"
                      className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-3">
                    <div>
                      <input
                        type="number"
                        step="0.01"
                        min="0"
                        value={expense.amount || ''}
                        onChange={(e) => updateExpense(index, 'amount', parseFloat(e.target.value) || 0)}
                        placeholder="Amount"
                        className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <div>
                      <input
                        type="text"
                        value={expense.category || ''}
                        onChange={(e) => updateExpense(index, 'category', e.target.value)}
                        placeholder="Category"
                        className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                  </div>

                  <div>
                    <input
                      type="date"
                      value={expense.date}
                      onChange={(e) => updateExpense(index, 'date', e.target.value)}
                      className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>
                </div>
              ))}
            </div>

            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex justify-between items-center">
                <span className="text-lg font-medium text-gray-900">Total</span>
                <span className="text-2xl font-bold text-gray-900">
                  {currency === 'PEN' ? 'S/' : '$'} {totalAmount.toFixed(2)}
                </span>
              </div>
            </div>
          </div>

          {error && (
            <div className="p-4 bg-red-50 border border-red-200 text-red-600 rounded-md">
              {error}
            </div>
          )}

          <div className="flex gap-3">
            <button
              type="button"
              onClick={() => navigate('/')}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? 'Creating...' : 'Create Bill'}
            </button>
          </div>
        </form>
      </div>
    </MainLayout>
  );
};
