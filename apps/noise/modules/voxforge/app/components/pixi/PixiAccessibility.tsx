'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import * as PIXI from 'pixi.js';
import { useAccessibility } from '../AccessibilityProvider';

// Text-based alternatives for canvas visualizations
export function PixiTextAlternative({ 
  visualizationData, 
  type 
}: { 
  visualizationData: any; 
  type: 'piano-roll' | 'visualizer' | 'rhythm-game'; 
}) {
  const { settings } = useAccessibility();

  if (!settings.screenReaderEnabled) {
    return null;
  }

  const generateTextDescription = () => {
    switch (type) {
      case 'piano-roll':
        return `Piano roll showing ${visualizationData?.notes?.length || 0} notes across ${visualizationData?.duration || 0} seconds. Notes range from ${visualizationData?.lowestNote || 'C3'} to ${visualizationData?.highestNote || 'C5'}.`;
      
      case 'visualizer':
        return `Audio visualizer showing ${visualizationData?.frequencyBands || 8} frequency bands. Current activity level: ${visualizationData?.activityLevel || 'medium'}.`;
      
      case 'rhythm-game':
        return `Rhythm game with ${visualizationData?.notes?.length || 0} notes to hit. Current score: ${visualizationData?.score || 0}. Combo: ${visualizationData?.combo || 0}.`;
      
      default:
        return 'Interactive audio visualization';
    }
  };

  return (
    <div className="sr-only" role="status" aria-live="polite">
      {generateTextDescription()}
    </div>
  );
}

