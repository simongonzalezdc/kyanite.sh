'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import { useAccessibility } from './AccessibilityProvider';

// Visual indicator for audio events
export function AudioVisualIndicator({ 
  isActive, 
  type = 'recording', 
  intensity = 'medium' 
}: { 
  isActive: boolean; 
  type?: 'recording' | 'playing' | 'processing' | 'error'; 
  intensity?: 'low' | 'medium' | 'high'; 
}) {
  const { settings } = useAccessibility();
  
  if (!settings.visualIndicators) {
    return null;
  }

  const getIndicatorStyles = () => {
    const baseStyles = 'fixed top-4 right-4 z-50 rounded-full transition-all duration-300';
    
    if (!isActive) {
      return `${baseStyles} opacity-0 scale-0`;
    }

    const typeStyles = {
      recording: 'bg-red-500 animate-pulse',
      playing: 'bg-green-500 animate-pulse',
      processing: 'bg-blue-500 animate-spin',
      error: 'bg-red-600 animate-bounce'
    };

    const sizeStyles = {
      low: 'w-8 h-8',
      medium: 'w-12 h-12',
      high: 'w-16 h-16'
    };

    return `${baseStyles} ${typeStyles[type]} ${sizeStyles[intensity]}`;
  };

  return (
    <div className={getIndicatorStyles()}>
      <span className="sr-only">
        {type === 'recording' && 'Recording in progress'}
        {type === 'playing' && 'Audio playing'}
        {type === 'processing' && 'Processing audio'}
        {type === 'error' && 'Audio error occurred'}
      </span>
    </div>
  );
}

// Vibration feedback for touch devices
export function VibrationFeedback({ 
  pattern = 'single', 
  trigger 
}: { 
  pattern?: 'single' | 'double' | 'triple' | 'long' | 'custom';
  trigger?: any;
}) {
  const { settings } = useAccessibility();
  const previousTrigger = useRef(trigger);

  useEffect(() => {
    if (trigger !== previousTrigger.current && settings.vibrationFeedback) {
      triggerVibration(pattern);
      previousTrigger.current = trigger;
    }
  }, [trigger, pattern, settings.vibrationFeedback]);

  const triggerVibration = (vibrationPattern: string) => {
    if (!('vibrate' in navigator)) {
      return;
    }

    const patterns = {
      single: [10],
      double: [10, 50, 10],
      triple: [10, 50, 10, 50, 10],
      long: [100],
      custom: [20, 30, 20, 30, 20]
    };

    navigator.vibrate(patterns[vibrationPattern as keyof typeof patterns] || patterns.single);
  };

  return null;
}

// Alternative input methods for recording
export function AlternativeRecordingInput({ 
  onRecordStart, 
  onRecordStop, 
  isRecording 
}: { 
  onRecordStart: () => void; 
  onRecordStop: () => void; 
  isRecording: boolean; 
}) {
  const [inputMethod, setInputMethod] = useState<'button' | 'spacebar' | 'voice'>('button');
  const { announce } = useAccessibility();

  const handleSpacebarRecord = useCallback((e: KeyboardEvent) => {
    if (e.code === 'Space' && inputMethod === 'spacebar') {
      e.preventDefault();
      if (isRecording) {
        onRecordStop();
        announce('Recording stopped');
      } else {
        onRecordStart();
        announce('Recording started');
      }
    }
  }, [inputMethod, isRecording, onRecordStart, onRecordStop, announce]);

  useEffect(() => {
    if (inputMethod === 'spacebar') {
      document.addEventListener('keydown', handleSpacebarRecord);
      return () => document.removeEventListener('keydown', handleSpacebarRecord);
    }
  }, [inputMethod, handleSpacebarRecord]);

  const handleVoiceCommand = useCallback(() => {
    if (!('webkitSpeechRecognition' in window) && !('SpeechRecognition' in window)) {
      announce('Voice commands not supported in this browser');
      return;
    }

    const SpeechRecognition = (window as any).SpeechRecognition || (window as any).webkitSpeechRecognition;
    const recognition = new SpeechRecognition();
    
    recognition.continuous = false;
    recognition.interimResults = false;
    recognition.lang = 'en-US';

    recognition.onresult = (event: any) => {
      const command = event.results[0][0].transcript.toLowerCase();
      
      if (command.includes('record') || command.includes('start')) {
        onRecordStart();
        announce('Recording started by voice command');
      } else if (command.includes('stop') || command.includes('end')) {
        onRecordStop();
        announce('Recording stopped by voice command');
      }
    };

    recognition.onerror = () => {
      announce('Voice command not recognized');
    };

    recognition.start();
    announce('Listening for voice commands...');
  }, [onRecordStart, onRecordStop, announce]);

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4 space-y-4">
      <h3 className="text-lg font-semibold">Alternative Recording Methods</h3>
      
      <div className="space-y-2">
        <label className="flex items-center space-x-2">
          <input
            type="radio"
            value="button"
            checked={inputMethod === 'button'}
            onChange={(e) => setInputMethod(e.target.value as any)}
            className="text-primary-500"
          />
          <span>Click Button</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="radio"
            value="spacebar"
            checked={inputMethod === 'spacebar'}
            onChange={(e) => setInputMethod(e.target.value as any)}
            className="text-primary-500"
          />
          <span>Press Spacebar</span>
        </label>
        
        <label className="flex items-center space-x-2">
          <input
            type="radio"
            value="voice"
            checked={inputMethod === 'voice'}
            onChange={(e) => setInputMethod(e.target.value as any)}
            className="text-primary-500"
          />
          <span>Voice Command</span>
        </label>
      </div>

      {inputMethod === 'spacebar' && (
        <div className="text-sm text-gray-400">
          Press and hold the spacebar to record
        </div>
      )}

      {inputMethod === 'voice' && (
        <button
          onClick={handleVoiceCommand}
          className="px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg transition-colors"
        >
          Enable Voice Commands
        </button>
      )}
    </div>
  );
}

