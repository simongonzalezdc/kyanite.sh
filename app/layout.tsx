import type { Metadata, Viewport } from 'next'
import './globals.css'
import { AnalyticsProvider } from './components/AnalyticsProvider'
import { AccessibilityProvider } from './components/AccessibilityProvider'
import { KeyboardNavigation } from './components/KeyboardNavigation'
import { ScreenReaderNavigation } from './components/ScreenReaderSupport'
import { HighContrastProvider, FontSizeProvider } from './components/VisualAccessibility'
import { VisualAccessibilityToolbar } from './components/VisualAccessibility'

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
}

export const metadata: Metadata = {
  title: 'VoxForge - Transform Your Voice into Music',
  description: 'Browser-based voice-to-song tool that transforms vocal recordings into complete music arrangements with full accessibility support',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'black-translucent',
    title: 'VoxForge'
  },
  manifest: '/manifest.json',
  robots: {
    index: true,
    follow: true
  },
  keywords: ['music', 'voice', 'accessibility', 'WCAG', 'screen reader', 'keyboard navigation']
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <head>
        <meta name="mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
        <meta name="format-detection" content="telephone=no" />
        <meta name="msapplication-tap-highlight" content="no" />
        <meta name="theme-color" content="#0a0a0a" />
        <meta name="description" content="VoxForge - Accessible voice-to-music transformation tool with WCAG 2.1 AA compliance" />
        <meta name="author" content="VoxForge Team" />
        <meta name="accessibility" content="WCAG 2.1 AA compliant" />
      </head>
      <body className="antialiased">
        <AccessibilityProvider>
          <HighContrastProvider>
            <FontSizeProvider>
              <KeyboardNavigation>
                <ScreenReaderNavigation />
                <AnalyticsProvider>
                  <main role="application" aria-label="VoxForge Music Application">
                    {children}
                  </main>
                  <VisualAccessibilityToolbar />
                </AnalyticsProvider>
              </KeyboardNavigation>
            </FontSizeProvider>
          </HighContrastProvider>
        </AccessibilityProvider>
      </body>
    </html>
  )
}

