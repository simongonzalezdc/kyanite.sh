'use client';

import { useState, useRef, useEffect } from 'react';
import Recorder from './components/Recorder';
import AnalysisDisplay from './components/AnalysisDisplay';
import PlaybackControls from './components/PlaybackControls';
import Metronome from './components/Metronome';
import InstrumentPicker from './components/InstrumentPicker';
import SectionManager from './components/SectionManager';
import ExportPanel from './components/ExportPanel';
import LyricsAssistant from './components/LyricsAssistant';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { BPMDetector } from '@/lib/audio/bpm-detector';
import { KeyDetector } from '@/lib/audio/key-detector';
import { MusicGenerator } from '@/lib/audio/music-generator';
import { PitchPoint, BPMAnalysis, KeyAnalysis, Section, InstrumentType } from '@/lib/types';

export default function Home() {
  const [audioBuffer, setAudioBuffer] = useState<AudioBuffer | null>(null);
  const [pitches, setPitches] = useState<PitchPoint[]>([]);
  const [bpm, setBpm] = useState<BPMAnalysis | null>(null);
  const [key, setKey] = useState<KeyAnalysis | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [sections, setSections] = useState<Section[]>([]);
  const [selectedInstruments, setSelectedInstruments] = useState<InstrumentType[]>(['drums', 'bass', 'chords']);
  const [isPlaying, setIsPlaying] = useState(false);
  const [arrangementMode, setArrangementMode] = useState<'sequential' | 'layered'>('sequential');
  
  const generatorRef = useRef<MusicGenerator | null>(null);

  useEffect(() => {
    // Initialize generator
    const generator = new MusicGenerator();
    generator.initialize().then(() => {
      generatorRef.current = generator;
    });

    return () => {
      generator?.dispose();
    };
  }, []);

  const handleRecordingComplete = async (buffer: AudioBuffer) => {
    setAudioBuffer(buffer);
    setIsAnalyzing(true);

    try {
      // Run pitch detection
      const pitchDetector = new PitchDetector(buffer.sampleRate);
      const detectedPitches = pitchDetector.analyze(buffer);
      setPitches(detectedPitches);

      // Run BPM detection
      const bpmDetector = new BPMDetector();
      const bpmAnalysis = await bpmDetector.analyze(buffer);
      setBpm(bpmAnalysis);

      // Run key detection
      const keyDetector = new KeyDetector();
      const keyAnalysis = keyDetector.analyze(detectedPitches);
      setKey(keyAnalysis);

      // Create section from recording
      const section: Section = {
        id: `section-${Date.now()}`,
        name: `Section ${sections.length + 1}`,
        audioBuffer: buffer,
        duration: buffer.duration,
        pitches: detectedPitches,
        bpm: bpmAnalysis,
        key: keyAnalysis,
        instruments: [...selectedInstruments],
      };

      setSections([...sections, section]);
    } catch (error) {
      console.error('Analysis failed:', error);
      alert('Failed to analyze audio. Please try again.');
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleGenerateMusic = () => {
    if (!generatorRef.current || !bpm || !key) {
      alert('Please record and analyze audio first');
      return;
    }

    const generator = generatorRef.current;
    generator.setBPM(bpm.bpm);
    generator.setKey(key.key);
    generator.setInstruments(selectedInstruments);
    generator.generateDrums('moderate');
    generator.generateBass(key.key);
    generator.generateChords(key.key);
  };

  const handlePlay = () => {
    if (!generatorRef.current) return;
    
    handleGenerateMusic();
    generatorRef.current.start();
    setIsPlaying(true);
  };

  const handlePause = () => {
    if (!generatorRef.current) return;
    generatorRef.current.stop();
    setIsPlaying(false);
  };

  const handleStop = () => {
    if (!generatorRef.current) return;
    generatorRef.current.stop();
    setIsPlaying(false);
  };

  const handleDeleteSection = (id: string) => {
    setSections(sections.filter(s => s.id !== id));
  };

  const handleRenameSection = (id: string, name: string) => {
    setSections(sections.map(s => s.id === id ? { ...s, name } : s));
  };

  const handleNewSection = () => {
    setAudioBuffer(null);
    setPitches([]);
    setBpm(null);
    setKey(null);
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

        {/* Recording Section */}
        <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">Step 1: Record Your Voice</h2>
            {sections.length > 0 && (
              <button
                onClick={handleNewSection}
                className="px-4 py-2 bg-primary-500/20 border border-primary-500/50 rounded-lg hover:bg-primary-500/30 transition-colors text-sm"
              >
                New Section
              </button>
            )}
          </div>
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

        {/* Analysis Display */}
        {pitches.length > 0 && !isAnalyzing && (
          <AnalysisDisplay 
            pitches={pitches} 
            bpm={bpm?.bpm || null}
            musicalKey={key?.key || null}
          />
        )}

        {/* Music Generation */}
        {bpm && key && !isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-6">
            <h2 className="text-xl font-semibold">Step 2: Generate Music</h2>
            
            <InstrumentPicker
              selected={selectedInstruments}
              onChange={setSelectedInstruments}
            />

            <div className="space-y-4">
              <button
                onClick={handleGenerateMusic}
                className="w-full px-6 py-3 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors"
              >
                Generate Music
              </button>

              <PlaybackControls
                isPlaying={isPlaying}
                onPlay={handlePlay}
                onPause={handlePause}
                onStop={handleStop}
              />

              {bpm && (
                <Metronome bpm={bpm.bpm} isPlaying={isPlaying} />
              )}
            </div>
          </div>
        )}

        {/* Section Management */}
        {sections.length > 0 && (
          <SectionManager
            sections={sections}
            onDelete={handleDeleteSection}
            onRename={handleRenameSection}
          />
        )}

        {/* Arrangement Mode */}
        {sections.length > 1 && (
          <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
            <h2 className="text-xl font-semibold mb-4">Arrangement Mode</h2>
            <div className="flex gap-4">
              <button
                onClick={() => setArrangementMode('sequential')}
                className={`px-4 py-2 rounded-lg border transition-colors ${
                  arrangementMode === 'sequential'
                    ? 'bg-primary-500/20 border-primary-500 text-primary-500'
                    : 'bg-gray-800 border-gray-700 text-gray-400'
                }`}
              >
                Sequential
              </button>
              <button
                onClick={() => setArrangementMode('layered')}
                className={`px-4 py-2 rounded-lg border transition-colors ${
                  arrangementMode === 'layered'
                    ? 'bg-primary-500/20 border-primary-500 text-primary-500'
                    : 'bg-gray-800 border-gray-700 text-gray-400'
                }`}
              >
                Layered
              </button>
            </div>
          </div>
        )}

        {/* Export Panel */}
        {pitches.length > 0 && (
          <ExportPanel
            audioBuffer={audioBuffer}
            pitches={pitches}
            bpm={bpm}
            key={key?.key || null}
            generator={generatorRef.current}
          />
        )}

        {/* AI Lyrics Assistant */}
        {pitches.length > 0 && key && (
          <LyricsAssistant
            pitches={pitches}
            musicalKey={key.key}
          />
        )}
      </div>
    </main>
  );
}
