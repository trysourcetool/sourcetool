export const CLOUD_DOMAIN = "trysourcetool.com";
export const LOCAL_DOMAIN = "local.trysourcetool.com";
export const LOCAL_API_BASE_URL = "local.trysourcetool.com:8080";

export const ENVIRONMENTS = {
  IS_CLOUD_EDITION: import.meta.env.VITE_IS_CLOUD_EDITION === "true",
  DOMAIN: import.meta.env.VITE_IS_CLOUD_EDITION === "true" ? CLOUD_DOMAIN : LOCAL_DOMAIN,
  API_BASE_URL: import.meta.env.VITE_IS_CLOUD_EDITION === "true" ? CLOUD_DOMAIN : LOCAL_API_BASE_URL,
  MODE: import.meta.env.MODE,
};
