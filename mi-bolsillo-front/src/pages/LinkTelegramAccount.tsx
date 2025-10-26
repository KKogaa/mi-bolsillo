import { useState, useEffect } from 'react';
import { MainLayout } from '../layouts/MainLayout';
import { authService } from '../services/api';

export const LinkTelegramAccount = () => {
  const [otpCode, setOtpCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [isLinked, setIsLinked] = useState(false);
  const [telegramId, setTelegramId] = useState<number | null>(null);
  const [checkingStatus, setCheckingStatus] = useState(true);

  useEffect(() => {
    checkLinkStatus();
  }, []);

  const checkLinkStatus = async () => {
    try {
      setCheckingStatus(true);
      const status = await authService.getLinkStatus();
      setIsLinked(status.isLinked);
      setTelegramId(status.telegramId || null);
    } catch (err) {
      console.error('Failed to check link status:', err);
    } finally {
      setCheckingStatus(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess(false);

    if (!otpCode || otpCode.length !== 6) {
      setError('Please enter a valid 6-digit code');
      return;
    }

    try {
      setIsLoading(true);
      const response = await authService.verifyOTP(otpCode);

      if (response.success) {
        setSuccess(true);
        setOtpCode('');
        // Refresh link status after successful linking
        await checkLinkStatus();
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || 'Failed to verify OTP code. Please try again.';
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  if (checkingStatus) {
    return (
      <MainLayout>
        <div className="max-w-2xl mx-auto px-4 sm:px-0">
          <div className="flex justify-center items-center py-12">
            <div className="text-center">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
              <p className="text-gray-600">Checking link status...</p>
            </div>
          </div>
        </div>
      </MainLayout>
    );
  }

  return (
    <MainLayout>
      <div className="max-w-2xl mx-auto px-4 sm:px-0">
        <div className="mb-6">
          <h2 className="text-2xl font-bold text-gray-900">Link Telegram Account</h2>
          <p className="mt-2 text-sm text-gray-600">
            Connect your Telegram account to manage your expenses from both platforms
          </p>
        </div>

        {isLinked && (
          <div className="mb-6 bg-green-50 border border-green-200 rounded-lg p-6">
            <div className="flex items-start">
              <div className="flex-shrink-0">
                <svg className="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" strokeWidth="2" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div className="ml-3 flex-1">
                <h3 className="text-lg font-medium text-green-900">
                  Account Already Linked
                </h3>
                <p className="mt-2 text-sm text-green-700">
                  Your account is already linked to Telegram.
                  {telegramId && (
                    <span className="block mt-1">
                      Telegram ID: <code className="px-2 py-1 bg-green-100 rounded text-xs font-mono">{telegramId}</code>
                    </span>
                  )}
                </p>
                <div className="mt-4">
                  <p className="text-sm text-green-700 font-medium">You can now:</p>
                  <ul className="mt-2 text-sm text-green-700 list-disc list-inside space-y-1">
                    <li>Add bills and expenses from Telegram</li>
                    <li>View your spending on both platforms</li>
                    <li>Get quick summaries via the bot</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="bg-white shadow sm:rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              {isLinked ? 'Link Another Account' : 'How to link your account'}
            </h3>

            <div className="mb-6 space-y-3">
              <div className="flex items-start">
                <span className="flex-shrink-0 w-6 h-6 flex items-center justify-center bg-blue-600 text-white text-sm rounded-full font-medium mr-3">
                  1
                </span>
                <p className="text-sm text-gray-700">
                  Open the Mi Bolsillo Telegram bot
                </p>
              </div>

              <div className="flex items-start">
                <span className="flex-shrink-0 w-6 h-6 flex items-center justify-center bg-blue-600 text-white text-sm rounded-full font-medium mr-3">
                  2
                </span>
                <p className="text-sm text-gray-700">
                  Send the command <code className="px-2 py-1 bg-gray-100 rounded text-xs font-mono">/link</code> to generate a verification code
                </p>
              </div>

              <div className="flex items-start">
                <span className="flex-shrink-0 w-6 h-6 flex items-center justify-center bg-blue-600 text-white text-sm rounded-full font-medium mr-3">
                  3
                </span>
                <p className="text-sm text-gray-700">
                  Enter the 6-digit code below within 5 minutes
                </p>
              </div>
            </div>

            {success && (
              <div className="mb-4 p-4 bg-green-50 border border-green-200 rounded-md">
                <div className="flex">
                  <div className="flex-shrink-0">
                    <svg className="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                    </svg>
                  </div>
                  <div className="ml-3">
                    <p className="text-sm font-medium text-green-800">
                      Success! Your Telegram account has been linked.
                    </p>
                    <p className="mt-1 text-sm text-green-700">
                      You can now manage your bills from both web and Telegram.
                    </p>
                  </div>
                </div>
              </div>
            )}

            {error && (
              <div className="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
                <p className="text-sm text-red-600">{error}</p>
              </div>
            )}

            <form onSubmit={handleSubmit}>
              <div className="mb-4">
                <label htmlFor="otpCode" className="block text-sm font-medium text-gray-700 mb-2">
                  Verification Code
                </label>
                <input
                  type="text"
                  id="otpCode"
                  value={otpCode}
                  onChange={(e) => {
                    const value = e.target.value.replace(/\D/g, '').slice(0, 6);
                    setOtpCode(value);
                  }}
                  placeholder="Enter 6-digit code"
                  className="w-full px-4 py-3 text-center text-2xl tracking-widest font-mono border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  maxLength={6}
                  disabled={isLoading || success}
                  autoComplete="off"
                />
                <p className="mt-2 text-xs text-gray-500">
                  The code expires in 5 minutes
                </p>
              </div>

              <button
                type="submit"
                disabled={isLoading || success || otpCode.length !== 6}
                className="w-full px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-gray-400 disabled:cursor-not-allowed"
              >
                {isLoading ? 'Verifying...' : success ? 'Account Linked' : 'Verify & Link Account'}
              </button>
            </form>
          </div>
        </div>

        <div className="mt-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3 flex-1">
              <h3 className="text-sm font-medium text-blue-800">
                Why link your Telegram account?
              </h3>
              <div className="mt-2 text-sm text-blue-700">
                <ul className="list-disc list-inside space-y-1">
                  <li>Add bills and expenses from Telegram messages</li>
                  <li>View your spending on both platforms</li>
                  <li>Get quick summaries via Telegram bot</li>
                  <li>All your data synced in real-time</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
};