// Keyboard navigation for interactive elements
export function PixiKeyboardNavigation({ 
  app, 
  onNavigate 
}: { 
  app: PIXI.Application | null; 
  onNavigate: (direction: 'up' | 'down' | 'left' | 'right' | 'select') => void; 
}) {
  const { settings } = useAccessibility();
  const [focusedElement, setFocusedElement] = useState<number>(0);

  useEffect(() => {
    if (!settings.keyboardNavigation || !app) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'ArrowUp':
          e.preventDefault();
          onNavigate('up');
          break;
        case 'ArrowDown':
          e.preventDefault();
          onNavigate('down');
          break;
        case 'ArrowLeft':
          e.preventDefault();
          onNavigate('left');
          break;
        case 'ArrowRight':
          e.preventDefault();
          onNavigate('right');
          break;
        case 'Enter':
        case ' ':
          e.preventDefault();
          onNavigate('select');
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [app, settings.keyboardNavigation, onNavigate]);

  return null;
}

// Screen reader announcements for visual changes
export function PixiScreenReaderAnnouncer({ 
  app, 
  events 
}: { 
  app: PIXI.Application | null; 
  events: Array<{ type: string; data: any; timestamp: number }>; 
}) {
  const { announce, settings } = useAccessibility();
  const [lastEvent, setLastEvent] = useState<typeof events[0] | null>(null);

  useEffect(() => {
    if (!settings.screenReaderEnabled || events.length === 0) return;
    
    const latestEvent = events[events.length - 1];
    
    if (latestEvent !== lastEvent) {
      const message = generateAnnouncementMessage(latestEvent);
      announce(message);
      setLastEvent(latestEvent);
    }
  }, [events, lastEvent, announce, settings.screenReaderEnabled]);

  const generateAnnouncementMessage = (event: any) => {
    switch (event.type) {
      case 'note-added':
        return `Note added at position ${event.data.x}, pitch ${event.data.pitch}`;
      case 'note-removed':
        return `Note removed`;
      case 'note-selected':
        return `Note selected at position ${event.data.x}, pitch ${event.data.pitch}`;
      case 'playback-started':
        return 'Playback started';
      case 'playback-stopped':
        return 'Playback stopped';
      case 'tempo-changed':
        return `Tempo changed to ${event.data.bpm} beats per minute`;
      case 'instrument-changed':
        return `Instrument changed to ${event.data.instrument}`;
      default:
        return 'Visualization updated';
    }
  };

  return null;
}

// High contrast rendering modes
export function PixiHighContrast({ 
  app, 
  enabled 
}: { 
  app: PIXI.Application | null; 
  enabled: boolean; 
}) {
  useEffect(() => {
    if (!app) return;

    const updateColors = () => {
      if (enabled) {
        // Apply high contrast colors
        (app.renderer as any).background.color = 0x000000;
        
        // Update all graphics to high contrast colors
        app.stage.children.forEach(child => {
          if (child instanceof PIXI.Graphics) {
            // Convert to high contrast colors
            child.tint = 0xFFFFFF;
          }
        });
      } else {
        // Restore original colors
        (app.renderer as any).background.color = 0x0a0a0a;
        
        // Restore original tints
        app.stage.children.forEach(child => {
          if (child instanceof PIXI.Graphics) {
            child.tint = 0xFFFFFF;
          }
        });
      }
    };

    updateColors();
  }, [app, enabled]);

  return null;
}

// Focus management in canvas
export function PixiFocusManager({ 
  app, 
  focusableElements 
}: { 
  app: PIXI.Application | null; 
  focusableElements: Array<{ id: string; x: number; y: number; width: number; height: number; label: string }>; 
}) {
  const [focusedIndex, setFocusedIndex] = useState(-1);
  const { announce, settings } = useAccessibility();

  const drawFocusIndicator = useCallback(() => {
    if (!app || focusedIndex < 0 || focusedIndex >= focusableElements.length) return;

    const element = focusableElements[focusedIndex];
    
    // Remove existing focus indicator
    const existingIndicator = app.stage.getChildByName('focus-indicator');
    if (existingIndicator) {
      app.stage.removeChild(existingIndicator);
    }

    // Create new focus indicator
    const focusIndicator = new PIXI.Graphics();
    focusIndicator.name = 'focus-indicator';
    focusIndicator.lineStyle(
      settings.highContrastMode ? 4 : 2,
      settings.highContrastMode ? 0xFFFFFF : 0x3B82F6,
      1
    );
    focusIndicator.drawRect(element.x - 2, element.y - 2, element.width + 4, element.height + 4);
    
    app.stage.addChild(focusIndicator);
    
    // Announce focused element
    if (settings.screenReaderEnabled) {
      announce(`Focused on ${element.label}`);
    }
  }, [app, focusedIndex, focusableElements, settings]);

  useEffect(() => {
    drawFocusIndicator();
  }, [drawFocusIndicator]);

  const setFocus = useCallback((index: number) => {
    setFocusedIndex(index);
  }, []);

  const nextFocus = useCallback(() => {
    setFocusedIndex(prev => (prev + 1) % focusableElements.length);
  }, [focusableElements.length]);

  const previousFocus = useCallback(() => {
    setFocusedIndex(prev => (prev - 1 + focusableElements.length) % focusableElements.length);
  }, [focusableElements.length]);

  return { focusedIndex, setFocus, nextFocus, previousFocus };
}

// Accessible piano roll editor
export function AccessiblePianoRoll({ 
  notes, 
  onNoteAdd, 
  onNoteRemove, 
  onNoteSelect 
}: { 
  notes: Array<{ id: string; x: number; y: number; pitch: string; time: number }>; 
  onNoteAdd: (x: number, y: number) => void; 
  onNoteRemove: (id: string) => void; 
  onNoteSelect: (id: string) => void; 
}) {
  const { settings } = useAccessibility();
  const [selectedNote, setSelectedNote] = useState<string | null>(null);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!settings.keyboardNavigation) return;

    switch (e.key) {
      case 'ArrowRight':
        // Move to next time position
        break;
      case 'ArrowLeft':
        // Move to previous time position
        break;
      case 'ArrowUp':
        // Move to higher pitch
        break;
      case 'ArrowDown':
        // Move to lower pitch
        break;
      case 'Enter':
      case ' ':
        // Add or select note
        if (selectedNote) {
          onNoteSelect(selectedNote);
        }
        break;
      case 'Delete':
      case 'Backspace':
        // Remove selected note
        if (selectedNote) {
          onNoteRemove(selectedNote);
          setSelectedNote(null);
        }
        break;
    }
  };

  return (
    <div 
      className="relative w-full h-64 bg-gray-800 rounded-lg"
      onKeyDown={handleKeyDown}
      tabIndex={0}
      role="grid"
      aria-label="Piano roll editor"
    >
      {/* Text-based representation for screen readers */}
      <div className="sr-only" role="status" aria-live="polite">
        Piano roll contains {notes.length} notes. 
        {notes.map((note, index) => (
          <span key={note.id}>
            Note {index + 1}: {note.pitch} at {note.time.toFixed(2)} seconds. 
          </span>
        ))}
      </div>

      {/* Visual representation */}
      <div className="absolute inset-0">
        {notes.map(note => (
          <div
            key={note.id}
            className={`absolute w-4 h-4 rounded cursor-pointer transition-colors ${
              selectedNote === note.id 
                ? 'bg-primary-500' 
                : 'bg-gray-600 hover:bg-gray-500'
            }`}
            style={{ left: `${note.x}px`, top: `${note.y}px` }}
            onClick={() => {
              setSelectedNote(note.id);
              onNoteSelect(note.id);
            }}
            role="gridcell"
            aria-label={`Note: ${note.pitch} at ${note.time.toFixed(2)} seconds`}
            tabIndex={-1}
          />
        ))}
      </div>

      {/* Instructions for keyboard users */}
      {settings.keyboardNavigation && (
        <div className="absolute bottom-2 left-2 text-xs text-gray-400">
          Use arrow keys to navigate, Enter to select, Delete to remove
        </div>
      )}
    </div>
  );
}

