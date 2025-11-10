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
import OnboardingTour from './components/OnboardingTour';
import WelcomeModal from './components/WelcomeModal';
import QuickStartGuide from './components/QuickStartGuide';
import MobileNavigation from './components/MobileNavigation';
import ResponsiveGrid, { ResponsiveCardGrid, ResponsiveStatsGrid } from './components/ResponsiveGrid';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { BPMDetector } from '@/lib/audio/bpm-detector';
import { KeyDetector } from '@/lib/audio/key-detector';
import { TimeSignatureDetector } from '@/lib/audio/time-signature-detector';
import { MusicGenerator } from '@/lib/audio/music-generator';
import PianoRoll from './components/PianoRoll';
import VisualizerCanvas from './components/pixi/VisualizerCanvas';
import PianoRollEditor from './components/pixi/PianoRollEditor';
import RhythmGameMode from './components/pixi/RhythmGameMode';
import { PitchPoint, BPMAnalysis, KeyAnalysis, Section, InstrumentType, TimeSignature } from '@/lib/types';
import { AnalyticsProvider, AnalyticsErrorBoundaryWithTracking } from './components/AnalyticsProvider';
import { AnalyticsDashboardDev } from './components/AnalyticsDashboard';
import {
  useTrackInteraction,
  useTrackAudio,
  useTrackFeature,
  useTrackPerformance,
  useAudioPerformanceMeasure,
  useJourneyTracking
} from '@/lib/analytics/hooks';
import {
  useMobileFeatures,
  useMobileNavigation,
  useResponsiveBreakpoints
} from '@/lib/hooks/useMobileFeatures';
import { usePerformanceOptimization } from '@/lib/hooks/usePerformanceOptimization';

