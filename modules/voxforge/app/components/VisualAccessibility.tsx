'use client';

import React, { useEffect, useState } from 'react';
import { useAccessibility } from './AccessibilityProvider';

// Visual accessibility control panel
export function VisualAccessibilityControls() {
  const { settings, updateSettings } = useAccessibility();

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4 space-y-4">
      <h3 className="text-lg font-semibold">Visual Accessibility</h3>
      
      {/* High Contrast Mode */}
      <div className="flex items-center justify-between">
        <label htmlFor="high-contrast" className="text-sm font-medium">
          High Contrast Mode
        </label>
        <button
          id="high-contrast"
          role="switch"
          aria-checked={settings.highContrastMode}
          onClick={() => updateSettings({ highContrastMode: !settings.highContrastMode })}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            settings.highContrastMode ? 'bg-primary-500' : 'bg-gray-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              settings.highContrastMode ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
      </div>

      {/* Font Size Control */}
      <div>
        <label htmlFor="font-size" className="text-sm font-medium block mb-2">
          Font Size
        </label>
        <select
          id="font-size"
          value={settings.fontSize}
          onChange={(e) => updateSettings({ fontSize: e.target.value as any })}
          className="w-full px-3 py-2 bg-gray-800 border border-gray-600 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          <option value="small">Small</option>
          <option value="medium">Medium</option>
          <option value="large">Large</option>
          <option value="extra-large">Extra Large</option>
        </select>
      </div>

      {/* Reduced Motion */}
      <div className="flex items-center justify-between">
        <label htmlFor="reduced-motion" className="text-sm font-medium">
          Reduced Motion
        </label>
        <button
          id="reduced-motion"
          role="switch"
          aria-checked={settings.reducedMotion}
          onClick={() => updateSettings({ reducedMotion: !settings.reducedMotion })}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            settings.reducedMotion ? 'bg-primary-500' : 'bg-gray-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              settings.reducedMotion ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
      </div>

      {/* Focus Indicators */}
      <div className="flex items-center justify-between">
        <label htmlFor="focus-visible" className="text-sm font-medium">
          Focus Indicators
        </label>
        <button
          id="focus-visible"
          role="switch"
          aria-checked={settings.focusVisible}
          onClick={() => updateSettings({ focusVisible: !settings.focusVisible })}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            settings.focusVisible ? 'bg-primary-500' : 'bg-gray-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              settings.focusVisible ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
      </div>
    </div>
  );
}

// High contrast theme provider
export function HighContrastProvider({ children }: { children: React.ReactNode }) {
  const { settings } = useAccessibility();

  useEffect(() => {
    const root = document.documentElement;
    
    if (settings.highContrastMode) {
      root.setAttribute('data-theme', 'high-contrast');
    } else {
      root.removeAttribute('data-theme');
    }
  }, [settings.highContrastMode]);

  return <>{children}</>;
}

// Font size provider
export function FontSizeProvider({ children }: { children: React.ReactNode }) {
  const { settings } = useAccessibility();

  useEffect(() => {
    const root = document.documentElement;
    
    // Apply font size classes
    root.classList.remove('text-size-small', 'text-size-medium', 'text-size-large', 'text-size-extra-large');
    root.classList.add(`text-size-${settings.fontSize}`);
  }, [settings.fontSize]);

  return <>{children}</>;
}

// Color blind friendly palette component
export function ColorBlindPalette() {
  const [palette, setPalette] = useState<'normal' | 'protanopia' | 'deuteranopia' | 'tritanopia'>('normal');

  useEffect(() => {
    const root = document.documentElement;
    root.setAttribute('data-color-blind', palette);
  }, [palette]);

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-4">Color Vision Support</h3>
      <div className="space-y-2">
        {[
          { value: 'normal', label: 'Normal Vision' },
          { value: 'protanopia', label: 'Protanopia (Red-Blind)' },
          { value: 'deuteranopia', label: 'Deuteranopia (Green-Blind)' },
          { value: 'tritanopia', label: 'Tritanopia (Blue-Blind)' }
        ].map(option => (
          <button
            key={option.value}
            onClick={() => setPalette(option.value as any)}
            className={`w-full px-3 py-2 text-left rounded-md transition-colors ${
              palette === option.value
                ? 'bg-primary-500 text-white'
                : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
            }`}
          >
            {option.label}
          </button>
        ))}
      </div>
    </div>
  );
}