// Accessible visualizer canvas
export function AccessibleVisualizer({ 
  frequencyData, 
  isRecording 
}: { 
  frequencyData: number[]; 
  isRecording: boolean; 
}) {
  const { settings, announce } = useAccessibility();

  useEffect(() => {
    if (settings.screenReaderEnabled && isRecording) {
      // Announce significant changes in audio levels
      const averageLevel = frequencyData.reduce((sum, val) => sum + val, 0) / frequencyData.length;
      
      if (averageLevel > 200) {
        announce('High audio level detected');
      } else if (averageLevel < 50) {
        announce('Low audio level detected');
      }
    }
  }, [frequencyData, isRecording, settings.screenReaderEnabled, announce]);

  return (
    <div className="relative w-full h-32 bg-gray-800 rounded-lg">
      {/* Text-based representation */}
      <div className="sr-only" role="status" aria-live="polite">
        {isRecording ? 'Recording in progress' : 'Not recording'}
        {frequencyData.length > 0 && (
          <span>
            Audio levels: {frequencyData.map((level, index) => 
              `Band ${index + 1}: ${Math.round(level)}`
            ).join(', ')}
          </span>
        )}
      </div>

      {/* Visual representation */}
      <div className="absolute inset-0 flex items-end justify-center space-x-1 p-2">
        {frequencyData.map((level, index) => (
          <div
            key={index}
            className="flex-1 bg-primary-500 transition-all duration-100"
            style={{ 
              height: `${(level / 255) * 100}%`,
              backgroundColor: settings.highContrastMode ? '#ffffff' : '#3B82F6'
            }}
            role="img"
            aria-label={`Frequency band ${index + 1} level: ${Math.round(level)}`}
          />
        ))}
      </div>

      {/* Recording indicator */}
      {isRecording && (
        <div className="absolute top-2 right-2">
          <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse" />
        </div>
      )}
    </div>
  );
}

// Accessible rhythm game mode
export function AccessibleRhythmGame({ 
  notes, 
  onHit, 
  onMiss, 
  score 
}: { 
  notes: Array<{ id: string; time: number; position: number; hit?: boolean }>; 
  onHit: (id: string) => void; 
  onMiss: (id: string) => void; 
  score: number; 
}) {
  const { settings, announce } = useAccessibility();
  const [nextNote, setNextNote] = useState<typeof notes[0] | null>(null);

  useEffect(() => {
    // Find the next unhit note
    const upcoming = notes.find(note => !note.hit && note.time > Date.now());
    if (upcoming !== nextNote) {
      setNextNote(upcoming || null);
      if (upcoming && settings.screenReaderEnabled) {
        announce(`Note approaching at position ${upcoming.position}`);
      }
    }
  }, [notes, nextNote, settings.screenReaderEnabled, announce]);

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (!settings.keyboardNavigation) return;

    if (e.key === ' ' || e.key === 'Enter') {
      // Hit the note
      if (nextNote) {
        onHit(nextNote.id);
        announce('Note hit!');
      }
    }
  };

  return (
    <div 
      className="relative w-full h-48 bg-gray-800 rounded-lg"
      onKeyDown={handleKeyPress}
      tabIndex={0}
      role="application"
      aria-label="Rhythm game"
    >
      {/* Game status for screen readers */}
      <div className="sr-only" role="status" aria-live="polite">
        Score: {score}. 
        Notes remaining: {notes.filter(note => !note.hit).length}.
        {nextNote && `Next note at position ${nextNote.position}.`}
      </div>

      {/* Visual game area */}
      <div className="absolute inset-0 flex items-center justify-center">
        <div className="w-full h-1 bg-gray-600 absolute" />
        
        {/* Hit zone */}
        <div className="w-16 h-16 border-2 border-primary-500 rounded-full flex items-center justify-center">
          <span className="text-primary-500 font-bold">HIT</span>
        </div>

        {/* Notes */}
        {notes.map(note => (
          <div
            key={note.id}
            className={`absolute w-8 h-8 rounded-full transition-all ${
              note.hit 
                ? 'bg-green-500' 
                : 'bg-red-500'
            }`}
            style={{ 
              left: `${note.position}%`,
              top: '50%',
              transform: 'translate(-50%, -50%)'
            }}
            role="img"
            aria-label={`Note at position ${note.position}${note.hit ? ' - hit' : ' - pending'}`}
          />
        ))}
      </div>

      {/* Score display */}
      <div className="absolute top-2 right-2 text-white font-bold">
        Score: {score}
      </div>

      {/* Instructions */}
      {settings.keyboardNavigation && (
        <div className="absolute bottom-2 left-2 text-xs text-gray-400">
          Press Space or Enter to hit notes
        </div>
      )}
    </div>
  );
}

// Main PIXI accessibility wrapper
export function PixiAccessibilityWrapper({ 
  children, 
  type, 
  data 
}: { 
  children: React.ReactNode; 
  type: 'piano-roll' | 'visualizer' | 'rhythm-game'; 
  data: any; 
}) {
  return (
    <div className="relative">
      {/* Text alternative for screen readers */}
      <PixiTextAlternative visualizationData={data} type={type} />
      
      {/* Main canvas content */}
      {children}
      
      {/* Accessibility controls overlay */}
      <div className="absolute top-2 left-2 bg-black bg-opacity-50 text-white p-2 rounded text-xs">
        Accessibility mode active
      </div>
    </div>
  );
}