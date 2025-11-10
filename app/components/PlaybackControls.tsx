'use client';

import { Play, Pause, Square } from 'lucide-react';

interface PlaybackControlsProps {
  isPlaying: boolean;
  onPlay: () => void;
  onPause: () => void;
  onStop: () => void;
}

export default function PlaybackControls({ isPlaying, onPlay, onPause, onStop }: PlaybackControlsProps) {
  return (
    <div className="flex gap-4 justify-center">
      {!isPlaying ? (
        <button
          onClick={onPlay}
          className="flex items-center gap-2 px-6 py-3 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors"
        >
          <Play size={20} />
          Play
        </button>
      ) : (
        <button
          onClick={onPause}
          className="flex items-center gap-2 px-6 py-3 bg-yellow-500 hover:bg-yellow-600 rounded-lg font-medium transition-colors"
        >
          <Pause size={20} />
          Pause
        </button>
      )}
      
      <button
        onClick={onStop}
        className="flex items-center gap-2 px-6 py-3 bg-red-500 hover:bg-red-600 rounded-lg font-medium transition-colors"
      >
        <Square size={20} />
        Stop
      </button>
    </div>
  );
}