// Text resizing without breaking layout
export function TextResizer() {
  const { settings, updateSettings } = useAccessibility();
  const [customScale, setCustomScale] = useState(100);

  const handleScaleChange = (scale: number) => {
    setCustomScale(scale);
    const root = document.documentElement;
    root.style.fontSize = `${scale}%`;
  };

  const presetSizes = [
    { label: '100%', value: 100 },
    { label: '125%', value: 125 },
    { label: '150%', value: 150 },
    { label: '175%', value: 175 },
    { label: '200%', value: 200 }
  ];

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-4">Text Size</h3>
      
      <div className="space-y-4">
        {/* Preset sizes */}
        <div className="grid grid-cols-3 gap-2">
          {presetSizes.map(size => (
            <button
              key={size.value}
              onClick={() => handleScaleChange(size.value)}
              className={`px-3 py-2 rounded-md transition-colors ${
                customScale === size.value
                  ? 'bg-primary-500 text-white'
                  : 'bg-gray-800 text-gray-300 hover:bg-gray-700'
              }`}
            >
              {size.label}
            </button>
          ))}
        </div>

        {/* Custom slider */}
        <div>
          <label htmlFor="text-scale" className="text-sm font-medium block mb-2">
            Custom Scale: {customScale}%
          </label>
          <input
            id="text-scale"
            type="range"
            min="100"
            max="200"
            step="5"
            value={customScale}
            onChange={(e) => handleScaleChange(Number(e.target.value))}
            className="w-full"
          />
        </div>

        {/* Reset button */}
        <button
          onClick={() => {
            handleScaleChange(100);
            updateSettings({ fontSize: 'medium' });
          }}
          className="w-full px-3 py-2 bg-gray-800 text-gray-300 rounded-md hover:bg-gray-700 transition-colors"
        >
          Reset to Default
        </button>
      </div>
    </div>
  );
}

// Contrast checker component
export function ContrastChecker() {
  const [foreground, setForeground] = useState('#ffffff');
  const [background, setBackground] = useState('#000000');
  const [ratio, setRatio] = useState(21);

  const calculateContrast = (fg: string, bg: string) => {
    const getLuminance = (color: string) => {
      const rgb = parseInt(color.slice(1), 16);
      const r = (rgb >> 16) & 0xff;
      const g = (rgb >> 8) & 0xff;
      const b = rgb & 0xff;
      
      const [rs, gs, bs] = [r, g, b].map(c => {
        c = c / 255;
        return c <= 0.03928 ? c / 12.92 : Math.pow((c + 0.055) / 1.055, 2.4);
      });
      
      return 0.2126 * rs + 0.7152 * gs + 0.0722 * bs;
    };

    const l1 = getLuminance(fg);
    const l2 = getLuminance(bg);
    const contrast = (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05);
    
    return Math.round(contrast * 100) / 100;
  };

  useEffect(() => {
    setRatio(calculateContrast(foreground, background));
  }, [foreground, background]);

  const getCompliance = (ratio: number) => {
    if (ratio >= 7) return { level: 'AAA', color: 'text-green-500' };
    if (ratio >= 4.5) return { level: 'AA', color: 'text-blue-500' };
    if (ratio >= 3) return { level: 'AA Large', color: 'text-yellow-500' };
    return { level: 'Fail', color: 'text-red-500' };
  };

  const compliance = getCompliance(ratio);

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-4">Contrast Checker</h3>
      
      <div className="space-y-4">
        {/* Color inputs */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label htmlFor="foreground" className="text-sm font-medium block mb-2">
              Foreground
            </label>
            <div className="flex gap-2">
              <input
                id="foreground"
                type="color"
                value={foreground}
                onChange={(e) => setForeground(e.target.value)}
                className="h-10 w-20"
              />
              <input
                type="text"
                value={foreground}
                onChange={(e) => setForeground(e.target.value)}
                className="flex-1 px-3 py-2 bg-gray-800 border border-gray-600 rounded-md text-sm"
              />
            </div>
          </div>
          
          <div>
            <label htmlFor="background" className="text-sm font-medium block mb-2">
              Background
            </label>
            <div className="flex gap-2">
              <input
                id="background"
                type="color"
                value={background}
                onChange={(e) => setBackground(e.target.value)}
                className="h-10 w-20"
              />
              <input
                type="text"
                value={background}
                onChange={(e) => setBackground(e.target.value)}
                className="flex-1 px-3 py-2 bg-gray-800 border border-gray-600 rounded-md text-sm"
              />
            </div>
          </div>
        </div>

        {/* Preview */}
        <div
          className="p-4 rounded-md text-center font-medium"
          style={{ backgroundColor: background, color: foreground }}
        >
          Sample Text
        </div>

        {/* Results */}
        <div className="space-y-2">
          <div className="flex justify-between">
            <span className="text-sm">Contrast Ratio:</span>
            <span className="text-sm font-mono">{ratio}:1</span>
          </div>
          <div className="flex justify-between">
            <span className="text-sm">WCAG Compliance:</span>
            <span className={`text-sm font-medium ${compliance.color}`}>
              {compliance.level}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

// Visual accessibility toolbar
export function VisualAccessibilityToolbar() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="fixed bottom-4 right-4 z-50">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="bg-primary-500 hover:bg-primary-600 text-white p-3 rounded-full shadow-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
        aria-label="Toggle accessibility options"
        aria-expanded={isOpen}
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute bottom-16 right-0 w-80 bg-gray-900 border border-gray-700 rounded-lg shadow-xl p-4 space-y-4 max-h-96 overflow-y-auto">
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold">Accessibility Options</h3>
            <button
              onClick={() => setIsOpen(false)}
              className="text-gray-400 hover:text-white"
              aria-label="Close accessibility options"
            >
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <VisualAccessibilityControls />
          <ColorBlindPalette />
          <TextResizer />
          <ContrastChecker />
        </div>
      )}
    </div>
  );
}