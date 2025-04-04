export const ENVIRONMENTS = {
  IS_CLOUD_EDITION: import.meta.env.VITE_IS_CLOUD_EDITION === "true",
  DOMAIN: import.meta.env.VITE_IS_CLOUD_EDITION === "true" ? "trysourcetool.com" : "local.trysourcetool.com",
  API_BASE_URL: import.meta.env.VITE_IS_CLOUD_EDITION === "true" ? "trysourcetool.com" : "local.trysourcetool.com:8080",
  MODE: import.meta.env.MODE,
};
