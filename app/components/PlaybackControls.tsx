'use client';

import { Play, Pause, Square } from 'lucide-react';
import { isMobile, triggerHaptic } from '@/lib/utils';

interface PlaybackControlsProps {
  isPlaying: boolean;
  onPlay: () => void;
  onPause: () => void;
  onStop: () => void;
}

export default function PlaybackControls({ isPlaying, onPlay, onPause, onStop }: PlaybackControlsProps) {
  return (
    <div className="flex gap-3 sm:gap-4 justify-center">
      {!isPlaying ? (
        <button
          onClick={() => {
            onPlay();
            triggerHaptic('light');
          }}
          className={`
            flex items-center justify-center gap-2 px-6 py-4 sm:py-3
            bg-primary-500 hover:bg-primary-600
            rounded-lg font-medium transition-all duration-200
            min-h-[44px] min-w-[44px] sm:min-w-0
            ${isMobile() ? 'text-lg px-8' : ''}
            active:scale-95 touch-manipulation
          `}
          style={{ minHeight: isMobile() ? '60px' : '48px' }}
        >
          <Play size={isMobile() ? 24 : 20} />
          <span>Play</span>
        </button>
      ) : (
        <button
          onClick={() => {
            onPause();
            triggerHaptic('medium');
          }}
          className={`
            flex items-center justify-center gap-2 px-6 py-4 sm:py-3
            bg-yellow-500 hover:bg-yellow-600
            rounded-lg font-medium transition-all duration-200
            min-h-[44px] min-w-[44px] sm:min-w-0
            ${isMobile() ? 'text-lg px-8' : ''}
            active:scale-95 touch-manipulation
          `}
          style={{ minHeight: isMobile() ? '60px' : '48px' }}
        >
          <Pause size={isMobile() ? 24 : 20} />
          <span>Pause</span>
        </button>
      )}
      
      <button
        onClick={() => {
          onStop();
          triggerHaptic('heavy');
        }}
        className={`
          flex items-center justify-center gap-2 px-6 py-4 sm:py-3
          bg-red-500 hover:bg-red-600
          rounded-lg font-medium transition-all duration-200
          min-h-[44px] min-w-[44px] sm:min-w-0
          ${isMobile() ? 'text-lg' : ''}
          active:scale-95 touch-manipulation
        `}
        style={{ minHeight: isMobile() ? '48px' : '44px' }}
      >
        <Square size={isMobile() ? 24 : 20} />
        <span>Stop</span>
      </button>
    </div>
  );
}

