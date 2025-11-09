import type { ReactNode } from 'react';
import { useState } from 'react';
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
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen(!isMobileMenuOpen);
  };

  const closeMobileMenu = () => {
    setIsMobileMenuOpen(false);
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
              {/* Desktop Navigation */}
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
              <span className="hidden sm:inline text-sm text-gray-600">
                {user?.firstName || user?.emailAddresses[0]?.emailAddress}
              </span>
              <UserButton afterSignOutUrl="/login" />
              {/* Mobile menu button */}
              <button
                onClick={toggleMobileMenu}
                className="md:hidden inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:text-blue-600 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-600"
                aria-expanded={isMobileMenuOpen}
              >
                <span className="sr-only">Open main menu</span>
                {isMobileMenuOpen ? (
                  <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                ) : (
                  <svg className="block h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                  </svg>
                )}
              </button>
            </div>
          </div>
          {/* Mobile menu */}
          {isMobileMenuOpen && (
            <div className="md:hidden">
              <div className="px-2 pt-2 pb-3 space-y-1">
                <Link
                  to="/"
                  onClick={closeMobileMenu}
                  className={`block px-3 py-2 rounded-md text-base font-medium ${
                    isActive('/')
                      ? 'text-blue-600 bg-blue-50'
                      : 'text-gray-700 hover:text-blue-600 hover:bg-gray-50'
                  }`}
                >
                  {t('nav.home')}
                </Link>
                <Link
                  to="/statistics"
                  onClick={closeMobileMenu}
                  className={`block px-3 py-2 rounded-md text-base font-medium ${
                    isActive('/statistics')
                      ? 'text-blue-600 bg-blue-50'
                      : 'text-gray-700 hover:text-blue-600 hover:bg-gray-50'
                  }`}
                >
                  {t('nav.statistics')}
                </Link>
              </div>
              {/* Mobile user info */}
              <div className="pt-4 pb-3 border-t border-gray-200">
                <div className="px-5">
                  <div className="text-base font-medium text-gray-800">
                    {user?.firstName || 'User'}
                  </div>
                  <div className="text-sm font-medium text-gray-500">
                    {user?.emailAddresses[0]?.emailAddress}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </nav>
      <main className="flex-grow max-w-7xl w-full mx-auto py-6 sm:px-6 lg:px-8">
        {children}
      </main>
      <Footer />
    </div>
  );
};
