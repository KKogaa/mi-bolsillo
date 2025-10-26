import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import en from './locales/en.json';
import es from './locales/es.json';

i18n
  .use(LanguageDetector) // Detects user language
  .use(initReactI18next) // Passes i18n down to react-i18next
  .init({
    resources: {
      en: {
        translation: en
      },
      es: {
        translation: es
      }
    },
    fallbackLng: 'en',

    detection: {
      // Order and from where user language should be detected
      // localStorage first to prioritize user's manual selection
      order: ['localStorage', 'cookie', 'navigator', 'htmlTag'],

      // Keys or params to lookup language from
      lookupQuerystring: 'lng',
      lookupCookie: 'i18next',
      lookupLocalStorage: 'i18nextLng',
      lookupFromPathIndex: 0,
      lookupFromSubdomainIndex: 0,

      // Cache user language on
      caches: ['localStorage', 'cookie'],
      excludeCacheFor: ['cimode'], // Languages to not persist (only on request)

      // Optional htmlTag with lang attribute
      htmlTag: document.documentElement,
    },

    interpolation: {
      escapeValue: false // React already safes from xss
    }
  });

export default i18n;
