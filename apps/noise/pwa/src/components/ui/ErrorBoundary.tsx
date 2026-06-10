"use client";

import React, { Component, ErrorInfo, ReactNode } from "react";

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

/**
 * ErrorBoundary catches JavaScript errors anywhere in the child component tree,
 * logs those errors, and displays a fallback UI instead of crashing the app.
 */
export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo): void {
    // Log error to console for debugging
    console.error("[ErrorBoundary] Caught error:", error);
    console.error("[ErrorBoundary] Error info:", errorInfo);
  }

  private handleReset = (): void => {
    this.setState({ hasError: false, error: null });
  };

  render(): ReactNode {
    if (this.state.hasError) {
      // Custom fallback UI if provided
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // Default fallback UI
      return (
        <div
          role="alert"
          className="min-h-screen flex items-center justify-center p-4"
          style={{ backgroundColor: "var(--background)", color: "var(--foreground)" }}
        >
          <div className="max-w-md w-full text-center space-y-4">
            <h1 className="text-xl font-semibold">Something went wrong</h1>
            <p className="text-sm opacity-70">
              The app encountered an unexpected error. Please try refreshing the page.
            </p>
            {this.state.error && (
              <details className="text-left text-xs opacity-50 mt-4">
                <summary className="cursor-pointer">Technical details</summary>
                <pre className="mt-2 p-2 rounded overflow-auto max-h-32" style={{ backgroundColor: "var(--muted)" }}>
                  {this.state.error.message}
                </pre>
              </details>
            )}
            <div className="flex gap-2 justify-center pt-4">
              <button
                onClick={this.handleReset}
                className="px-4 py-2 rounded text-sm font-medium"
                style={{
                  backgroundColor: "var(--primary)",
                  color: "var(--primary-foreground)",
                }}
              >
                Try again
              </button>
              <button
                onClick={() => window.location.reload()}
                className="px-4 py-2 rounded text-sm font-medium"
                style={{
                  backgroundColor: "var(--muted)",
                  color: "var(--foreground)",
                }}
              >
                Refresh page
              </button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
