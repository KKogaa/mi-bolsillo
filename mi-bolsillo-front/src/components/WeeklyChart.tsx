import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { useTranslation } from 'react-i18next';
import type { WeeklyStatistics } from '../services/api';

interface WeeklyChartProps {
  data: WeeklyStatistics[];
}

export const WeeklyChart = ({ data }: WeeklyChartProps) => {
  const { t } = useTranslation();
  // Format data for the chart
  const chartData = data.map(stat => ({
    week: stat.weekLabel,
    PEN: stat.totalPen,
    USD: stat.totalUsd,
  })).reverse(); // Reverse to show oldest to newest

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">{t('statistics.weeklySpending')}</h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="week" angle={-45} textAnchor="end" height={100} />
          <YAxis />
          <Tooltip />
          <Legend />
          <Bar dataKey="PEN" fill="#3b82f6" name="PEN (S/)" />
          <Bar dataKey="USD" fill="#10b981" name="USD ($)" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};
