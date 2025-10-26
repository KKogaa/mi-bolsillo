import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { MainLayout } from '../layouts/MainLayout';
import { billService, authService } from '../services/api';
import type { Bill } from '../types';

export const Dashboard = () => {
  const { t } = useTranslation();
  const [bills, setBills] = useState<Bill[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [expandedBillId, setExpandedBillId] = useState<string | null>(null);
  const [isLinked, setIsLinked] = useState(false);

  useEffect(() => {
    loadBills();
    checkLinkStatus();
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
      setError(t('dashboard.errorLoad'));
    } finally {
      setIsLoading(false);
    }
  };

  const checkLinkStatus = async () => {
    try {
      const { isLinked } = await authService.getLinkStatus();
      setIsLinked(isLinked);
    } catch (err) {
      console.error('Failed to check link status:', err);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await billService.delete(id);
      await loadBills();
    } catch (err) {
      setError(t('dashboard.errorDelete'));
    }
  };

  return (
    <MainLayout>
      <div className="px-4 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900">{t('dashboard.title')}</h2>
          <div className="flex gap-3">
            {!isLinked && (
              <Link
                to="/link-telegram"
                className="px-4 py-2 bg-white text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                {t('dashboard.linkTelegram')}
              </Link>
            )}
            {isLinked && (
              <a
                href={`https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_USERNAME || 'mi_bolsillo_bot'}`}
                target="_blank"
                rel="noopener noreferrer"
                className="px-4 py-2 bg-white text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500 flex items-center gap-2"
              >
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.223-.548.223l.188-2.85 5.18-4.68c.223-.198-.054-.308-.346-.11l-6.4 4.03-2.76-.918c-.6-.187-.612-.6.125-.89l10.782-4.156c.498-.187.935.112.77.89z"/>
                </svg>
                {t('dashboard.openBot')}
              </a>
            )}
            <Link
              to="/bills/new"
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {t('dashboard.newBill')}
            </Link>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 text-red-600 rounded-md">
            {error}
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <div className="text-gray-600">{t('dashboard.loading')}</div>
          </div>
        ) : bills.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600 mb-4">{t('dashboard.noBills')}</p>
            <Link
              to="/bills/new"
              className="text-blue-600 hover:text-blue-500 font-medium"
            >
              {t('dashboard.createFirst')}
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
                                    • {bill.expenses.length} {bill.expenses.length !== 1 ? t('dashboard.expenses_plural') : t('dashboard.expenses')}
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
                            {t('dashboard.delete')}
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
