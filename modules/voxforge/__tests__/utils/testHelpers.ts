/**
 * Test helpers for common VoxForge testing scenarios
 */

import { render, RenderOptions } from '@testing-library/react'
import React, { ReactElement, ReactNode } from 'react'
import { StoreProvider } from '@/app/components/StoreProvider'
import { AccessibilityProvider } from '@/app/components/AccessibilityProvider'
import { KeyboardNavigation } from '@/app/components/KeyboardNavigation'
import { setupAudioMocks } from './audioMocks'

declare const jest: any

// Custom render function with providers
export const renderWithProviders = (
  ui: ReactElement,
  options: RenderOptions = {}
) => {
  const Wrapper = ({ children }: { children: ReactNode }) => {
    return React.createElement(
      StoreProvider,
      null,
      React.createElement(
        AccessibilityProvider,
        null,
        React.createElement(KeyboardNavigation, null, children)
      )
    )
  }

  return render(ui, { wrapper: Wrapper, ...options })
}

// Helper to wait for async operations
export const waitForAsync = () => new Promise(resolve => setTimeout(resolve, 0))

// Helper to create a mock audio context with test data
export const createMockAudioContext = () => {
  setupAudioMocks()
  return new (global.AudioContext as any)()
}

// Helper to simulate user media permission
export const simulateMediaPermission = (granted: boolean) => {
  const mockGetUserMedia = granted
    ? jest.fn().mockResolvedValue({
        getTracks: () => [{
          stop: jest.fn(),
          getSettings: () => ({
            deviceId: 'default',
            groupId: 'default',
            kind: 'audioinput',
            label: 'Default',
            sampleRate: 44100,
          }),
        }],
      })
    : jest.fn().mockRejectedValue(new Error('Permission denied'))

  Object.defineProperty(global.navigator, 'mediaDevices', {
    value: {
      getUserMedia: mockGetUserMedia,
      enumerateDevices: jest.fn().mockResolvedValue([
        {
          deviceId: 'default',
          groupId: 'default',
          kind: 'audioinput',
          label: 'Default',
        },
      ]),
    },
    writable: true,
  })

  return mockGetUserMedia
}

// Helper to simulate audio recording
export const simulateAudioRecording = (duration: number = 1000) => {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve({
        data: new ArrayBuffer(duration),
        blob: new Blob([new ArrayBuffer(duration)], { type: 'audio/wav' }),
      })
    }, duration)
  })
}

