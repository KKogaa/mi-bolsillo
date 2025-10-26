import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { MainLayout } from '../layouts/MainLayout';
import { statisticsService, type DashboardStatistics } from '../services/api';
import { MonthlyChart } from '../components/MonthlyChart';
import { WeeklyChart } from '../components/WeeklyChart';
import { CategoryChart } from '../components/CategoryChart';

export const Statistics = () => {
  const { t } = useTranslation();
  const [statistics, setStatistics] = useState<DashboardStatistics | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [months, setMonths] = useState(6);

  useEffect(() => {
    loadStatistics();
  }, [months]);

  const loadStatistics = async () => {
    try {
      setIsLoading(true);
      setError('');
      const data = await statisticsService.getDashboardStatistics(months);
      setStatistics(data);
    } catch (err) {
      console.error('Failed to load statistics:', err);
      setError(t('statistics.errorLoad'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <MainLayout>
      <div className="px-4 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-gray-900">{t('statistics.title')}</h2>
          <div className="flex items-center gap-2">
            <label htmlFor="months" className="text-sm text-gray-600">
              {t('statistics.showLast')}
            </label>
            <select
              id="months"
              value={months}
              onChange={(e) => setMonths(Number(e.target.value))}
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value={3}>3 {t('statistics.months')}</option>
              <option value={6}>6 {t('statistics.months')}</option>
              <option value={12}>12 {t('statistics.months')}</option>
              <option value={24}>24 {t('statistics.months')}</option>
            </select>
          </div>
        </div>

        {error && (
          <div className="mb-4 p-4 bg-red-50 border border-red-200 text-red-600 rounded-md">
            {error}
          </div>
        )}

        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">{t('statistics.loading')}</p>
          </div>
        ) : statistics ? (
          <>
            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
              <div className="bg-white p-6 rounded-lg shadow">
                <h3 className="text-sm font-medium text-gray-500 mb-2">{t('statistics.totalSpentPEN')}</h3>
                <p className="text-3xl font-bold text-gray-900">S/ {statistics.totalPen.toFixed(2)}</p>
                <p className="text-sm text-gray-500 mt-2">{statistics.totalBills} {t('statistics.bills')}</p>
              </div>
              <div className="bg-white p-6 rounded-lg shadow">
                <h3 className="text-sm font-medium text-gray-500 mb-2">{t('statistics.totalSpentUSD')}</h3>
                <p className="text-3xl font-bold text-gray-900">$ {statistics.totalUsd.toFixed(2)}</p>
                <p className="text-sm text-gray-500 mt-2">{statistics.totalBills} {t('statistics.bills')}</p>
              </div>
              <div className="bg-white p-6 rounded-lg shadow">
                <h3 className="text-sm font-medium text-gray-500 mb-2">{t('statistics.averagePerBill')}</h3>
                <p className="text-3xl font-bold text-gray-900">
                  S/ {statistics.totalBills > 0 ? (statistics.totalPen / statistics.totalBills).toFixed(2) : '0.00'}
                </p>
                <p className="text-sm text-gray-500 mt-2">
                  $ {statistics.totalBills > 0 ? (statistics.totalUsd / statistics.totalBills).toFixed(2) : '0.00'}
                </p>
              </div>
            </div>

            {/* Charts */}
            <div className="space-y-6">
              {statistics.monthlyStats.length > 0 && (
                <MonthlyChart data={statistics.monthlyStats} />
              )}

              {statistics.weeklyStats.length > 0 && (
                <WeeklyChart data={statistics.weeklyStats} />
              )}

              {statistics.categoryStats.length > 0 && (
                <CategoryChart data={statistics.categoryStats} />
              )}
            </div>
          </>
        ) : (
          <div className="text-center py-12">
            <p className="text-gray-600">No statistics available</p>
          </div>
        )}
      </div>
    </MainLayout>
  );
};