// Visual waveform representation
export function VisualWaveform({ 
  audioBuffer, 
  isPlaying, 
  currentTime = 0 
}: { 
  audioBuffer: AudioBuffer | null; 
  isPlaying: boolean; 
  currentTime?: number; 
}) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const { settings } = useAccessibility();

  useEffect(() => {
    if (!audioBuffer || !canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const width = canvas.width;
    const height = canvas.height;
    const data = audioBuffer.getChannelData(0);
    const step = Math.ceil(data.length / width);
    const amp = height / 2;

    ctx.clearRect(0, 0, width, height);

    // Set colors based on accessibility settings
    const waveColor = settings.highContrastMode ? '#ffffff' : '#3B82F6';
    const progressColor = settings.highContrastMode ? '#ffff00' : '#10B981';

    // Draw waveform
    ctx.beginPath();
    ctx.moveTo(0, amp);
    ctx.strokeStyle = waveColor;
    ctx.lineWidth = settings.highContrastMode ? 3 : 2;

    for (let i = 0; i < width; i++) {
      let min = 1.0;
      let max = -1.0;
      
      for (let j = 0; j < step; j++) {
        const datum = data[(i * step) + j];
        if (datum < min) min = datum;
        if (datum > max) max = datum;
      }
      
      ctx.lineTo(i, (1 + min) * amp);
      ctx.lineTo(i, (1 + max) * amp);
    }

    ctx.stroke();

    // Draw progress indicator
    if (isPlaying && currentTime > 0) {
      const progress = (currentTime / audioBuffer.duration) * width;
      ctx.beginPath();
      ctx.moveTo(progress, 0);
      ctx.lineTo(progress, height);
      ctx.strokeStyle = progressColor;
      ctx.lineWidth = 2;
      ctx.stroke();
    }
  }, [audioBuffer, isPlaying, currentTime, settings.highContrastMode]);

  if (!audioBuffer) {
    return (
      <div className="w-full h-32 bg-gray-800 rounded-lg flex items-center justify-center">
        <span className="text-gray-400">No audio data</span>
      </div>
    );
  }

  return (
    <div className="relative">
      <canvas
        ref={canvasRef}
        width={800}
        height={200}
        className="w-full h-32 bg-gray-800 rounded-lg"
        role="img"
        aria-label="Audio waveform visualization"
      />
      {isPlaying && (
        <div className="absolute top-2 right-2">
          <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse" />
        </div>
      )}
    </div>
  );
}

// Caption support for audio content
export function AudioCaptions({ 
  captions, 
  currentTime, 
  isPlaying 
}: { 
  captions: Array<{ start: number; end: number; text: string }>; 
  currentTime: number; 
  isPlaying: boolean; 
}) {
  const [currentCaption, setCurrentCaption] = useState('');

  useEffect(() => {
    const activeCaption = captions.find(
      caption => currentTime >= caption.start && currentTime <= caption.end
    );
    
    setCurrentCaption(activeCaption ? activeCaption.text : '');
  }, [currentTime, captions]);

  if (!currentCaption) {
    return null;
  }

  return (
    <div className="bg-black bg-opacity-75 text-white p-4 rounded-lg text-center max-w-2xl mx-auto">
      <p className="text-lg">{currentCaption}</p>
    </div>
  );
}