function HomeContent() {
  const [audioBuffer, setAudioBuffer] = useState<AudioBuffer | null>(null);
  const [pitches, setPitches] = useState<PitchPoint[]>([]);
  const [bpm, setBpm] = useState<BPMAnalysis | null>(null);
  const [key, setKey] = useState<KeyAnalysis | null>(null);
  const [timeSignature, setTimeSignature] = useState<TimeSignature | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [sections, setSections] = useState<Section[]>([]);
  const [selectedInstruments, setSelectedInstruments] = useState<InstrumentType[]>(['drums', 'bass', 'chords']);
  const [isPlaying, setIsPlaying] = useState(false);
  const [arrangementMode, setArrangementMode] = useState<'sequential' | 'layered'>('sequential');
  const [pixiMode, setPixiMode] = useState<'record' | 'edit' | 'game' | null>(null);
  
  // Onboarding state
  const [showWelcomeModal, setShowWelcomeModal] = useState(false);
  const [showTour, setShowTour] = useState(false);
  const [hasSeenOnboarding, setHasSeenOnboarding] = useState(true);
  
  const generatorRef = useRef<MusicGenerator | null>(null);

  // Mobile and performance hooks
  const {
    isMobile,
    isTouchDevice,
    orientation,
    triggerHaptic,
    scrollToSection
  } = useMobileFeatures();
  
  const {
    activeSection,
    setActiveSection,
    isVisible: navVisible
  } = useMobileNavigation();
  
  const { breakpoint, isSmallScreen } = useResponsiveBreakpoints();
  
  const {
    shouldReduceAnimations,
    shouldLazyLoad,
    measurePerformance
  } = usePerformanceOptimization();

  // Analytics hooks
  const trackInteraction = useTrackInteraction();
  const trackAudio = useTrackAudio();
  const trackFeature = useTrackFeature();
  const trackPerformance = useTrackPerformance();
  const measureAudioProcessing = useAudioPerformanceMeasure('audio_analysis');
  const { trackStepStart, trackStepComplete, trackJourneyComplete } = useJourneyTracking('voxforge_workflow');

  useEffect(() => {
    // Check if user has seen onboarding
    const onboardingStatus = localStorage.getItem('voxforge-onboarding-complete');
    if (!onboardingStatus) {
      setHasSeenOnboarding(false);
      setShowWelcomeModal(true);
    }

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
    // Track recording completion
    trackAudio('record_complete', {
      duration: buffer.duration,
      sampleRate: buffer.sampleRate
    });

    // Clear previous analysis state
    setPitches([]);
    setBpm(null);
    setKey(null);
    setTimeSignature(null);
    
    setAudioBuffer(buffer);
    setIsAnalyzing(true);

    // Track analysis start
    trackStepStart('audio_analysis');
    trackAudio('analysis_start');

    try {
      // Run pitch detection with performance tracking
      const detectedPitches = await measureAudioProcessing(async () => {
        const pitchDetector = new PitchDetector(buffer.sampleRate);
        return pitchDetector.analyze(buffer);
      });
      setPitches(detectedPitches);

      // Run BPM detection with performance tracking
      const bpmAnalysis = await measureAudioProcessing(async () => {
        const bpmDetector = new BPMDetector();
        return bpmDetector.analyze(buffer);
      });
      setBpm(bpmAnalysis);

      // Run key detection with performance tracking
      const keyAnalysis = await measureAudioProcessing(async () => {
        const keyDetector = new KeyDetector();
        return keyDetector.analyze(detectedPitches);
      });
      setKey(keyAnalysis);

      // Run time signature detection with performance tracking
      const timeSig = await measureAudioProcessing(async () => {
        const timeSigDetector = new TimeSignatureDetector();
        return timeSigDetector.analyze(detectedPitches, (bpmAnalysis as BPMAnalysis).bpm);
      });
      setTimeSignature(timeSig);

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
      
      // Track successful analysis
      trackStepComplete('audio_analysis');
      trackAudio('analysis_complete', {
        duration: buffer.duration,
        success: true
      });
    } catch (error) {
      console.error('Analysis failed:', error);
      trackAudio('analysis_complete', {
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
      alert('Failed to analyze audio. Please try again.');
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleGenerateMusic = () => {
    if (!generatorRef.current || !bpm || !key) {
      trackInteraction('click', 'generate_music_error');
      alert('Please record and analyze audio first');
      return;
    }

    // Track music generation
    trackStepStart('music_generation');
    trackFeature('generation', 'start', {
      bpm: bpm.bpm,
      key: key.key,
      instruments: selectedInstruments
    });

    const generator = generatorRef.current;
    generator.setBPM(bpm.bpm);
    generator.setKey(key.key);
    generator.setInstruments(selectedInstruments);
    generator.generateDrums('moderate');
    generator.generateBass(key.key);
    generator.generateChords(key.key);
    
    trackStepComplete('music_generation');
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

  const handleBPMChange = (newBPM: number) => {
    if (bpm) {
      setBpm({ ...bpm, bpm: newBPM });
    }
    if (generatorRef.current) {
      generatorRef.current.setBPM(newBPM);
    }
  };

  const handleKeyChange = (newKey: string) => {
    if (key) {
      setKey({ ...key, key: newKey });
    }
    if (generatorRef.current) {
      generatorRef.current.setKey(newKey);
    }
  };

  const handleTimeSignatureChange = (newTimeSig: TimeSignature) => {
    setTimeSignature(newTimeSig);
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

  // Onboarding handlers
  const handleStartTour = () => {
    trackFeature('onboarding', 'start');
    setShowWelcomeModal(false);
    setShowTour(true);
  };

  const handleJumpIn = () => {
    trackFeature('onboarding', 'skip');
    setShowWelcomeModal(false);
    setHasSeenOnboarding(true);
    localStorage.setItem('voxforge-onboarding-complete', 'true');
  };

  const handleTourComplete = () => {
    trackFeature('onboarding', 'complete');
    trackJourneyComplete({
      tour_completed: true
    });
    setShowTour(false);
    setHasSeenOnboarding(true);
    localStorage.setItem('voxforge-onboarding-complete', 'true');
  };

  const handleTourSkip = () => {
    trackFeature('onboarding', 'skip');
    setShowTour(false);
    setHasSeenOnboarding(true);
    localStorage.setItem('voxforge-onboarding-complete', 'true');
  };

  const handleRestartTour = () => {
    trackInteraction('click', 'restart_tour');
    setShowTour(true);
  };

  return (
    <main className={`
      min-h-screen-safe pb-nav-h md:pb-0
      ${isMobile ? 'px-4' : 'px-6 sm:px-8'}
      ${shouldReduceAnimations() ? 'reduce-animation' : ''}
    `}>
      {/* Onboarding Components */}
      <WelcomeModal
        isOpen={showWelcomeModal}
        onStartTour={handleStartTour}
        onJumpIn={handleJumpIn}
      />
      
      <OnboardingTour
        isOpen={showTour}
        onComplete={handleTourComplete}
        onSkip={handleTourSkip}
      />

      <div className="max-w-4xl mx-auto space-y-6 sm:space-y-8">
        <header className="text-center space-y-2 py-4 sm:py-6">
          <h1 className="fluid-4xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            VoxForge
          </h1>
          <p className="fluid-base text-gray-400">
            Transform your voice into music
          </p>
          {hasSeenOnboarding && (
            <button
              onClick={() => {
                handleRestartTour();
                triggerHaptic('light');
              }}
              className="mt-2 fluid-sm text-primary-500 hover:text-primary-400 transition-colors min-h-[44px] px-4 py-2 rounded-lg touch-manipulation active:scale-95"
            >
              Take a tour
            </button>
          )}
        </header>

        {/* Quick Start Guide */}
        {!hasSeenOnboarding && <QuickStartGuide />}

        {/* Recording Section */}
        <div className="bg-gray-900 rounded-xl p-4 sm:p-6 lg:p-8 border border-gray-800" id="record">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4">
            <h2 className="fluid-xl font-semibold">Step 1: Record Your Voice</h2>
            {sections.length > 0 && (
              <button
                onClick={() => {
                  handleNewSection();
                  triggerHaptic('light');
                }}
                className="px-4 py-2 bg-primary-500/20 border border-primary-500/50 rounded-lg hover:bg-primary-500/30 transition-colors fluid-sm min-h-[44px] touch-manipulation active:scale-95"
              >
                New Section
              </button>
            )}
          </div>
          <Recorder onRecordingComplete={handleRecordingComplete} />
        </div>

        {isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-6 sm:p-8 border border-gray-800">
            <div className="flex items-center justify-center gap-3">
              <div className={`animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500 ${shouldReduceAnimations() ? 'animate-pulse' : ''}`}></div>
              <p className="fluid-base text-gray-400">Analyzing audio...</p>
            </div>
          </div>
        )}

        {/* Analysis Display */}
        {pitches.length > 0 && !isAnalyzing && (
          <>
            <AnalysisDisplay 
              pitches={pitches} 
              bpm={bpm?.bpm || null}
              musicalKey={key?.key || null}
              timeSignature={timeSignature}
              onBPMChange={handleBPMChange}
              onKeyChange={handleKeyChange}
              onTimeSignatureChange={handleTimeSignatureChange}
            />
            
            {/* Pixi.js Mode Selector */}
            <div className="bg-gray-900 rounded-xl p-4 border border-gray-800">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-lg font-semibold">Interactive Visualizations</h3>
                <div className="flex gap-2">
                  <button
                    onClick={() => setPixiMode(pixiMode === 'record' ? null : 'record')}
                    className={`px-4 py-2 rounded-lg transition-colors text-sm ${
                      pixiMode === 'record'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                    }`}
                  >
                    Real-Time Visualizer
                  </button>
                  <button
                    onClick={() => setPixiMode(pixiMode === 'edit' ? null : 'edit')}
                    className={`px-4 py-2 rounded-lg transition-colors text-sm ${
                      pixiMode === 'edit'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                    }`}
                  >
                    Piano Roll Editor
                  </button>
                  <button
                    onClick={() => setPixiMode(pixiMode === 'game' ? null : 'game')}
                    className={`px-4 py-2 rounded-lg transition-colors text-sm ${
                      pixiMode === 'game'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                    }`}
                  >
                    Rhythm Game
                  </button>
                </div>
              </div>
              
              {pixiMode === 'record' && (
                <VisualizerCanvas 
                  pitches={pitches}
                  isRecording={false}
                />
              )}
              
              {pixiMode === 'edit' && bpm && (
                <PianoRollEditor
                  initialPitches={pitches}
                  bpm={bpm.bpm}
                  onExport={(newPitches) => {
                    setPitches(newPitches);
                  }}
                />
              )}
              
              {pixiMode === 'game' && bpm && (
                <RhythmGameMode
                  bpm={bpm.bpm}
                  onComplete={(notes) => {
                    // Convert game notes to pitches
                    const gamePitches: PitchPoint[] = notes.map(note => ({
                      frequency: 440 * Math.pow(2, (note.midi - 69) / 12),
                      time: note.time,
                      midi: note.midi,
                      confidence: 1
                    }));
                    setPitches([...pitches, ...gamePitches]);
                  }}
                />
              )}
            </div>
            
            {audioBuffer && !pixiMode && (
              <PianoRoll
                pitches={pitches}
                duration={audioBuffer.duration}
                bpm={bpm?.bpm || null}
              />
            )}
          </>
        )}

        {/* Music Generation */}
        {bpm && key && !isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-6 sm:p-8 border border-gray-800 space-y-6" id="generate" data-tour="instruments">
            <h2 className="fluid-xl font-semibold">Step 2: Generate Music</h2>
            
            <InstrumentPicker
              selected={selectedInstruments}
              onChange={setSelectedInstruments}
            />

            <div className="space-y-4">
              <button
                onClick={() => {
                  handleGenerateMusic();
                  triggerHaptic('medium');
                }}
                className="w-full px-6 py-4 sm:py-3 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-all duration-200 min-h-[44px] touch-manipulation active:scale-95"
                data-tour="generate"
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
            musicalKey={key?.key || null}
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

export default function Home() {
  return (
    <AnalyticsErrorBoundaryWithTracking>
      <HomeContent />
      <AnalyticsDashboardDev />
    </AnalyticsErrorBoundaryWithTracking>
  );
}
