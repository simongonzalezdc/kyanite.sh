'use client';

import { useState } from 'react';
import Recorder from './components/Recorder';
import AnalysisDisplay from './components/AnalysisDisplay';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { PitchPoint } from '@/lib/types';

export default function Home() {
  const [audioBuffer, setAudioBuffer] = useState<AudioBuffer | null>(null);
  const [pitches, setPitches] = useState<PitchPoint[]>([]);
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  const handleRecordingComplete = async (buffer: AudioBuffer) => {
    setAudioBuffer(buffer);
    setIsAnalyzing(true);

    // Run pitch detection
    try {
      const detector = new PitchDetector(buffer.sampleRate);
      const detectedPitches = detector.analyze(buffer);
      setPitches(detectedPitches);
    } catch (error) {
      console.error('Pitch detection failed:', error);
      alert('Failed to analyze pitch. Please try again.');
    } finally {
      setIsAnalyzing(false);
    }
  };

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <header className="text-center space-y-2">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            VoxForge
          </h1>
          <p className="text-gray-400">
            Transform your voice into music
          </p>
        </header>

        <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
          <h2 className="text-xl font-semibold mb-4">Step 1: Record Your Voice</h2>
          <Recorder onRecordingComplete={handleRecordingComplete} />
        </div>

        {isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
            <div className="flex items-center justify-center gap-3">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500"></div>
              <p className="text-gray-400">Analyzing audio...</p>
            </div>
          </div>
        )}

        {pitches.length > 0 && !isAnalyzing && (
          <>
            <AnalysisDisplay pitches={pitches} />
            
            <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
              <p className="text-gray-400 text-center">
                Next: BPM and Key detection will be added in Day 3
              </p>
            </div>
          </>
        )}
      </div>
    </main>
  );
}

