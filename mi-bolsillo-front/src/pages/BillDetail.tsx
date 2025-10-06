import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { billService } from '../services/api';
import type { Bill } from '../types';

export const BillDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [bill, setBill] = useState<Bill | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      loadBill(id);
    }
  }, [id]);

  const loadBill = async (billId: string) => {
    try {
      setIsLoading(true);
      const data = await billService.getById(billId);
      setBill(data);
    } catch (err) {
      setError('Failed to load bill');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!id || !window.confirm('Are you sure you want to delete this bill?')) {
      return;
    }

    try {
      await billService.delete(id);
      navigate('/');
    } catch (err) {
      setError('Failed to delete bill');
    }
  };

  if (isLoading) {
    return (
      <MainLayout>
        <div className="text-center py-12">
          <div className="text-gray-600">Loading...</div>
        </div>
      </MainLayout>
    );
  }

  if (error || !bill) {
    return (
      <MainLayout>
        <div className="text-center py-12">
          <div className="text-red-600 mb-4">{error || 'Bill not found'}</div>
          <Link to="/" className="text-blue-600 hover:text-blue-500">
            Back to dashboard
          </Link>
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="max-w-3xl mx-auto px-4 sm:px-0">
        <div className="mb-6 flex justify-between items-start">
          <div>
            <Link to="/" className="text-blue-600 hover:text-blue-500 text-sm mb-2 inline-block">
              ← Back to bills
            </Link>
            <h2 className="text-2xl font-bold text-gray-900">{bill.name}</h2>
            {bill.description && (
              <p className="text-gray-600 mt-1">{bill.description}</p>
            )}
            <p className="text-sm text-gray-500 mt-2">
              Created {new Date(bill.createdAt).toLocaleDateString()}
            </p>
          </div>
          <button
            onClick={handleDelete}
            className="px-4 py-2 text-red-600 hover:text-red-800 border border-red-600 rounded-md"
          >
            Delete
          </button>
        </div>

        <div className="bg-white shadow rounded-lg overflow-hidden">
          <div className="px-6 py-4 bg-gray-50 border-b border-gray-200">
            <h3 className="text-lg font-medium text-gray-900">Expenses</h3>
          </div>
          <ul className="divide-y divide-gray-200">
            {bill.expenses.map((expense, index) => (
              <li key={expense.id || index} className="px-6 py-4">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-medium text-gray-900">{expense.name}</p>
                    {expense.category && (
                      <p className="text-sm text-gray-600 mt-1">
                        {expense.category}
                      </p>
                    )}
                  </div>
                  <span className="font-semibold text-gray-900">
                    ${expense.amount.toFixed(2)}
                  </span>
                </div>
              </li>
            ))}
          </ul>
          <div className="px-6 py-4 bg-gray-50 border-t border-gray-200">
            <div className="flex justify-between items-center">
              <span className="text-lg font-medium text-gray-900">Total</span>
              <span className="text-2xl font-bold text-gray-900">
                ${bill.totalAmount.toFixed(2)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
};
