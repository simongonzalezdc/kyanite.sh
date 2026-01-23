"use client";

import { useEffect, useState, useCallback, createContext, useContext, ReactNode } from "react";

// =============================================================================
// Types
// =============================================================================

export type ToastType = "success" | "error" | "warning" | "info";

export interface Toast {
  id: string;
  type: ToastType;
  message: string;
  duration?: number;
  action?: {
    label: string;
    onClick: () => void;
  };
}

interface ToastContextValue {
  toasts: Toast[];
  addToast: (toast: Omit<Toast, "id">) => string;
  removeToast: (id: string) => void;
  clearAll: () => void;
}

// =============================================================================
// Constants
// =============================================================================

const DEFAULT_DURATION_MS = 4000;
const ANIMATION_DURATION_MS = 300;

const TOAST_ICONS: Record<ToastType, string> = {
  success: "✓",
  error: "✕",
  warning: "⚠",
  info: "ℹ",
};

const TOAST_COLORS: Record<ToastType, { bg: string; border: string; icon: string }> = {
  success: {
    bg: "bg-[var(--color-success)]/10",
    border: "border-[var(--color-success)]",
    icon: "text-[var(--color-success)]",
  },
  error: {
    bg: "bg-[var(--color-error)]/10",
    border: "border-[var(--color-error)]",
    icon: "text-[var(--color-error)]",
  },
  warning: {
    bg: "bg-[var(--color-warning)]/10",
    border: "border-[var(--color-warning)]",
    icon: "text-[var(--color-warning)]",
  },
  info: {
    bg: "bg-[var(--color-primary)]/10",
    border: "border-[var(--color-primary)]",
    icon: "text-[var(--color-primary)]",
  },
};

// =============================================================================
// Context
// =============================================================================

const ToastContext = createContext<ToastContextValue | null>(null);

export function useToast(): ToastContextValue {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return context;
}

// =============================================================================
// Provider
// =============================================================================

interface ToastProviderProps {
  children: ReactNode;
}

export function ToastProvider({ children }: ToastProviderProps) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback((toast: Omit<Toast, "id">): string => {
    const id = `toast-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    const newToast: Toast = {
      ...toast,
      id,
      duration: toast.duration ?? DEFAULT_DURATION_MS,
    };

    setToasts((prev) => [...prev, newToast]);
    return id;
  }, []);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const clearAll = useCallback(() => {
    setToasts([]);
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast, clearAll }}>
      {children}
      <ToastContainer toasts={toasts} onRemove={removeToast} />
    </ToastContext.Provider>
  );
}

// =============================================================================
// Toast Container
// =============================================================================

interface ToastContainerProps {
  toasts: Toast[];
  onRemove: (id: string) => void;
}

function ToastContainer({ toasts, onRemove }: ToastContainerProps) {
  return (
    <div
      aria-live="polite"
      aria-label="Notifications"
      className="fixed bottom-4 left-4 right-4 z-50 flex flex-col gap-2 pointer-events-none"
      style={{ paddingBottom: "env(safe-area-inset-bottom)" }}
    >
      {toasts.map((toast) => (
        <ToastItem key={toast.id} toast={toast} onRemove={onRemove} />
      ))}
    </div>
  );
}

// =============================================================================
// Toast Item
// =============================================================================

interface ToastItemProps {
  toast: Toast;
  onRemove: (id: string) => void;
}

function ToastItem({ toast, onRemove }: ToastItemProps) {
  const [isVisible, setIsVisible] = useState(false);
  const [isExiting, setIsExiting] = useState(false);

  const colors = TOAST_COLORS[toast.type];
  const icon = TOAST_ICONS[toast.type];

  // Define dismiss handler first
  const handleDismiss = useCallback(() => {
    setIsExiting(true);
    setTimeout(() => {
      onRemove(toast.id);
    }, ANIMATION_DURATION_MS);
  }, [onRemove, toast.id]);

  const handleActionClick = useCallback(() => {
    toast.action?.onClick();
    handleDismiss();
  }, [toast.action, handleDismiss]);

  // Animate in on mount
  useEffect(() => {
    const timer = setTimeout(() => setIsVisible(true), 10);
    return () => clearTimeout(timer);
  }, []);

  // Auto-dismiss after duration
  useEffect(() => {
    if (toast.duration && toast.duration > 0) {
      const timer = setTimeout(() => {
        handleDismiss();
      }, toast.duration);
      return () => clearTimeout(timer);
    }
  }, [toast.duration, handleDismiss]);

  return (
    <div
      role="alert"
      className={`
        pointer-events-auto
        flex items-center gap-3
        px-4 py-3
        rounded-lg
        border
        backdrop-blur-sm
        shadow-lg
        transition-all duration-300 ease-out
        ${colors.bg}
        ${colors.border}
        ${isVisible && !isExiting ? "opacity-100 translate-y-0" : "opacity-0 translate-y-4"}
      `}
      style={{ backgroundColor: "var(--color-surface)" }}
    >
      {/* Icon */}
      <span className={`text-lg font-bold ${colors.icon}`} aria-hidden="true">
        {icon}
      </span>

      {/* Message */}
      <p className="flex-1 text-sm text-[var(--color-text)]">{toast.message}</p>

      {/* Action button */}
      {toast.action && (
        <button
          onClick={handleActionClick}
          className="text-sm font-medium text-[var(--color-primary)] hover:underline focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] rounded px-2 py-1 min-h-[44px] min-w-[44px] flex items-center justify-center"
        >
          {toast.action.label}
        </button>
      )}

      {/* Dismiss button */}
      <button
        onClick={handleDismiss}
        className="text-[var(--color-text-muted)] hover:text-[var(--color-text)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] rounded min-h-[44px] min-w-[44px] flex items-center justify-center"
        aria-label="Dismiss notification"
      >
        <span aria-hidden="true">×</span>
      </button>
    </div>
  );
}

// =============================================================================
// Convenience hooks for common toast types
// =============================================================================

export function useToastActions() {
  const { addToast } = useToast();

  return {
    success: (message: string, action?: Toast["action"]) =>
      addToast({ type: "success", message, action }),

    error: (message: string, action?: Toast["action"]) =>
      addToast({ type: "error", message, duration: 6000, action }),

    warning: (message: string, action?: Toast["action"]) =>
      addToast({ type: "warning", message, action }),

    info: (message: string, action?: Toast["action"]) =>
      addToast({ type: "info", message, action }),
  };
}
