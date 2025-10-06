import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { billService } from '../services/api';
import type { Expense } from '../types';

export const CreateBill = () => {
  const navigate = useNavigate();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [expenses, setExpenses] = useState<Expense[]>([
    { name: '', amount: 0, category: '' }
  ]);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const addExpense = () => {
    setExpenses([...expenses, { name: '', amount: 0, category: '' }]);
  };

  const removeExpense = (index: number) => {
    setExpenses(expenses.filter((_, i) => i !== index));
  };

  const updateExpense = (index: number, field: keyof Expense, value: string | number) => {
    const newExpenses = [...expenses];
    newExpenses[index] = { ...newExpenses[index], [field]: value };
    setExpenses(newExpenses);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');

    const validExpenses = expenses.filter(exp => exp.name && exp.amount > 0);

    if (validExpenses.length === 0) {
      setError('Please add at least one expense with a name and amount');
      return;
    }

    setIsLoading(true);

    try {
      await billService.create({
        name,
        description: description || undefined,
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
                <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                  Bill Name
                </label>
                <input
                  id="name"
                  type="text"
                  required
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="e.g., Groceries, Restaurant Bill"
                />
              </div>

              <div>
                <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                  Description (optional)
                </label>
                <textarea
                  id="description"
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  rows={2}
                  className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Add any notes about this bill"
                />
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
                <div key={index} className="flex gap-3 items-start">
                  <div className="flex-1">
                    <input
                      type="text"
                      value={expense.name}
                      onChange={(e) => updateExpense(index, 'name', e.target.value)}
                      placeholder="Expense name"
                      className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>
                  <div className="w-32">
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
                  <div className="flex-1">
                    <input
                      type="text"
                      value={expense.category || ''}
                      onChange={(e) => updateExpense(index, 'category', e.target.value)}
                      placeholder="Category (optional)"
                      className="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>
                  {expenses.length > 1 && (
                    <button
                      type="button"
                      onClick={() => removeExpense(index)}
                      className="px-3 py-2 text-red-600 hover:text-red-800"
                    >
                      ×
                    </button>
                  )}
                </div>
              ))}
            </div>

            <div className="mt-4 pt-4 border-t border-gray-200">
              <div className="flex justify-between items-center">
                <span className="text-lg font-medium text-gray-900">Total</span>
                <span className="text-2xl font-bold text-gray-900">
                  ${totalAmount.toFixed(2)}
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
