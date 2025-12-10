import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { MainLayout } from '../layouts/MainLayout';
import { billService, authService } from '../services/api';
import type { Bill } from '../types';
import { config } from '../config';

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
      <div className="px-4 sm:px-6">
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-6">
          <h2 className="text-xl sm:text-2xl font-bold text-gray-900">{t('dashboard.title')}</h2>
          <div className="flex flex-col sm:flex-row gap-2 sm:gap-3">
            {!isLinked && (
              <Link
                to="/link-telegram"
                className="px-4 py-2.5 bg-white text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500 text-center text-sm sm:text-base"
              >
                {t('dashboard.linkTelegram')}
              </Link>
            )}
            {isLinked && (
              <a
                href={`https://t.me/${config.telegramBotUsername}`}
                target="_blank"
                rel="noopener noreferrer"
                className="px-4 py-2.5 bg-white text-blue-600 border border-blue-600 rounded-md hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500 flex items-center justify-center gap-2 text-sm sm:text-base"
              >
                <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.223-.548.223l.188-2.85 5.18-4.68c.223-.198-.054-.308-.346-.11l-6.4 4.03-2.76-.918c-.6-.187-.612-.6.125-.89l10.782-4.156c.498-.187.935.112.77.89z"/>
                </svg>
                {t('dashboard.openBot')}
              </a>
            )}
            <Link
              to="/bills/new"
              className="px-4 py-2.5 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 text-center text-sm sm:text-base"
            >
              {t('dashboard.newBill')}
            </Link>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-3 sm:p-4 bg-red-50 border border-red-200 text-red-600 rounded-md text-sm sm:text-base">
            {error}
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-12 px-4">
            <div className="text-gray-600 text-sm sm:text-base">{t('dashboard.loading')}</div>
          </div>
        ) : bills.length === 0 ? (
          <div className="text-center py-12 px-4">
            <p className="text-gray-600 mb-4 text-sm sm:text-base">{t('dashboard.noBills')}</p>
            <Link
              to="/bills/new"
              className="text-blue-600 hover:text-blue-500 font-medium text-sm sm:text-base"
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
                      <div className="flex flex-col gap-3">
                        <div className="flex items-start justify-between gap-3">
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center gap-2 flex-wrap">
                              <p className="text-base sm:text-lg font-medium text-gray-900 break-words">
                                {bill.description || 'Untitled Bill'}
                              </p>
                              <span className="text-xs px-2 py-0.5 bg-blue-100 text-blue-700 rounded font-medium whitespace-nowrap">
                                {bill.currency}
                              </span>
                            </div>
                            {bill.category && (
                              <p className="text-sm text-gray-600 mt-1.5">
                                {bill.category}
                              </p>
                            )}
                            <p className="text-xs sm:text-sm text-gray-500 mt-1.5">
                              {new Date(bill.date).toLocaleDateString()}
                              {bill.expenses && bill.expenses.length > 0 && (
                                <span className="ml-2">
                                  • {bill.expenses.length} {bill.expenses.length !== 1 ? t('dashboard.expenses_plural') : t('dashboard.expenses')}
                                </span>
                              )}
                            </p>
                          </div>
                          <div className="flex items-start gap-2 sm:gap-3">
                            <div className="text-right">
                              <p className="text-lg sm:text-xl font-semibold text-gray-900 whitespace-nowrap">
                                {bill.currency === 'PEN' ? 'S/' : '$'} {(bill.currency === 'PEN' ? bill.amountPen : bill.amountUsd).toFixed(2)}
                              </p>
                            </div>
                            <span className="text-gray-400 text-lg">
                              {expandedBillId === bill.billId ? '▼' : '▶'}
                            </span>
                          </div>
                        </div>
                        <div className="flex justify-end">
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              handleDelete(bill.billId);
                            }}
                            className="text-red-600 hover:text-red-800 text-sm font-medium px-3 py-1.5 rounded hover:bg-red-50"
                          >
                            {t('dashboard.delete')}
                          </button>
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
                              className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 bg-white px-3 py-3 rounded-md"
                            >
                              <div className="flex-1 min-w-0">
                                <p className="text-sm font-medium text-gray-900 break-words">
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
                              <div className="text-left sm:text-right">
                                <p className="text-sm sm:text-base font-semibold text-gray-900 whitespace-nowrap">
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