// Audio event announcer for screen readers
export function AudioEventAnnouncer({ 
  events 
}: { 
  events: Array<{ type: string; message: string; timestamp: number }>; 
}) {
  const { announce, settings } = useAccessibility();
  const [lastEvent, setLastEvent] = useState<typeof events[0] | null>(null);

  useEffect(() => {
    if (events.length === 0) return;
    
    const latestEvent = events[events.length - 1];
    
    if (latestEvent !== lastEvent && settings.screenReaderEnabled) {
      announce(latestEvent.message);
      setLastEvent(latestEvent);
    }
  }, [events, lastEvent, announce, settings.screenReaderEnabled]);

  return null;
}

// Audio accessibility control panel
export function AudioAccessibilityControls() {
  const { settings, updateSettings } = useAccessibility();

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4 space-y-4">
      <h3 className="text-lg font-semibold">Audio Accessibility</h3>
      
      {/* Visual Indicators */}
      <div className="flex items-center justify-between">
        <label htmlFor="visual-indicators" className="text-sm font-medium">
          Visual Audio Indicators
        </label>
        <button
          id="visual-indicators"
          role="switch"
          aria-checked={settings.visualIndicators}
          onClick={() => updateSettings({ visualIndicators: !settings.visualIndicators })}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            settings.visualIndicators ? 'bg-primary-500' : 'bg-gray-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              settings.visualIndicators ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
      </div>

      {/* Vibration Feedback */}
      <div className="flex items-center justify-between">
        <label htmlFor="vibration-feedback" className="text-sm font-medium">
          Vibration Feedback
        </label>
        <button
          id="vibration-feedback"
          role="switch"
          aria-checked={settings.vibrationFeedback}
          onClick={() => updateSettings({ vibrationFeedback: !settings.vibrationFeedback })}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            settings.vibrationFeedback ? 'bg-primary-500' : 'bg-gray-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              settings.vibrationFeedback ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
      </div>

      {/* Audio Descriptions */}
      <div className="flex items-center justify-between">
        <label htmlFor="audio-descriptions" className="text-sm font-medium">
          Audio Descriptions
        </label>
        <button
          id="audio-descriptions"
          role="switch"
          aria-checked={settings.audioDescriptions}
          onClick={() => updateSettings({ audioDescriptions: !settings.audioDescriptions })}
          className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
            settings.audioDescriptions ? 'bg-primary-500' : 'bg-gray-600'
          }`}
        >
          <span
            className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
              settings.audioDescriptions ? 'translate-x-6' : 'translate-x-1'
            }`}
          />
        </button>
      </div>
    </div>
  );
}

// Test vibration functionality
export function VibrationTest() {
  const { settings } = useAccessibility();

  const testVibration = (pattern: 'single' | 'double' | 'triple' | 'long') => {
    if (!('vibrate' in navigator) || !settings.vibrationFeedback) {
      return;
    }

    const patterns = {
      single: [10],
      double: [10, 50, 10],
      triple: [10, 50, 10, 50, 10],
      long: [100]
    };

    navigator.vibrate(patterns[pattern]);
  };

  if (!('vibrate' in navigator)) {
    return (
      <div className="bg-gray-900 border border-gray-700 rounded-lg p-4">
        <h3 className="text-lg font-semibold mb-2">Vibration Test</h3>
        <p className="text-gray-400">Vibration API not supported on this device</p>
      </div>
    );
  }

  return (
    <div className="bg-gray-900 border border-gray-700 rounded-lg p-4">
      <h3 className="text-lg font-semibold mb-4">Vibration Test</h3>
      
      <div className="grid grid-cols-2 gap-2">
        <button
          onClick={() => testVibration('single')}
          className="px-3 py-2 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
        >
          Single
        </button>
        <button
          onClick={() => testVibration('double')}
          className="px-3 py-2 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
        >
          Double
        </button>
        <button
          onClick={() => testVibration('triple')}
          className="px-3 py-2 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
        >
          Triple
        </button>
        <button
          onClick={() => testVibration('long')}
          className="px-3 py-2 bg-gray-800 hover:bg-gray-700 rounded-md transition-colors"
        >
          Long
        </button>
      </div>
    </div>
  );
}