import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';
import { useTranslation } from 'react-i18next';
import type { CategoryStatistics } from '../services/api';

interface CategoryChartProps {
  data: CategoryStatistics[];
}

const COLORS = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#14b8a6', '#f97316'];

export const CategoryChart = ({ data }: CategoryChartProps) => {
  const { t } = useTranslation();
  // Format data for the pie chart
  const chartData = data.map(stat => ({
    name: stat.category,
    value: stat.totalPen,
    percentage: stat.percentage,
  }));

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">{t('statistics.spendingByCategory')}</h3>

      {chartData.length === 0 ? (
        <div className="flex items-center justify-center h-64 text-gray-500">
          {t('statistics.noDataAvailable')}
        </div>
      ) : (
        <>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={chartData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={(entry: any) => `${(entry.percent * 100).toFixed(1)}%`}
                outerRadius={100}
                fill="#8884d8"
                dataKey="value"
              >
                {chartData.map((_entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value: number) => `S/ ${value.toFixed(2)}`} />
              <Legend />
            </PieChart>
          </ResponsiveContainer>

          <div className="mt-4 space-y-2">
            {data.map((stat, index) => (
              <div key={stat.category} className="flex items-center justify-between text-sm">
                <div className="flex items-center gap-2">
                  <div
                    className="w-3 h-3 rounded-full"
                    style={{ backgroundColor: COLORS[index % COLORS.length] }}
                  ></div>
                  <span className="text-gray-700">{stat.category}</span>
                </div>
                <div className="text-right">
                  <span className="font-medium text-gray-900">S/ {stat.totalPen.toFixed(2)}</span>
                  <span className="text-gray-500 ml-2">({stat.percentage.toFixed(1)}%)</span>
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
};
