import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { useTranslation } from 'react-i18next';
import type { MonthlyStatistics } from '../services/api';

interface MonthlyChartProps {
  data: MonthlyStatistics[];
}

export const MonthlyChart = ({ data }: MonthlyChartProps) => {
  const { t } = useTranslation();
  // Format data for the chart
  const chartData = data.map(stat => ({
    month: getMonthLabel(stat.monthNum, stat.year),
    PEN: stat.totalPen,
    USD: stat.totalUsd,
  })).reverse(); // Reverse to show oldest to newest

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      <h3 className="text-lg font-semibold text-gray-900 mb-4">{t('statistics.monthlySpending')}</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="month" />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="PEN" stroke="#3b82f6" strokeWidth={2} name="PEN (S/)" />
          <Line type="monotone" dataKey="USD" stroke="#10b981" strokeWidth={2} name="USD ($)" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

function getMonthLabel(monthNum: number, year: number): string {
  const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  return `${months[monthNum - 1]} ${year}`;
}
