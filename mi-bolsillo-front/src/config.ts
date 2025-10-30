// Runtime configuration helper
// This allows environment variables to be injected at runtime (Docker) or build time (Vite)

interface RuntimeConfig {
  VITE_CLERK_PUBLISHABLE_KEY?: string;
  VITE_API_URL?: string;
  VITE_TELEGRAM_BOT_USERNAME?: string;
}

declare global {
  interface Window {
    RUNTIME_CONFIG?: RuntimeConfig;
  }
}

// Helper function to get environment variables from runtime config or build-time env
export const getEnv = (key: keyof RuntimeConfig): string | undefined => {
  return window.RUNTIME_CONFIG?.[key] || import.meta.env[key];
};

// Export specific config values
export const config = {
  clerkPublishableKey: getEnv('VITE_CLERK_PUBLISHABLE_KEY'),
  apiUrl: getEnv('VITE_API_URL') || 'http://localhost:8080',
  telegramBotUsername: getEnv('VITE_TELEGRAM_BOT_USERNAME') || 'mi_bolsillo_bot',
};
