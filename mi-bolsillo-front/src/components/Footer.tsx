import { useTranslation } from 'react-i18next';

export const Footer = () => {
  const { t } = useTranslation();
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-white border-t border-gray-200 mt-auto">
      <div className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {/* About Section */}
          <div>
            <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">
              {t('footer.about')}
            </h3>
            <p className="text-sm text-gray-600">
              {t('footer.aboutText')}
            </p>
          </div>

          {/* Quick Links */}
          <div>
            <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">
              {t('footer.quickLinks')}
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href="/"
                  className="text-sm text-gray-600 hover:text-blue-600 transition-colors"
                >
                  {t('footer.home')}
                </a>
              </li>
              <li>
                <a
                  href="/statistics"
                  className="text-sm text-gray-600 hover:text-blue-600 transition-colors"
                >
                  {t('footer.statistics')}
                </a>
              </li>
              <li>
                <a
                  href="/bills/new"
                  className="text-sm text-gray-600 hover:text-blue-600 transition-colors"
                >
                  {t('footer.createBill')}
                </a>
              </li>
              <li>
                <a
                  href="/link-telegram"
                  className="text-sm text-gray-600 hover:text-blue-600 transition-colors"
                >
                  {t('footer.linkTelegram')}
                </a>
              </li>
            </ul>
          </div>

          {/* Support & Resources */}
          <div>
            <h3 className="text-sm font-semibold text-gray-900 uppercase tracking-wider mb-4">
              {t('footer.support')}
            </h3>
            <ul className="space-y-2">
              <li>
                <a
                  href={`https://t.me/${import.meta.env.VITE_TELEGRAM_BOT_USERNAME || 'mi_bolsillo_bot'}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-gray-600 hover:text-blue-600 transition-colors flex items-center gap-1"
                >
                  <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 0C5.373 0 0 5.373 0 12s5.373 12 12 12 12-5.373 12-12S18.627 0 12 0zm5.894 8.221l-1.97 9.28c-.145.658-.537.818-1.084.508l-3-2.21-1.446 1.394c-.14.18-.357.223-.548.223l.188-2.85 5.18-4.68c.223-.198-.054-.308-.346-.11l-6.4 4.03-2.76-.918c-.6-.187-.612-.6.125-.89l10.782-4.156c.498-.187.935.112.77.89z"/>
                  </svg>
                  {t('footer.telegramBot')}
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/KKogaa/mi-bolsillo"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-sm text-gray-600 hover:text-blue-600 transition-colors flex items-center gap-1"
                >
                  <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                    <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
                  </svg>
                  {t('footer.github')}
                </a>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="mt-8 pt-8 border-t border-gray-200">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <p className="text-sm text-gray-500">
              © {currentYear} Mi Bolsillo. {t('footer.rights')}.
            </p>
            <div className="flex items-center gap-6 mt-4 md:mt-0">
              <a
                href="#"
                className="text-sm text-gray-500 hover:text-gray-700 transition-colors"
              >
                {t('footer.privacy')}
              </a>
              <a
                href="#"
                className="text-sm text-gray-500 hover:text-gray-700 transition-colors"
              >
                {t('footer.terms')}
              </a>
            </div>
          </div>
        </div>
      </div>
    </footer>
  );
};
