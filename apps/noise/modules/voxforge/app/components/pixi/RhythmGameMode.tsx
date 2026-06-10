'use client';

import { useEffect, useRef, useState } from 'react';
import { RhythmGameRecorder } from '@/lib/pixi/rhythm-game';
import { Play, Square, Trophy } from 'lucide-react';

interface RhythmGameModeProps {
  bpm: number;
  onComplete?: (notes: Array<{ time: number; midi: number }>) => void;
}

export default function RhythmGameMode({ bpm, onComplete }: RhythmGameModeProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const gameRef = useRef<RhythmGameRecorder | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [score, setScore] = useState(0);
  const [combo, setCombo] = useState(0);
  const [highScore, setHighScore] = useState(0);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create game
    const game = new RhythmGameRecorder(
      { fallSpeed: 200 },
      (newScore, newCombo) => {
        setScore(newScore);
        setCombo(newCombo);
      }
    );

    // Wait for initialization, then add canvas to DOM
    game.waitForInit().then(() => {
      const canvas = game.getCanvas();
      canvas.style.display = 'block';
      canvas.style.width = '100%';
      canvas.style.height = '100%';
      containerRef.current?.appendChild(canvas);
      gameRef.current = game;

      // Load demo sequence
      const demoNotes = generateDemoSequence(bpm);
      game.loadSequence(demoNotes, bpm);
    }).catch(console.error);

    return () => {
      game.destroy();
    };
  }, [bpm]);

  const handleStart = () => {
    if (!gameRef.current) return;
    
    gameRef.current.start();
    setIsPlaying(true);
    setScore(0);
    setCombo(0);
  };

  const handleStop = () => {
    if (!gameRef.current) return;
    
    const recordedNotes = gameRef.current.stop();
    setIsPlaying(false);
    
    // Update high score
    if (score > highScore) {
      setHighScore(score);
    }
    
    onComplete?.(recordedNotes);
  };

  return (
    <div className="space-y-4">
      {/* Score Display */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-gray-900 p-4 rounded-lg border border-gray-800">
          <div className="text-sm text-gray-400">Score</div>
          <div className="text-3xl font-bold text-primary-500">{score}</div>
        </div>
        
        <div className="bg-gray-900 p-4 rounded-lg border border-gray-800">
          <div className="text-sm text-gray-400">Combo</div>
          <div className="text-3xl font-bold text-secondary-500">
            {combo > 0 ? `×${combo}` : '0'}
          </div>
        </div>
        
        <div className="bg-gray-900 p-4 rounded-lg border border-gray-800 flex items-center gap-2">
          <Trophy className="text-yellow-500" size={20} />
          <div>
            <div className="text-sm text-gray-400">High Score</div>
            <div className="text-2xl font-bold text-yellow-500">{highScore}</div>
          </div>
        </div>
      </div>

      {/* Instructions */}
      <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-800">
        <p className="text-sm text-gray-400">
          <strong>How to Play:</strong> Press keys 1-5 or spacebar when notes reach the target line. 
          Perfect timing = higher score & combo multiplier!
        </p>
      </div>

      {/* Game Canvas */}
      <div 
        ref={containerRef}
        className="rounded-lg overflow-hidden border border-gray-800 relative"
        style={{ minHeight: '600px', maxHeight: '600px' }}
      />

      {/* Controls */}
      <div className="flex justify-center">
        {!isPlaying ? (
          <button
            onClick={handleStart}
            className="flex items-center gap-2 px-8 py-4 bg-primary-500 hover:bg-primary-600 rounded-lg text-lg font-semibold"
          >
            <Play size={24} />
            Start Game
          </button>
        ) : (
          <button
            onClick={handleStop}
            className="flex items-center gap-2 px-8 py-4 bg-red-500 hover:bg-red-600 rounded-lg text-lg font-semibold"
          >
            <Square size={24} />
            Stop & Save
          </button>
        )}
      </div>
    </div>
  );
}

// Generate demo sequence for testing
function generateDemoSequence(bpm: number): Array<{ time: number; midi: number }> {
  const beatDuration = 60 / bpm;
  const notes: Array<{ time: number; midi: number }> = [];
  
  // Generate a simple melody
  const melody = [60, 62, 64, 65, 67, 65, 64, 62];
  
  melody.forEach((midi, i) => {
    notes.push({
      time: i * beatDuration + 2, // Start after 2 seconds
      midi
    });
  });
  
  return notes;
}

