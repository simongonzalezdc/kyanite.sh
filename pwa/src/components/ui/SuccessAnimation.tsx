"use client";

import { useEffect, useRef } from "react";

/** Duration of the success animation (ms) */
const SUCCESS_ANIMATION_DURATION_MS = 1500;

interface SuccessAnimationProps {
  /** Whether to show the animation */
  show: boolean;
  /** Callback when animation completes */
  onComplete?: () => void;
  /** Size of the animation */
  size?: "sm" | "md" | "lg";
  /** Text to show below the checkmark */
  message?: string;
}

/**
 * Animated checkmark success indicator
 */
export function SuccessAnimation({
  show,
  onComplete,
  size = "md",
  message,
}: SuccessAnimationProps) {
  const prevShowRef = useRef(show);

  // Trigger onComplete callback after animation duration
  useEffect(() => {
    if (show && !prevShowRef.current) {
      // Transition from false to true - animation started
      const timer = setTimeout(() => {
        onComplete?.();
      }, SUCCESS_ANIMATION_DURATION_MS);

      prevShowRef.current = show;
      return () => clearTimeout(timer);
    }
    prevShowRef.current = show;
  }, [show, onComplete]);

  // Only render when show is true
  if (!show) {
    return null;
  }

  // Animation is controlled by CSS - plays when component mounts (show becomes true)
  const isAnimating = true;

  const sizes = {
    sm: { circle: 48, stroke: 3, icon: 20 },
    md: { circle: 72, stroke: 4, icon: 32 },
    lg: { circle: 96, stroke: 5, icon: 44 },
  };

  const s = sizes[size];
  const radius = (s.circle - s.stroke * 2) / 2;
  const circumference = 2 * Math.PI * radius;

  return (
    <div
      className={`
        fixed inset-0 z-50 flex flex-col items-center justify-center
        bg-[var(--color-background)]/80 backdrop-blur-sm
        transition-opacity duration-300
        ${show ? "opacity-100" : "opacity-0"}
      `}
      role="status"
      aria-label={message ?? "Success"}
    >
      <div className="relative" style={{ width: s.circle, height: s.circle }}>
        {/* Circle */}
        <svg
          className="absolute inset-0"
          width={s.circle}
          height={s.circle}
          viewBox={`0 0 ${s.circle} ${s.circle}`}
        >
          {/* Background circle */}
          <circle
            cx={s.circle / 2}
            cy={s.circle / 2}
            r={radius}
            fill="none"
            stroke="var(--color-surface)"
            strokeWidth={s.stroke}
          />
          {/* Animated circle */}
          <circle
            cx={s.circle / 2}
            cy={s.circle / 2}
            r={radius}
            fill="none"
            stroke="var(--color-success)"
            strokeWidth={s.stroke}
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={circumference}
            style={{
              animation: isAnimating ? "success-circle 0.6s ease-out forwards" : "none",
              transformOrigin: "center",
              transform: "rotate(-90deg)",
            }}
          />
        </svg>

        {/* Checkmark */}
        <svg
          className="absolute inset-0 m-auto"
          width={s.icon}
          height={s.icon}
          viewBox="0 0 24 24"
          fill="none"
          stroke="var(--color-success)"
          strokeWidth="3"
          strokeLinecap="round"
          strokeLinejoin="round"
          style={{
            animation: isAnimating ? "success-check 0.4s ease-out 0.4s forwards" : "none",
            opacity: 0,
          }}
        >
          <path d="M5 12l5 5L20 7" />
        </svg>
      </div>

      {message && (
        <p
          className="mt-4 text-lg font-medium text-[var(--color-text)]"
          style={{
            animation: isAnimating ? "success-text 0.3s ease-out 0.6s forwards" : "none",
            opacity: 0,
          }}
        >
          {message}
        </p>
      )}

      <style jsx>{`
        @keyframes success-circle {
          to {
            stroke-dashoffset: 0;
          }
        }

        @keyframes success-check {
          0% {
            opacity: 0;
            transform: scale(0.5);
          }
          50% {
            transform: scale(1.2);
          }
          100% {
            opacity: 1;
            transform: scale(1);
          }
        }

        @keyframes success-text {
          from {
            opacity: 0;
            transform: translateY(10px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
      `}</style>
    </div>
  );
}

/**
 * Compact inline success indicator (no overlay)
 */
interface InlineSuccessProps {
  show: boolean;
  className?: string;
}

export function InlineSuccess({ show, className = "" }: InlineSuccessProps) {
  if (!show) return null;

  return (
    <span
      className={`inline-flex items-center gap-1 text-[var(--color-success)] ${className}`}
      style={{ animation: "inline-success 0.3s ease-out" }}
    >
      <svg
        className="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M5 13l4 4L19 7"
        />
      </svg>

      <style jsx>{`
        @keyframes inline-success {
          0% {
            opacity: 0;
            transform: scale(0.8);
          }
          50% {
            transform: scale(1.1);
          }
          100% {
            opacity: 1;
            transform: scale(1);
          }
        }
      `}</style>
    </span>
  );
}