// Helper to simulate audio analysis
export const simulateAudioAnalysis = () => {
  return {
    pitches: [
      { frequency: 440, time: 0, midi: 69, confidence: 0.8 },
      { frequency: 523.25, time: 0.5, midi: 72, confidence: 0.9 },
      { frequency: 659.25, time: 1, midi: 76, confidence: 0.7 },
    ],
    bpm: { bpm: 120, confidence: 0.8, stable: true },
    key: { key: 'C Major', tonic: 'C', scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'] },
    timeSignature: { numerator: 4, denominator: 4, display: '4/4' },
  }
}

// Helper to simulate file upload
export const simulateFileUpload = (fileName: string = 'test-audio.wav') => {
  const file = new File([''], fileName, { type: 'audio/wav' })
  return file
}

// Helper to simulate keyboard events
export const simulateKeyPress = (key: string, options: KeyboardEventInit = {}) => {
  const event = new KeyboardEvent('keydown', {
    key,
    bubbles: true,
    ...options,
  })
  document.dispatchEvent(event)
  return event
}

// Helper to simulate mouse events
export const simulateMouseClick = (element: HTMLElement) => {
  const event = new MouseEvent('click', {
    bubbles: true,
    cancelable: true,
    view: window,
  })
  element.dispatchEvent(event)
  return event
}

// Helper to simulate drag and drop
export const simulateDragDrop = (dragElement: HTMLElement, dropElement: HTMLElement) => {
  // Drag start
  const dragStartEvent = new DragEvent('dragstart', {
    bubbles: true,
    cancelable: true,
  })
  dragElement.dispatchEvent(dragStartEvent)

  // Drag over
  const dragOverEvent = new DragEvent('dragover', {
    bubbles: true,
    cancelable: true,
  })
  dropElement.dispatchEvent(dragOverEvent)

  // Drop
  const dropEvent = new DragEvent('drop', {
    bubbles: true,
    cancelable: true,
  })
  dropElement.dispatchEvent(dropEvent)

  // Drag end
  const dragEndEvent = new DragEvent('dragend', {
    bubbles: true,
    cancelable: true,
  })
  dragElement.dispatchEvent(dragEndEvent)
}

// Helper to simulate window resize
export const simulateResize = (width: number, height: number) => {
  Object.defineProperty(window, 'innerWidth', {
    writable: true,
    configurable: true,
    value: width,
  })
  Object.defineProperty(window, 'innerHeight', {
    writable: true,
    configurable: true,
    value: height,
  })

  const event = new Event('resize')
  window.dispatchEvent(event)
}

// Helper to simulate network conditions
export const simulateNetworkCondition = (condition: 'slow' | 'offline' | 'fast') => {
  const conditions = {
    slow: { effectiveType: 'slow-2g', downlink: 0.1, rtt: 2000 },
    offline: { effectiveType: 'slow-2g', downlink: 0, rtt: 9999, onLine: false },
    fast: { effectiveType: '4g', downlink: 10, rtt: 50 },
  }

  Object.defineProperty(global.navigator, 'connection', {
    value: conditions[condition],
    writable: true,
  })

  Object.defineProperty(global.navigator, 'onLine', {
    value: condition !== 'offline',
    writable: true,
  })
}

// Helper to simulate device orientation
export const simulateDeviceOrientation = (alpha: number, beta: number, gamma: number) => {
  const event = new DeviceOrientationEvent('deviceorientation', {
    alpha,
    beta,
    gamma,
  })
  window.dispatchEvent(event)
}

// Helper to simulate touch events
export const simulateTouch = (element: HTMLElement, type: 'start' | 'move' | 'end', x: number, y: number) => {
  const touch = new Touch({
    identifier: 0,
    target: element,
    clientX: x,
    clientY: y,
    pageX: x,
    pageY: y,
    screenX: x,
    screenY: y,
  })

  const event = new TouchEvent(`touch${type}`, {
    bubbles: true,
    cancelable: true,
    touches: type === 'end' ? [] : [touch],
    changedTouches: [touch],
  })

  element.dispatchEvent(event)
  return event
}

// Helper to test accessibility
export const testAccessibility = async (container: HTMLElement) => {
  // This would integrate with axe-core for accessibility testing
  // For now, just check for basic ARIA attributes
  const interactiveElements = container.querySelectorAll('button, input, select, textarea, a')
  const hasAriaLabels = Array.from(interactiveElements).every(el => 
    el.hasAttribute('aria-label') || 
    el.hasAttribute('aria-labelledby') || 
    el.textContent?.trim()
  )

  return {
    passed: hasAriaLabels,
    issues: hasAriaLabels ? [] : ['Some interactive elements lack proper ARIA labels'],
  }
}

// Helper to measure performance
export const measurePerformance = (fn: () => void | Promise<void>) => {
  const start = performance.now()
  const result = fn()
  
  if (result instanceof Promise) {
    return result.then(() => performance.now() - start)
  } else {
    return Promise.resolve(performance.now() - start)
  }
}

// Helper to create a mock fetch response
export const mockFetchResponse = (data: any, status: number = 200) => {
  return Promise.resolve({
    ok: status >= 200 && status < 300,
    status,
    json: () => Promise.resolve(data),
    text: () => Promise.resolve(JSON.stringify(data)),
    blob: () => Promise.resolve(new Blob([JSON.stringify(data)])),
  } as Response)
}

// Helper to mock console methods
export const mockConsole = () => {
  const originalConsole = { ...console }
  const mockConsole = {
    log: jest.fn(),
    warn: jest.fn(),
    error: jest.fn(),
    info: jest.fn(),
    debug: jest.fn(),
    trace: jest.fn(),
    assert: jest.fn(),
    clear: jest.fn(),
    count: jest.fn(),
    countReset: jest.fn(),
    dir: jest.fn(),
    dirxml: jest.fn(),
    group: jest.fn(),
    groupCollapsed: jest.fn(),
    groupEnd: jest.fn(),
    profile: jest.fn(),
    profileEnd: jest.fn(),
    table: jest.fn(),
    time: jest.fn(),
    timeEnd: jest.fn(),
    timeLog: jest.fn(),
    timeStamp: jest.fn(),
  }

  global.console = mockConsole as any
  return { originalConsole, mockConsole }
}

// Helper to restore console
export const restoreConsole = (originalConsole: Console) => {
  global.console = originalConsole
}

// Helper to test responsive behavior
export const testResponsive = async (testFn: () => Promise<void>) => {
  const viewports = [
    { width: 320, height: 568 },  // Mobile
    { width: 768, height: 1024 }, // Tablet
    { width: 1024, height: 768 }, // Desktop
    { width: 1920, height: 1080 }, // Large desktop
  ]

  const results = []

  for (const viewport of viewports) {
    simulateResize(viewport.width, viewport.height)
    await waitForAsync()
    
    const result = await testFn()
    results.push({
      viewport,
      result,
    })
  }

  return results
}

// Helper to create a test store with initial state
export const createTestStore = (initialState: any = {}) => {
  return {
    getState: () => initialState,
    setState: jest.fn(),
    subscribe: jest.fn(),
    destroy: jest.fn(),
  }
}

// Helper to test component lifecycle
export const testComponentLifecycle = async (Component: any, props: any = {}) => {
  const { unmount } = renderWithProviders(Component(props))
  await waitForAsync()
  
  // Test unmount
  unmount()
  await waitForAsync()
}

// Helper to test error boundaries
export const testErrorBoundary = (Component: any, error: Error) => {
  const ThrowError = () => {
    throw error
  }

  const { container } = renderWithProviders(
    React.createElement(Component, {}, React.createElement(ThrowError))
  )

  return container
}
