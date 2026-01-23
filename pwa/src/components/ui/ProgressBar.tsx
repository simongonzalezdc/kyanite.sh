"use client";

interface ProgressBarProps {
  /** Progress value from 0 to 100 */
  value: number;
  /** Optional label to show above the progress bar */
  label?: string;
  /** Show percentage text */
  showPercentage?: boolean;
  /** Size variant */
  size?: "sm" | "md" | "lg";
  /** Color variant */
  variant?: "primary" | "success" | "warning" | "error";
  /** Indeterminate/loading state (animates) */
  indeterminate?: boolean;
  /** Accessible label for screen readers */
  "aria-label"?: string;
}

/** Width of the indeterminate progress indicator (%) */
const INDETERMINATE_PROGRESS_WIDTH = "30%";

export function ProgressBar({
  value,
  label,
  showPercentage = false,
  size = "md",
  variant = "primary",
  indeterminate = false,
  "aria-label": ariaLabel,
}: ProgressBarProps) {
  // Clamp value between 0 and 100
  const clampedValue = Math.min(100, Math.max(0, value));

  const heights = {
    sm: "h-1",
    md: "h-2",
    lg: "h-3",
  };

  const colors = {
    primary: "bg-[var(--color-primary)]",
    success: "bg-[var(--color-success)]",
    warning: "bg-[var(--color-warning)]",
    error: "bg-[var(--color-error)]",
  };

  return (
    <div className="w-full">
      {/* Label and percentage */}
      {(label || showPercentage) && (
        <div className="flex justify-between items-center mb-1">
          {label && (
            <span className="text-sm text-[var(--color-text-muted)]">{label}</span>
          )}
          {showPercentage && !indeterminate && (
            <span className="text-sm text-[var(--color-text-muted)]">
              {Math.round(clampedValue)}%
            </span>
          )}
        </div>
      )}

      {/* Progress bar track */}
      <div
        role="progressbar"
        aria-valuenow={indeterminate ? undefined : clampedValue}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-label={ariaLabel ?? label ?? "Progress"}
        className={`w-full rounded-full overflow-hidden bg-[var(--color-surface)] ${heights[size]}`}
      >
        {/* Progress bar fill */}
        <div
          className={`
            ${heights[size]}
            ${colors[variant]}
            rounded-full
            transition-all duration-300 ease-out
            ${indeterminate ? "animate-progress-indeterminate" : ""}
          `}
          style={{
            width: indeterminate ? INDETERMINATE_PROGRESS_WIDTH : `${clampedValue}%`,
          }}
        />
      </div>
    </div>
  );
}

/**
 * Circular progress indicator for compact spaces
 */
interface CircularProgressProps {
  /** Progress value from 0 to 100 */
  value: number;
  /** Size in pixels */
  size?: number;
  /** Stroke width */
  strokeWidth?: number;
  /** Show percentage in center */
  showPercentage?: boolean;
  /** Indeterminate state */
  indeterminate?: boolean;
  /** Accessible label */
  "aria-label"?: string;
}

export function CircularProgress({
  value,
  size = 40,
  strokeWidth = 4,
  showPercentage = false,
  indeterminate = false,
  "aria-label": ariaLabel,
}: CircularProgressProps) {
  const clampedValue = Math.min(100, Math.max(0, value));
  const radius = (size - strokeWidth) / 2;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (clampedValue / 100) * circumference;

  return (
    <div
      role="progressbar"
      aria-valuenow={indeterminate ? undefined : clampedValue}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label={ariaLabel ?? "Progress"}
      className="relative inline-flex items-center justify-center"
      style={{ width: size, height: size }}
    >
      <svg
        className={indeterminate ? "animate-spin" : ""}
        width={size}
        height={size}
        style={{ transform: "rotate(-90deg)" }}
      >
        {/* Background circle */}
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="var(--color-surface)"
          strokeWidth={strokeWidth}
        />
        {/* Progress circle */}
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="var(--color-primary)"
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={indeterminate ? circumference * 0.75 : offset}
          className="transition-all duration-300 ease-out"
        />
      </svg>
      {showPercentage && !indeterminate && (
        <span className="absolute text-xs font-medium text-[var(--color-text)]">
          {Math.round(clampedValue)}
        </span>
      )}
    </div>
  );
}
