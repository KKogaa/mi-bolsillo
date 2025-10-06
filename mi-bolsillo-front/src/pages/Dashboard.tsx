import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { billService } from '../services/api';
import type { Bill } from '../types';

export const Dashboard = () => {
  const [bills, setBills] = useState<Bill[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadBills();
  }, []);

  const loadBills = async () => {
    try {
      setIsLoading(true);
      const data = await billService.getAll();
      setBills(data);
    } catch (err) {
      setError('Failed to load bills');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this bill?')) {
      return;
    }

    try {
      await billService.delete(id);
      setBills(bills.filter((bill) => bill.id !== id));
    } catch (err) {
      setError('Failed to delete bill');
    }
  };

  return (
    <MainLayout>
      <div className="px-4 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900">Bills</h2>
          <Link
            to="/bills/new"
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            New Bill
          </Link>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 text-red-600 rounded-md">
            {error}
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <div className="text-gray-600">Loading...</div>
          </div>
        ) : bills.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600 mb-4">No bills yet</p>
            <Link
              to="/bills/new"
              className="text-blue-600 hover:text-blue-500 font-medium"
            >
              Create your first bill
            </Link>
          </div>
        ) : (
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
              {bills.map((bill) => (
                <li key={bill.id}>
                  <div className="px-4 py-4 sm:px-6 hover:bg-gray-50">
                    <div className="flex items-center justify-between">
                      <Link
                        to={`/bills/${bill.id}`}
                        className="flex-1"
                      >
                        <div className="flex items-center justify-between">
                          <div>
                            <p className="text-lg font-medium text-gray-900">
                              {bill.name}
                            </p>
                            {bill.description && (
                              <p className="text-sm text-gray-600 mt-1">
                                {bill.description}
                              </p>
                            )}
                            <p className="text-sm text-gray-500 mt-1">
                              {bill.expenses.length} expense{bill.expenses.length !== 1 ? 's' : ''}
                            </p>
                          </div>
                          <div className="text-right">
                            <p className="text-lg font-semibold text-gray-900">
                              ${bill.totalAmount.toFixed(2)}
                            </p>
                            <p className="text-xs text-gray-500 mt-1">
                              {new Date(bill.createdAt).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                      </Link>
                      <button
                        onClick={(e) => {
                          e.preventDefault();
                          handleDelete(bill.id);
                        }}
                        className="ml-4 text-red-600 hover:text-red-800 text-sm"
                      >
                        Delete
                      </button>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </MainLayout>
  );
};
