import type { ReactNode } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { UserButton, useUser } from '@clerk/clerk-react';
import { useTranslation } from 'react-i18next';
import { PocketLogo } from '../components/PocketLogo';
import { Footer } from '../components/Footer';
import { LanguageToggle } from '../components/LanguageToggle';

interface MainLayoutProps {
  children: ReactNode;
}

export const MainLayout = ({ children }: MainLayoutProps) => {
  const { user } = useUser();
  const location = useLocation();
  const { t } = useTranslation();

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex items-center">
              <Link to="/" className="flex items-center gap-3 hover:opacity-80 transition-opacity">
                <PocketLogo className="text-blue-600" size={36} />
                <h1 className="text-xl font-bold text-gray-900">Mi Bolsillo</h1>
              </Link>
              <div className="hidden md:flex items-center ml-10 space-x-8">
                <Link
                  to="/"
                  className={`text-sm font-medium transition-colors ${
                    isActive('/')
                      ? 'text-blue-600 border-b-2 border-blue-600 pb-[22px]'
                      : 'text-gray-700 hover:text-blue-600'
                  }`}
                >
                  {t('nav.home')}
                </Link>
                <Link
                  to="/statistics"
                  className={`text-sm font-medium transition-colors ${
                    isActive('/statistics')
                      ? 'text-blue-600 border-b-2 border-blue-600 pb-[22px]'
                      : 'text-gray-700 hover:text-blue-600'
                  }`}
                >
                  {t('nav.statistics')}
                </Link>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <LanguageToggle />
              <span className="text-sm text-gray-600">
                {user?.firstName || user?.emailAddresses[0]?.emailAddress}
              </span>
              <UserButton afterSignOutUrl="/login" />
            </div>
          </div>
        </div>
      </nav>
      <main className="flex-grow max-w-7xl w-full mx-auto py-6 sm:px-6 lg:px-8">
        {children}
      </main>
      <Footer />
    </div>
  );
};
