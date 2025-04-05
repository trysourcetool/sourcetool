/**
 * Log levels
 */
export enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
  FATAL = 'fatal',
}

/**
 * Current log level
 */
let currentLevel: LogLevel = LogLevel.INFO;

/**
 * Set the log level
 * @param level Log level
 */
export function setLevel(level: LogLevel): void {
  currentLevel = level;
}

/**
 * Check if the message should be logged
 * @param level Log level
 * @returns True if the message should be logged
 */
function shouldLog(level: LogLevel): boolean {
  const levels = [
    LogLevel.DEBUG,
    LogLevel.INFO,
    LogLevel.WARN,
    LogLevel.ERROR,
    LogLevel.FATAL,
  ];
  const currentIndex = levels.indexOf(currentLevel);
  const messageIndex = levels.indexOf(level);

  return messageIndex >= currentIndex;
}

/**
 * Log a debug message
 * @param message Message
 * @param args Arguments
 */
export function debug(message: string, ...args: any[]): void {
  if (shouldLog(LogLevel.DEBUG)) {
    console.debug(`[DEBUG] ${message}`, ...args);
  }
}

/**
 * Log an info message
 * @param message Message
 * @param args Arguments
 */
export function info(message: string, ...args: any[]): void {
  if (shouldLog(LogLevel.INFO)) {
    console.info(`[INFO] ${message}`, ...args);
  }
}

/**
 * Log a warning message
 * @param message Message
 * @param args Arguments
 */
export function warn(message: string, ...args: any[]): void {
  if (shouldLog(LogLevel.WARN)) {
    console.warn(`[WARN] ${message}`, ...args);
  }
}

/**
 * Log an error message
 * @param message Message
 * @param args Arguments
 */
export function error(message: string, ...args: any[]): void {
  if (shouldLog(LogLevel.ERROR)) {
    console.error(`[ERROR] ${message}`, ...args);
  }
}

/**
 * Log a fatal message
 * @param message Message
 * @param args Arguments
 */
export function fatal(message: string, ...args: any[]): void {
  if (shouldLog(LogLevel.FATAL)) {
    console.error(`[FATAL] ${message}`, ...args);
  }
}

/**
 * Initialize the logger
 */
export function init(): void {
  // Set log level from environment variable
  const logLevel = process.env.LOG_LEVEL as LogLevel;
  if (logLevel && Object.values(LogLevel).includes(logLevel)) {
    setLevel(logLevel);
  }
}

/**
 * Sync the logger
 */
export function sync(): void {
  // No-op in this implementation
}
