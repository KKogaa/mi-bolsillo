import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { MainLayout } from '../layouts/MainLayout';
import { billService } from '../services/api';
import type { Bill } from '../types';

export const Dashboard = () => {
  const [bills, setBills] = useState<Bill[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [expandedBillId, setExpandedBillId] = useState<string | null>(null);

  useEffect(() => {
    loadBills();
  }, []);

  const toggleExpand = (billId: string) => {
    setExpandedBillId(expandedBillId === billId ? null : billId);
  };

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
    try {
      await billService.delete(id);
      await loadBills();
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
                <li key={bill.billId}>
                  <div>
                    <div
                      className="px-4 py-4 sm:px-6 hover:bg-gray-50 cursor-pointer"
                      onClick={() => toggleExpand(bill.billId)}
                    >
                      <div className="flex items-center justify-between">
                        <div className="flex-1">
                          <div className="flex items-center justify-between">
                            <div>
                              <div className="flex items-center gap-2">
                                <p className="text-lg font-medium text-gray-900">
                                  {bill.description || 'Untitled Bill'}
                                </p>
                                <span className="text-xs px-2 py-0.5 bg-blue-100 text-blue-700 rounded font-medium">
                                  {bill.currency}
                                </span>
                              </div>
                              {bill.category && (
                                <p className="text-sm text-gray-600 mt-1">
                                  {bill.category}
                                </p>
                              )}
                              <p className="text-sm text-gray-500 mt-1">
                                {new Date(bill.date).toLocaleDateString()}
                                {bill.expenses && bill.expenses.length > 0 && (
                                  <span className="ml-2">
                                    • {bill.expenses.length} expense{bill.expenses.length !== 1 ? 's' : ''}
                                  </span>
                                )}
                              </p>
                            </div>
                            <div className="text-right">
                              <p className="text-xl font-semibold text-gray-900">
                                {bill.currency === 'PEN' ? 'S/' : '$'} {(bill.currency === 'PEN' ? bill.amountPen : bill.amountUsd).toFixed(2)}
                              </p>
                            </div>
                          </div>
                        </div>
                        <div className="ml-4 flex items-center gap-2">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDelete(bill.billId);
                            }}
                            className="text-red-600 hover:text-red-800 text-sm px-2"
                          >
                            Delete
                          </button>
                          <span className="text-gray-400">
                            {expandedBillId === bill.billId ? '▼' : '▶'}
                          </span>
                        </div>
                      </div>
                    </div>

                    {expandedBillId === bill.billId && bill.expenses && bill.expenses.length > 0 && (
                      <div className="px-4 py-3 sm:px-6 bg-gray-50 border-t border-gray-200">
                        <h4 className="text-sm font-medium text-gray-700 mb-3">Expenses</h4>
                        <ul className="space-y-2">
                          {bill.expenses.map((expense) => (
                            <li
                              key={expense.expenseId}
                              className="flex items-center justify-between bg-white px-3 py-2 rounded-md"
                            >
                              <div className="flex-1">
                                <p className="text-sm font-medium text-gray-900">
                                  {expense.description}
                                </p>
                                {expense.category && (
                                  <p className="text-xs text-gray-500 mt-1">
                                    {expense.category}
                                  </p>
                                )}
                                <p className="text-xs text-gray-400 mt-1">
                                  {new Date(expense.date).toLocaleDateString()}
                                </p>
                              </div>
                              <div className="text-right">
                                <p className="text-sm font-semibold text-gray-900">
                                  {bill.currency === 'PEN' ? 'S/' : '$'} {(bill.currency === 'PEN' ? expense.amountPen : expense.amountUsd).toFixed(2)}
                                </p>
                              </div>
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}
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
