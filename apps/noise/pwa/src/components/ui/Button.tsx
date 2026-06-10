"use client";

import { forwardRef, type ButtonHTMLAttributes } from "react";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className = "", variant = "primary", size = "md", loading, disabled, children, ...props }, ref) => {
    const baseStyles = "inline-flex items-center justify-center font-mono font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50";
    
    const variants = {
      primary: "bg-[var(--color-primary)] text-[var(--color-background)] hover:opacity-90 focus-visible:ring-[var(--color-primary)]",
      secondary: "border border-[var(--color-secondary)] text-[var(--color-text)] hover:bg-[var(--color-surface)] focus-visible:ring-[var(--color-secondary)]",
      ghost: "text-[var(--color-text)] hover:bg-[var(--color-surface)] focus-visible:ring-[var(--color-text-muted)]",
      danger: "bg-[var(--color-error)] text-[var(--color-background)] hover:opacity-90 focus-visible:ring-[var(--color-error)]",
    };

    // All sizes meet 44px minimum touch target (WCAG 2.5.5)
    const sizes = {
      sm: "h-11 min-h-[44px] px-4 text-sm rounded",
      md: "h-12 min-h-[44px] px-5 text-base rounded-md",
      lg: "h-14 min-h-[44px] px-6 text-lg rounded-lg",
    };

    return (
      <button
        ref={ref}
        className={`${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`}
        disabled={disabled || loading}
        {...props}
      >
        {loading ? (
          <>
            <svg className="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
            Loading...
          </>
        ) : (
          children
        )}
      </button>
    );
  }
);

Button.displayName = "Button";
