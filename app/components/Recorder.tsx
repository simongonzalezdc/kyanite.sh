'use client';

import { useRef, useEffect } from 'react';
import { AudioRecorder } from '@/lib/audio/recorder';
import Waveform from './Waveform';
import Tooltip from './Tooltip';
import { Mic, Square, Play, Scissors, RotateCcw } from 'lucide-react';
import { trimAudioBuffer } from '@/lib/utils/audio-utils';
import { isMobile, isTouchDevice, triggerHaptic, enableSwipeGestures } from '@/lib/utils';

// Import store hooks
import {
  useRecordingState,
  useAudioActions
} from '@/lib/store/hooks';

// Import accessibility components
import { useAccessibility } from './AccessibilityProvider';
import { useKeyboardNavigation } from './KeyboardNavigation';
import { AudioVisualIndicator, AlternativeRecordingInput } from './AudioAccessibility';
import { ScreenReaderOnly, LiveRegion, AriaLabel } from './ScreenReaderSupport';

interface RecorderProps {
  onRecordingComplete: (audioBuffer: AudioBuffer) => void;
}

export default function Recorder({ onRecordingComplete }: RecorderProps) {
  // Use store hooks instead of local state
  const recording = useRecordingState();
  const audioActions = useAudioActions();
  
  // Accessibility hooks
  const { announce, settings } = useAccessibility();
  const { registerShortcut, unregisterShortcut } = useKeyboardNavigation();
  
  const recorderRef = useRef<AudioRecorder | null>(null);
  const audioContextRef = useRef<AudioContext | null>(null);

  const requestPermission = async () => {
    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }
    
    const granted = await recorderRef.current.requestPermission();
    audioActions.setPermission(granted);
    
    if (!granted) {
      alert('Microphone permission is required to record audio.');
    }
  };

  const startRecording = async () => {
    if (!recording.hasPermission) {
      await requestPermission();
      if (!recording.hasPermission) return;
    }

    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }

    try {
      await recorderRef.current.startRecording();
      audioActions.setRecording(true);
      audioActions.setRecordedAudio(null);
      
      // Announce to screen readers
      announce('Recording started');
    } catch (error) {
      console.error('Failed to start recording:', error);
      announce('Failed to start recording. Please check microphone permissions.');
      alert('Failed to start recording. Please check microphone permissions.');
    }
  };

  const stopRecording = async () => {
    if (!recorderRef.current || !recording.isRecording) return;

    try {
      const audioBuffer = await recorderRef.current.stopRecording();
      audioActions.setRecording(false);
      audioActions.setRecordedAudio(audioBuffer);
      audioActions.setOriginalAudio(audioBuffer);
      audioActions.setTrimStart(0);
      audioActions.setTrimEnd(audioBuffer.duration);
      audioActions.setTrimming(false);
      
      // Announce to screen readers
      announce(`Recording stopped. Duration: ${audioBuffer.duration.toFixed(2)} seconds`);
      
      onRecordingComplete(audioBuffer);
    } catch (error) {
      console.error('Failed to stop recording:', error);
      announce('Failed to process recording.');
      alert('Failed to process recording.');
    }
  };

  const playRecording = async () => {
    if (!recording.recordedAudio) return;

    if (!audioContextRef.current) {
      audioContextRef.current = new AudioContext();
    }

    const source = audioContextRef.current.createBufferSource();
    source.buffer = recording.recordedAudio;
    source.connect(audioContextRef.current.destination);
    
    source.onended = () => {
      audioActions.setPlaying(false);
      announce('Playback finished');
    };

    source.start();
    audioActions.setPlaying(true);
    
    // Announce to screen readers
    announce('Playing recording');
  };

  const handleTrim = async () => {
    // Always trim from originalAudio to allow re-trimming with different values
    const sourceAudio = recording.originalAudio || recording.recordedAudio;
    if (!sourceAudio) return;
    
    audioActions.setTrimming(true);
    try {
      const trimmed = await trimAudioBuffer(sourceAudio, recording.trimStart, recording.trimEnd);
      audioActions.setRecordedAudio(trimmed);
      // Ensure originalAudio is set for future resets
      if (!recording.originalAudio) {
        audioActions.setOriginalAudio(sourceAudio);
      }
      onRecordingComplete(trimmed);
    } catch (error) {
      console.error('Failed to trim audio:', error);
      alert('Failed to trim audio. Please check your trim values.');
    } finally {
      audioActions.setTrimming(false);
    }
  };

  const handleResetTrim = () => {
    if (!recording.originalAudio) return;
    audioActions.setRecordedAudio(recording.originalAudio);
    audioActions.setTrimStart(0);
    audioActions.setTrimEnd(recording.originalAudio.duration);
    onRecordingComplete(recording.originalAudio);
  };

  // Update trim end when audio changes
  useEffect(() => {
    if (recording.recordedAudio && !recording.isRecording) {
      audioActions.setTrimEnd(recording.recordedAudio.duration);
    }
  }, [recording.recordedAudio?.duration, recording.isRecording]); // Use specific property and remove audioActions

  // Add swipe gesture support for mobile
  const waveformRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    if (isTouchDevice() && waveformRef.current) {
      const cleanup = enableSwipeGestures(waveformRef.current, {
        onSwipeLeft: () => {
          if (recording.recordedAudio && !recording.isRecording && !recording.isPlaying) {
            playRecording();
          }
        },
        onSwipeRight: () => {
          if (recording.isRecording) {
            stopRecording();
          } else if (!recording.recordedAudio) {
            startRecording();
          }
        }
      });
      
      return cleanup;
    }
  }, [recording.isRecording, recording.recordedAudio?.duration, recording.isPlaying]); // Use specific property

  // Register keyboard shortcuts
  useEffect(() => {
    registerShortcut({
      key: 'r',
      description: 'Start/stop recording',
      action: () => {
        if (recording.isRecording) {
          stopRecording();
        } else {
          startRecording();
        }
      },
      category: 'recording'
    });

    registerShortcut({
      key: 'p',
      description: 'Play recording',
      action: () => {
        if (recording.recordedAudio && !recording.isRecording) {
          playRecording();
        }
      },
      category: 'recording'
    });

    return () => {
      unregisterShortcut('r');
      unregisterShortcut('p');
    };
  }, [recording.isRecording, recording.recordedAudio?.duration]); // Use specific property and remove functions

  return (
    <div className="space-y-4" data-tour="recorder" id="record" role="region" aria-labelledby="recorder-heading">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-2">
        <h2 id="recorder-heading" className="fluid-xl font-semibold flex items-center gap-2">
          Voice Recorder
          <Tooltip content="Record your voice, melody, or any sound. The recorder will capture audio from your microphone and analyze it for pitch and timing." />
        </h2>
      </div>
      
      <div ref={waveformRef} className="touch-manipulation">
        <Waveform
          isRecording={recording.isRecording}
          getWaveformData={recording.isRecording ? () => recorderRef.current!.getWaveformData() : undefined}
          audioBuffer={recording.recordedAudio}
        />
        
        {/* Audio visual indicator for accessibility */}
        <AudioVisualIndicator
          isActive={recording.isRecording}
          type="recording"
          intensity="medium"
        />
      </div>

      <div className="flex flex-col sm:flex-row gap-4 justify-center">
        {!recording.isRecording ? (
          <AriaLabel
            label="Start recording"
            description="Begin recording audio from microphone"
          >
            <button
              onClick={() => {
                startRecording();
                triggerHaptic('medium');
              }}
              className={`
                flex items-center justify-center gap-2 px-6 py-4 sm:py-3
                bg-primary-500 hover:bg-primary-600
                rounded-lg font-medium transition-all duration-200
                min-h-[44px] min-w-[44px] sm:min-w-0
                ${isMobile() ? 'text-lg px-8 py-6' : ''}
                active:scale-95 touch-manipulation
                focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2
              `}
              style={{ minHeight: isMobile() ? '60px' : '48px' }}
              aria-pressed={recording.isRecording}
            >
              <Mic size={isMobile() ? 24 : 20} aria-hidden="true" />
              <span>Start Recording</span>
            </button>
          </AriaLabel>
        ) : (
          <AriaLabel
            label="Stop recording"
            description="Stop recording and process audio"
          >
            <button
              onClick={() => {
                stopRecording();
                triggerHaptic('heavy');
              }}
              className={`
                flex items-center justify-center gap-2 px-6 py-4 sm:py-3
                bg-red-500 hover:bg-red-600
                rounded-lg font-medium transition-all duration-200
                min-h-[44px] min-w-[44px] sm:min-w-0
                ${isMobile() ? 'text-lg px-8 py-6 animate-pulse' : ''}
                active:scale-95 touch-manipulation
                focus-visible:ring-2 focus-visible:ring-red-500 focus-visible:ring-offset-2
              `}
              style={{ minHeight: isMobile() ? '60px' : '48px' }}
              aria-pressed={recording.isRecording}
            >
              <Square size={isMobile() ? 24 : 20} aria-hidden="true" />
              <span>Stop Recording</span>
            </button>
          </AriaLabel>
        )}

        {recording.recordedAudio && !recording.isRecording && (
          <AriaLabel
            label="Play recording"
            description="Play back the recorded audio"
          >
            <button
              onClick={() => {
                playRecording();
                triggerHaptic('light');
              }}
              disabled={recording.isPlaying}
              className={`
                flex items-center justify-center gap-2 px-6 py-4 sm:py-3
                bg-secondary-500 hover:bg-secondary-600
                rounded-lg font-medium transition-all duration-200
                disabled:opacity-50 disabled:cursor-not-allowed
                min-h-[44px] min-w-[44px] sm:min-w-0
                ${isMobile() ? 'text-lg' : ''}
                active:scale-95 touch-manipulation
                focus-visible:ring-2 focus-visible:ring-secondary-500 focus-visible:ring-offset-2
              `}
              style={{ minHeight: isMobile() ? '48px' : '44px' }}
              aria-pressed={recording.isPlaying}
            >
              <Play size={isMobile() ? 24 : 20} aria-hidden="true" />
              <span>Play Back</span>
            </button>
          </AriaLabel>
        )}
      </div>

      {recording.recordedAudio && !recording.isRecording && (
        <div className="space-y-4">
          <div className="text-center fluid-sm text-gray-400" role="status" aria-live="polite">
            Duration: {recording.recordedAudio.duration.toFixed(2)}s
            {recording.originalAudio && recording.originalAudio !== recording.recordedAudio && (
              <span className="ml-2 text-primary-500 block sm:inline">
                (Trimmed from {recording.originalAudio.duration.toFixed(2)}s)
              </span>
            )}
          </div>

          {/* Trim Controls */}
          <div className="bg-gray-800 rounded-lg p-4 sm:p-6 border border-gray-700 space-y-4" role="region" aria-labelledby="trim-heading">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
              <h3 id="trim-heading" className="fluid-sm font-medium flex items-center gap-2">
                <Scissors size={16} className="text-secondary-500" aria-hidden="true" />
                Trim Audio
                <Tooltip content="Remove unwanted parts from your recording. Adjust start and end points to keep only the best portion of your audio." />
              </h3>
              {recording.originalAudio && recording.originalAudio !== recording.recordedAudio && (
                <button
                  onClick={() => {
                    handleResetTrim();
                    triggerHaptic('light');
                  }}
                  className="flex items-center gap-1 px-3 py-2 text-xs sm:text-sm bg-gray-700 hover:bg-gray-600 rounded transition-colors min-h-[44px] touch-manipulation active:scale-95"
                >
                  <RotateCcw size={14} />
                  Reset
                </button>
              )}
            </div>

            <div className="space-y-4">
              <div>
                <label htmlFor="trim-start" className="block fluid-xs text-gray-400 mb-2">
                  Start: {recording.trimStart.toFixed(2)}s
                </label>
                <input
                  id="trim-start"
                  type="range"
                  min="0"
                  max={recording.originalAudio?.duration || recording.recordedAudio.duration}
                  step="0.01"
                  value={recording.trimStart}
                  onChange={(e) => {
                    const newStart = parseFloat(e.target.value);
                    audioActions.setTrimStart(Math.min(newStart, recording.trimEnd - 0.1));
                    announce(`Trim start: ${newStart.toFixed(2)} seconds`);
                  }}
                  className="w-full h-3 sm:h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-primary-500 touch-manipulation focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2"
                  style={{ minHeight: isMobile() ? '24px' : '16px' }}
                  aria-label={`Trim start position: ${recording.trimStart.toFixed(2)} seconds`}
                  aria-valuemin={0}
                  aria-valuemax={recording.originalAudio?.duration || recording.recordedAudio.duration}
                  aria-valuenow={recording.trimStart}
                />
              </div>

              <div>
                <label htmlFor="trim-end" className="block fluid-xs text-gray-400 mb-2">
                  End: {recording.trimEnd.toFixed(2)}s
                </label>
                <input
                  id="trim-end"
                  type="range"
                  min={recording.trimStart + 0.1}
                  max={recording.originalAudio?.duration || recording.recordedAudio.duration}
                  step="0.01"
                  value={recording.trimEnd}
                  onChange={(e) => {
                    const newEnd = parseFloat(e.target.value);
                    audioActions.setTrimEnd(Math.max(newEnd, recording.trimStart + 0.1));
                    announce(`Trim end: ${newEnd.toFixed(2)} seconds`);
                  }}
                  className="w-full h-3 sm:h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-primary-500 touch-manipulation focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2"
                  style={{ minHeight: isMobile() ? '24px' : '16px' }}
                  aria-label={`Trim end position: ${recording.trimEnd.toFixed(2)} seconds`}
                  aria-valuemin={recording.trimStart + 0.1}
                  aria-valuemax={recording.originalAudio?.duration || recording.recordedAudio.duration}
                  aria-valuenow={recording.trimEnd}
                />
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between fluid-xs text-gray-400 gap-2">
                <span>Trimmed length: {(recording.trimEnd - recording.trimStart).toFixed(2)}s</span>
                <span>
                  {recording.originalAudio && (
                    <>Remove: {(recording.originalAudio.duration - (recording.trimEnd - recording.trimStart)).toFixed(2)}s</>
                  )}
                </span>
              </div>

              <button
                onClick={() => {
                  handleTrim();
                  triggerHaptic('medium');
                }}
                disabled={recording.isTrimming || !recording.recordedAudio || (recording.trimStart === 0 && recording.trimEnd === (recording.originalAudio?.duration || recording.recordedAudio.duration))}
                className="w-full px-4 py-3 sm:py-2 bg-secondary-500 hover:bg-secondary-600 rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 fluid-sm min-h-[44px] touch-manipulation active:scale-95 focus-visible:ring-2 focus-visible:ring-secondary-500 focus-visible:ring-offset-2"
                aria-describedby="trim-status"
              >
                {recording.isTrimming ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" aria-hidden="true"></div>
                    <span id="trim-status">Trimming...</span>
                  </>
                ) : (
                  <>
                    <Scissors size={16} aria-hidden="true" />
                    <span id="trim-status">Apply Trim & Re-analyze</span>
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
      
      {/* Alternative recording input for accessibility */}
      {settings.visualIndicators && (
        <AlternativeRecordingInput
          onRecordStart={startRecording}
          onRecordStop={stopRecording}
          isRecording={recording.isRecording}
        />
      )}
      
      {/* Screen reader announcements */}
      <LiveRegion>
        {recording.isRecording && 'Recording in progress'}
        {recording.isPlaying && 'Playing recording'}
        {recording.isTrimming && 'Trimming audio'}
      </LiveRegion>
    </div>
  );
}
