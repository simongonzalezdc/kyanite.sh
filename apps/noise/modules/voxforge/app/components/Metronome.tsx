'use client';

import { useEffect, useState } from 'react';
import * as Tone from 'tone';

interface MetronomeProps {
  bpm: number;
  isPlaying: boolean;
}

export default function Metronome({ bpm, isPlaying }: MetronomeProps) {
  const [currentBeat, setCurrentBeat] = useState(0);

  useEffect(() => {
    if (!isPlaying) {
      setCurrentBeat(0);
      return;
    }

    const interval = (60 / bpm) * 1000; // ms per beat
    let beat = 0;

    const timer = setInterval(() => {
      beat = (beat + 1) % 4;
      setCurrentBeat(beat);
    }, interval);

    return () => clearInterval(timer);
  }, [bpm, isPlaying]);

  return (
    <div className="flex gap-2 justify-center">
      {[0, 1, 2, 3].map((beat) => (
        <div
          key={beat}
          className={`w-3 h-3 rounded-full transition-colors ${
            beat === currentBeat && isPlaying
              ? 'bg-primary-500 scale-125'
              : 'bg-gray-600'
          }`}
        />
      ))}
    </div>
  );
}

