'use client';

import { useRef, useEffect } from 'react';
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

// Import store hooks
import { StoreProvider, HydrationLoader } from './components/StoreProvider';
import {
  useAudioStore,
  useRecordingState,
  useAnalysisState,
  useGeneratedAudioState,
  useAudioActions,
  useNavigationState,
  useModalState,
  useThemePreferences,
  useLayoutPreferences,
  useVisualizationState,
  useUIActions,
  useOnboardingState,
  useUserActions,
  useProjectState,
  useProjectActions,
  useAudioAndAnalysis,
  useModalsAndNavigation,
  useOnboardingAndHelp,
  useMobileState,
  usePersistenceState,
  useAnalyticsState
} from '@/lib/store/hooks';

function HomeContent() {
  // Store hooks for state
  const recording = useRecordingState();
  const analysis = useAnalysisState();
  const generated = useGeneratedAudioState();
  const audioActions = useAudioActions();
  const audioStore = useAudioStore(); // Get the full store for accessing state
  
  const navigation = useNavigationState();
  const modals = useModalState();
  const theme = useThemePreferences();
  const layout = useLayoutPreferences();
  const visualization = useVisualizationState();
  const uiActions = useUIActions();
  
  const onboarding = useOnboardingState();
  const userActions = useUserActions();
  
  const project = useProjectState();
  const projectActions = useProjectActions();
  
  // Combined hooks for convenience
  const audioAndAnalysis = useAudioAndAnalysis();
  const modalsAndNavigation = useModalsAndNavigation();
  const onboardingAndHelp = useOnboardingAndHelp();
  const mobileState = useMobileState();
  const persistenceState = usePersistenceState();
  const analyticsState = useAnalyticsState();
  
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
    // Check if user has seen onboarding using store state
    if (!onboarding.hasSeenOnboarding) {
      uiActions.setWelcomeModal(true);
    }

    // Initialize generator
    const generator = new MusicGenerator();
    generator.initialize().then(() => {
      generatorRef.current = generator;
    });

    return () => {
      generator?.dispose();
    };
  }, [onboarding.hasSeenOnboarding, uiActions]);

  const handleRecordingComplete = async (buffer: AudioBuffer) => {
    // Track recording completion
    trackAudio('record_complete', {
      duration: buffer.duration,
      sampleRate: buffer.sampleRate
    });

    // Clear previous analysis state
    audioActions.resetAnalysis();
    audioActions.setRecordedAudio(buffer);
    audioActions.setAnalyzing(true);

    // Track analysis start
    trackStepStart('audio_analysis');
    trackAudio('analysis_start');

    try {
      // Run pitch detection with performance tracking
      const detectedPitches = await measureAudioProcessing(async () => {
        const pitchDetector = new PitchDetector(buffer.sampleRate);
        return pitchDetector.analyze(buffer);
      });
      audioActions.setPitches(detectedPitches);

      // Run BPM detection with performance tracking
      const bpmAnalysis = await measureAudioProcessing(async () => {
        const bpmDetector = new BPMDetector();
        return bpmDetector.analyze(buffer);
      });
      audioActions.setBPM(bpmAnalysis);

      // Run key detection with performance tracking
      const keyAnalysis = await measureAudioProcessing(async () => {
        const keyDetector = new KeyDetector();
        return keyDetector.analyze(detectedPitches);
      });
      audioActions.setKey(keyAnalysis);

      // Run time signature detection with performance tracking
      const timeSig = await measureAudioProcessing(async () => {
        const timeSigDetector = new TimeSignatureDetector();
        return timeSigDetector.analyze(detectedPitches, (bpmAnalysis as BPMAnalysis).bpm);
      });
      audioActions.setTimeSignature(timeSig);

      // Create section from recording
      const section: Section = {
        id: `section-${Date.now()}`,
        name: `Section ${audioStore.sections?.length || 0 + 1}`,
        audioBuffer: buffer,
        duration: buffer.duration,
        pitches: detectedPitches,
        bpm: bpmAnalysis,
        key: keyAnalysis,
        instruments: [...audioStore.selectedInstruments],
      };

      audioActions.addSection(section);
      
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
      audioActions.setAnalyzing(false);
    }
  };

  const handleGenerateMusic = () => {
    if (!generatorRef.current || !analysis.bpm || !analysis.key) {
      trackInteraction('click', 'generate_music_error');
      alert('Please record and analyze audio first');
      return;
    }

    // Track music generation
    trackStepStart('music_generation');
    trackFeature('generation', 'start', {
      bpm: analysis.bpm.bpm,
      key: analysis.key.key,
      instruments: audioStore.selectedInstruments
    });

    const generator = generatorRef.current;
    generator.setBPM(analysis.bpm.bpm);
    generator.setKey(analysis.key.key);
    generator.setInstruments(audioStore.selectedInstruments);
    generator.generateDrums('moderate');
    generator.generateBass(analysis.key.key);
    generator.generateChords(analysis.key.key);
    
    trackStepComplete('music_generation');
  };

  const handlePlay = () => {
    if (!generatorRef.current) return;
    
    handleGenerateMusic();
    generatorRef.current.start();
    audioActions.setGeneratedPlaying(true);
  };

  const handlePause = () => {
    if (!generatorRef.current) return;
    generatorRef.current.stop();
    audioActions.setGeneratedPlaying(false);
  };

  const handleStop = () => {
    if (!generatorRef.current) return;
    generatorRef.current.stop();
    audioActions.setGeneratedPlaying(false);
  };

  const handleBPMChange = (newBPM: number) => {
    if (analysis.bpm) {
      audioActions.setBPM({ ...analysis.bpm, bpm: newBPM });
    }
    if (generatorRef.current) {
      generatorRef.current.setBPM(newBPM);
    }
  };

  const handleKeyChange = (newKey: string) => {
    if (analysis.key) {
      audioActions.setKey({ ...analysis.key, key: newKey });
    }
    if (generatorRef.current) {
      generatorRef.current.setKey(newKey);
    }
  };

  const handleTimeSignatureChange = (newTimeSig: TimeSignature) => {
    audioActions.setTimeSignature(newTimeSig);
  };

  const handleDeleteSection = (id: string) => {
    audioActions.deleteSection(id);
  };

  const handleRenameSection = (id: string, name: string) => {
    audioActions.updateSection(id, { name });
  };

  const handleNewSection = () => {
    audioActions.resetRecording();
    audioActions.resetAnalysis();
  };

  // Onboarding handlers
  const handleStartTour = () => {
    trackFeature('onboarding', 'start');
    uiActions.setWelcomeModal(false);
    uiActions.setOnboardingTour(true);
  };

  const handleJumpIn = () => {
    trackFeature('onboarding', 'skip');
    uiActions.setWelcomeModal(false);
    userActions.completeOnboarding();
  };

  const handleTourComplete = () => {
    trackFeature('onboarding', 'complete');
    trackJourneyComplete({
      tour_completed: true
    });
    uiActions.setOnboardingTour(false);
    userActions.completeOnboarding();
  };

  const handleTourSkip = () => {
    trackFeature('onboarding', 'skip');
    uiActions.setOnboardingTour(false);
    userActions.completeOnboarding();
  };

  const handleRestartTour = () => {
    trackInteraction('click', 'restart_tour');
    uiActions.setOnboardingTour(true);
  };

  return (
    <main className={`
      min-h-screen-safe pb-nav-h md:pb-0
      ${isMobile ? 'px-4' : 'px-6 sm:px-8'}
      ${shouldReduceAnimations() ? 'reduce-animation' : ''}
    `}>
      {/* Onboarding Components */}
      <WelcomeModal />
      
      <OnboardingTour
        isOpen={modals.onboardingTour}
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
          {onboarding.hasSeenOnboarding && (
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
        {!onboarding.hasSeenOnboarding && <QuickStartGuide />}

        {/* Recording Section */}
        <div className="bg-gray-900 rounded-xl p-4 sm:p-6 lg:p-8 border border-gray-800" id="record">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4">
            <h2 className="fluid-xl font-semibold">Step 1: Record Your Voice</h2>
            {audioStore.sections && audioStore.sections.length > 0 && (
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

        {analysis.isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-6 sm:p-8 border border-gray-800">
            <div className="flex items-center justify-center gap-3">
              <div className={`animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500 ${shouldReduceAnimations() ? 'animate-pulse' : ''}`}></div>
              <p className="fluid-base text-gray-400">Analyzing audio...</p>
            </div>
          </div>
        )}

        {/* Analysis Display */}
        {analysis.pitches.length > 0 && !analysis.isAnalyzing && (
          <>
            <AnalysisDisplay 
              pitches={analysis.pitches} 
              bpm={analysis.bpm?.bpm || null}
              musicalKey={analysis.key?.key || null}
              timeSignature={analysis.timeSignature}
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
                    onClick={() => uiActions.setPixiMode(visualization.pixiMode === 'record' ? null : 'record')}
                    className={`px-4 py-2 rounded-lg transition-colors text-sm ${
                      visualization.pixiMode === 'record'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                    }`}
                  >
                    Real-Time Visualizer
                  </button>
                  <button
                    onClick={() => uiActions.setPixiMode(visualization.pixiMode === 'edit' ? null : 'edit')}
                    className={`px-4 py-2 rounded-lg transition-colors text-sm ${
                      visualization.pixiMode === 'edit'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                    }`}
                  >
                    Piano Roll Editor
                  </button>
                  <button
                    onClick={() => uiActions.setPixiMode(visualization.pixiMode === 'game' ? null : 'game')}
                    className={`px-4 py-2 rounded-lg transition-colors text-sm ${
                      visualization.pixiMode === 'game'
                        ? 'bg-primary-500 text-white'
                        : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                    }`}
                  >
                    Rhythm Game
                  </button>
                </div>
              </div>
              
              {visualization.pixiMode === 'record' && (
                <VisualizerCanvas 
                  pitches={analysis.pitches}
                  isRecording={recording.isRecording}
                />
              )}
              
              {visualization.pixiMode === 'edit' && analysis.bpm && (
                <PianoRollEditor
                  initialPitches={analysis.pitches}
                  bpm={analysis.bpm.bpm}
                  onExport={(newPitches) => {
                    audioActions.setPitches(newPitches);
                  }}
                />
              )}
              
              {visualization.pixiMode === 'game' && analysis.bpm && (
                <RhythmGameMode
                  bpm={analysis.bpm.bpm}
                  onComplete={(notes) => {
                    // Convert game notes to pitches
                    const gamePitches: PitchPoint[] = notes.map(note => ({
                      frequency: 440 * Math.pow(2, (note.midi - 69) / 12),
                      time: note.time,
                      midi: note.midi,
                      confidence: 1
                    }));
                    audioActions.setPitches([...analysis.pitches, ...gamePitches]);
                  }}
                />
              )}
            </div>
            
            {recording.recordedAudio && !visualization.pixiMode && (
              <PianoRoll
                pitches={analysis.pitches}
                duration={recording.recordedAudio.duration}
                bpm={analysis.bpm?.bpm || null}
              />
            )}
          </>
        )}

        {/* Music Generation */}
        {analysis.bpm && analysis.key && !analysis.isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-6 sm:p-8 border border-gray-800 space-y-6" id="generate" data-tour="instruments">
            <h2 className="fluid-xl font-semibold">Step 2: Generate Music</h2>
            
            <InstrumentPicker
              selected={audioStore.selectedInstruments}
              onChange={audioActions.setSelectedInstruments}
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
                isPlaying={generated.generatedPlaying}
                onPlay={handlePlay}
                onPause={handlePause}
                onStop={handleStop}
              />

              {analysis.bpm && (
                <Metronome bpm={analysis.bpm.bpm} isPlaying={generated.generatedPlaying} />
              )}
            </div>
          </div>
        )}

        {/* Section Management */}
        {audioStore.sections && audioStore.sections.length > 0 && (
          <SectionManager
            sections={audioStore.sections}
            onDelete={handleDeleteSection}
            onRename={handleRenameSection}
          />
        )}

        {/* Arrangement Mode */}
        {audioStore.sections && audioStore.sections.length > 1 && (
          <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
            <h2 className="text-xl font-semibold mb-4">Arrangement Mode</h2>
            <div className="flex gap-4">
              <button
                onClick={() => audioActions.setArrangementMode('sequential')}
                className={`px-4 py-2 rounded-lg border transition-colors ${
                  generated.arrangementMode === 'sequential'
                    ? 'bg-primary-500/20 border-primary-500 text-primary-500'
                    : 'bg-gray-800 border-gray-700 text-gray-400'
                }`}
              >
                Sequential
              </button>
              <button
                onClick={() => audioActions.setArrangementMode('layered')}
                className={`px-4 py-2 rounded-lg border transition-colors ${
                  generated.arrangementMode === 'layered'
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
        {analysis.pitches.length > 0 && (
          <ExportPanel
            audioBuffer={recording.recordedAudio}
            pitches={analysis.pitches}
            bpm={analysis.bpm}
            musicalKey={analysis.key?.key || null}
            generator={generatorRef.current}
          />
        )}

        {/* AI Lyrics Assistant */}
        {analysis.pitches.length > 0 && analysis.key && (
          <LyricsAssistant
            pitches={analysis.pitches}
            musicalKey={analysis.key.key}
          />
        )}
      </div>
    </main>
  );
}

export default function Home() {
  return (
    <StoreProvider>
      <HydrationLoader>
        <AnalyticsErrorBoundaryWithTracking>
          <HomeContent />
          <AnalyticsDashboardDev />
        </AnalyticsErrorBoundaryWithTracking>
      </HydrationLoader>
    </StoreProvider>
  );
}
