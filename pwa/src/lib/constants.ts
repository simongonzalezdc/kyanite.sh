/**
 * Application-wide constants
 * Extract magic numbers here to improve code clarity and maintainability
 */

// =============================================================================
// Sync Configuration
// =============================================================================

/** Default interval for automatic sync polling (ms) */
export const SYNC_INTERVAL_MS = 30000;

/** Timeout for API ping requests (ms) */
export const API_PING_TIMEOUT_MS = 5000;

// =============================================================================
// WebSocket Configuration
// =============================================================================

/** Base delay before attempting WebSocket reconnection (ms) */
export const WEBSOCKET_RECONNECT_DELAY_MS = 1000;

/** Maximum number of WebSocket reconnection attempts */
export const WEBSOCKET_MAX_RECONNECT_ATTEMPTS = 5;

/** Interval between WebSocket ping messages (ms) */
export const WEBSOCKET_PING_INTERVAL_MS = 30000;

/** Exponential backoff multiplier for reconnection attempts */
export const WEBSOCKET_BACKOFF_MULTIPLIER = 2;

// =============================================================================
// Camera/Photo Configuration
// =============================================================================

/** Ideal camera capture width (pixels) */
export const CAMERA_WIDTH = 1920;

/** Ideal camera capture height (pixels) */
export const CAMERA_HEIGHT = 1080;

/** JPEG quality for captured photos (0-1) */
export const JPEG_QUALITY = 0.9;

// =============================================================================
// Audio/Recording Configuration
// =============================================================================

/** MediaRecorder timeslice - how often to collect data (ms) */
export const MEDIA_RECORDER_TIMESLICE_MS = 100;

/** Interval for updating recording duration display (ms) */
export const RECORDING_TIMER_INTERVAL_MS = 1000;

// =============================================================================
// Tap Tempo Configuration
// =============================================================================

/** Minimum number of taps required for BPM calculation */
export const TAP_TEMPO_MIN_TAPS = 4;

/** Maximum number of taps to store for BPM calculation */
export const TAP_TEMPO_MAX_TAPS = 16;

/** Timeout before tap sequence resets (ms) */
export const TAP_TEMPO_TIMEOUT_MS = 2000;

/** Minimum valid BPM value */
export const BPM_MIN = 20;

/** Maximum valid BPM value */
export const BPM_MAX = 300;

/** Milliseconds per minute (for BPM calculations) */
export const MS_PER_MINUTE = 60000;

/** Haptic feedback duration for tap (ms) */
export const TAP_HAPTIC_DURATION_MS = 10;

// =============================================================================
// Time Display Thresholds
// =============================================================================

/** Threshold for showing "just now" (seconds) */
export const TIME_JUST_NOW_THRESHOLD = 60;

/** Threshold for showing "1m ago" (seconds) */
export const TIME_ONE_MINUTE_THRESHOLD = 120;

/** Threshold for showing minutes (seconds) */
export const TIME_ONE_HOUR_THRESHOLD = 3600;

/** Threshold for showing "1h ago" (seconds) */
export const TIME_TWO_HOURS_THRESHOLD = 7200;

/** Seconds per minute */
export const SECONDS_PER_MINUTE = 60;

/** Seconds per hour */
export const SECONDS_PER_HOUR = 3600;
