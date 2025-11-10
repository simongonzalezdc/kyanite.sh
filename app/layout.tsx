import type { Metadata } from 'next'
import './globals.css'

export const metadata: Metadata = {
  title: 'VoxForge - Transform Your Voice into Music',
  description: 'Browser-based voice-to-song tool that transforms vocal recordings into complete music arrangements',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  )
}

