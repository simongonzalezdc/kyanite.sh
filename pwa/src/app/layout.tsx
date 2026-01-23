import type { Metadata, Viewport } from "next";
import { SkipLink } from "@/components/ui/SkipLink";
import "./globals.css";

export const metadata: Metadata = {
  title: "noise.sh Companion",
  description: "Capture song ideas on the go - voice memos, photos, tap tempo, and quick notes that sync to your desktop.",
  appleWebApp: {
    capable: true,
    statusBarStyle: "black-translucent",
    title: "noise.sh",
  },
  formatDetection: {
    telephone: false,
  },
  other: {
    "mobile-web-app-capable": "yes",
  },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
  viewportFit: "cover",
  themeColor: "#D4A574",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        {/* PWA meta tags */}
        <link rel="apple-touch-icon" href="/icons/icon-192.svg" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent" />
      </head>
      <body className="antialiased">
        <SkipLink />
        <div id="main-content">
          {children}
        </div>
      </body>
    </html>
  );
}
