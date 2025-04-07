const cookie = {
  key: {
    auth: 'auth',
  },
};

const locale = {
  languages: {
    en: 'en',
    // ja: 'ja',
  },
} as const;

export const CONSTANTS = {
  locale,
  cookie,
};
