import 'i18next';
import commonEn from '../locales/en/common.json';

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'common';
    resources: {
      common: typeof commonEn;
    };
  }
}
