'use client';

import { useState, useRef, useEffect } from 'react';
import { AudioRecorder } from '@/lib/audio/recorder';
import Waveform from './Waveform';
import Tooltip from './Tooltip';
import { Mic, Square, Play, Scissors, RotateCcw } from 'lucide-react';
import { trimAudioBuffer } from '@/lib/utils/audio-utils';
import { isMobile, isTouchDevice, triggerHaptic, enableSwipeGestures } from '@/lib/utils';

interface RecorderProps {
  onRecordingComplete: (audioBuffer: AudioBuffer) => void;
}

export default function Recorder({ onRecordingComplete }: RecorderProps) {
  const [isRecording, setIsRecording] = useState(false);
  const [hasPermission, setHasPermission] = useState(false);
  const [recordedAudio, setRecordedAudio] = useState<AudioBuffer | null>(null);
  const [originalAudio, setOriginalAudio] = useState<AudioBuffer | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [trimStart, setTrimStart] = useState(0);
  const [trimEnd, setTrimEnd] = useState(0);
  const [isTrimming, setIsTrimming] = useState(false);
  
  const recorderRef = useRef<AudioRecorder | null>(null);
  const audioContextRef = useRef<AudioContext | null>(null);

  const requestPermission = async () => {
    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }
    
    const granted = await recorderRef.current.requestPermission();
    setHasPermission(granted);
    
    if (!granted) {
      alert('Microphone permission is required to record audio.');
    }
  };

  const startRecording = async () => {
    if (!hasPermission) {
      await requestPermission();
      if (!hasPermission) return;
    }

    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }

    try {
      await recorderRef.current.startRecording();
      setIsRecording(true);
      setRecordedAudio(null);
    } catch (error) {
      console.error('Failed to start recording:', error);
      alert('Failed to start recording. Please check microphone permissions.');
    }
  };

  const stopRecording = async () => {
    if (!recorderRef.current || !isRecording) return;

    try {
      const audioBuffer = await recorderRef.current.stopRecording();
      setIsRecording(false);
      setRecordedAudio(audioBuffer);
      setOriginalAudio(audioBuffer);
      setTrimStart(0);
      setTrimEnd(audioBuffer.duration);
      setIsTrimming(false);
      onRecordingComplete(audioBuffer);
    } catch (error) {
      console.error('Failed to stop recording:', error);
      alert('Failed to process recording.');
    }
  };

  const playRecording = async () => {
    if (!recordedAudio) return;

    if (!audioContextRef.current) {
      audioContextRef.current = new AudioContext();
    }

    const source = audioContextRef.current.createBufferSource();
    source.buffer = recordedAudio;
    source.connect(audioContextRef.current.destination);
    
    source.onended = () => {
      setIsPlaying(false);
    };

    source.start();
    setIsPlaying(true);
  };

  const handleTrim = async () => {
    // Always trim from originalAudio to allow re-trimming with different values
    const sourceAudio = originalAudio || recordedAudio;
    if (!sourceAudio) return;
    
    setIsTrimming(true);
    try {
      const trimmed = await trimAudioBuffer(sourceAudio, trimStart, trimEnd);
      setRecordedAudio(trimmed);
      // Ensure originalAudio is set for future resets
      if (!originalAudio) {
        setOriginalAudio(sourceAudio);
      }
      onRecordingComplete(trimmed);
    } catch (error) {
      console.error('Failed to trim audio:', error);
      alert('Failed to trim audio. Please check your trim values.');
    } finally {
      setIsTrimming(false);
    }
  };

  const handleResetTrim = () => {
    if (!originalAudio) return;
    setRecordedAudio(originalAudio);
    setTrimStart(0);
    setTrimEnd(originalAudio.duration);
    onRecordingComplete(originalAudio);
  };

  // Update trim end when audio changes
  useEffect(() => {
    if (recordedAudio && !isRecording) {
      setTrimEnd(recordedAudio.duration);
    }
  }, [recordedAudio, isRecording]);

  // Add swipe gesture support for mobile
  const waveformRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    if (isTouchDevice() && waveformRef.current) {
      const cleanup = enableSwipeGestures(waveformRef.current, {
        onSwipeLeft: () => {
          if (recordedAudio && !isRecording && !isPlaying) {
            playRecording();
          }
        },
        onSwipeRight: () => {
          if (isRecording) {
            stopRecording();
          } else if (!recordedAudio) {
            startRecording();
          }
        }
      });
      
      return cleanup;
    }
  }, [isRecording, recordedAudio, isPlaying]);

  return (
    <div className="space-y-4" data-tour="recorder" id="record">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 mb-2">
        <h2 className="fluid-xl font-semibold flex items-center gap-2">
          Voice Recorder
          <Tooltip content="Record your voice, melody, or any sound. The recorder will capture audio from your microphone and analyze it for pitch and timing." />
        </h2>
      </div>
      
      <div ref={waveformRef} className="touch-manipulation">
        <Waveform
          isRecording={isRecording}
          getWaveformData={isRecording ? () => recorderRef.current!.getWaveformData() : undefined}
          audioBuffer={recordedAudio}
        />
      </div>

      <div className="flex flex-col sm:flex-row gap-4 justify-center">
        {!isRecording ? (
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
            `}
            style={{ minHeight: isMobile() ? '60px' : '48px' }}
          >
            <Mic size={isMobile() ? 24 : 20} />
            <span>Start Recording</span>
          </button>
        ) : (
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
            `}
            style={{ minHeight: isMobile() ? '60px' : '48px' }}
          >
            <Square size={isMobile() ? 24 : 20} />
            <span>Stop Recording</span>
          </button>
        )}

        {recordedAudio && !isRecording && (
          <button
            onClick={() => {
              playRecording();
              triggerHaptic('light');
            }}
            disabled={isPlaying}
            className={`
              flex items-center justify-center gap-2 px-6 py-4 sm:py-3
              bg-secondary-500 hover:bg-secondary-600
              rounded-lg font-medium transition-all duration-200
              disabled:opacity-50 disabled:cursor-not-allowed
              min-h-[44px] min-w-[44px] sm:min-w-0
              ${isMobile() ? 'text-lg' : ''}
              active:scale-95 touch-manipulation
            `}
            style={{ minHeight: isMobile() ? '48px' : '44px' }}
          >
            <Play size={isMobile() ? 24 : 20} />
            <span>Play Back</span>
          </button>
        )}
      </div>

      {recordedAudio && !isRecording && (
        <div className="space-y-4">
          <div className="text-center fluid-sm text-gray-400">
            Duration: {recordedAudio.duration.toFixed(2)}s
            {originalAudio && originalAudio !== recordedAudio && (
              <span className="ml-2 text-primary-500 block sm:inline">
                (Trimmed from {originalAudio.duration.toFixed(2)}s)
              </span>
            )}
          </div>

          {/* Trim Controls */}
          <div className="bg-gray-800 rounded-lg p-4 sm:p-6 border border-gray-700 space-y-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
              <h3 className="fluid-sm font-medium flex items-center gap-2">
                <Scissors size={16} className="text-secondary-500" />
                Trim Audio
                <Tooltip content="Remove unwanted parts from your recording. Adjust the start and end points to keep only the best portion of your audio." />
              </h3>
              {originalAudio && originalAudio !== recordedAudio && (
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
                <label className="block fluid-xs text-gray-400 mb-2">
                  Start: {trimStart.toFixed(2)}s
                </label>
                <input
                  type="range"
                  min="0"
                  max={originalAudio?.duration || recordedAudio.duration}
                  step="0.01"
                  value={trimStart}
                  onChange={(e) => {
                    const newStart = parseFloat(e.target.value);
                    setTrimStart(Math.min(newStart, trimEnd - 0.1));
                  }}
                  className="w-full h-3 sm:h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-primary-500 touch-manipulation"
                  style={{ minHeight: isMobile() ? '24px' : '16px' }}
                />
              </div>

              <div>
                <label className="block fluid-xs text-gray-400 mb-2">
                  End: {trimEnd.toFixed(2)}s
                </label>
                <input
                  type="range"
                  min={trimStart + 0.1}
                  max={originalAudio?.duration || recordedAudio.duration}
                  step="0.01"
                  value={trimEnd}
                  onChange={(e) => {
                    const newEnd = parseFloat(e.target.value);
                    setTrimEnd(Math.max(newEnd, trimStart + 0.1));
                  }}
                  className="w-full h-3 sm:h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer accent-primary-500 touch-manipulation"
                  style={{ minHeight: isMobile() ? '24px' : '16px' }}
                />
              </div>

              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between fluid-xs text-gray-400 gap-2">
                <span>Trimmed length: {(trimEnd - trimStart).toFixed(2)}s</span>
                <span>
                  {originalAudio && (
                    <>Remove: {(originalAudio.duration - (trimEnd - trimStart)).toFixed(2)}s</>
                  )}
                </span>
              </div>

              <button
                onClick={() => {
                  handleTrim();
                  triggerHaptic('medium');
                }}
                disabled={isTrimming || !recordedAudio || (trimStart === 0 && trimEnd === (originalAudio?.duration || recordedAudio.duration))}
                className="w-full px-4 py-3 sm:py-2 bg-secondary-500 hover:bg-secondary-600 rounded-lg font-medium transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 fluid-sm min-h-[44px] touch-manipulation active:scale-95"
              >
                {isTrimming ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    Trimming...
                  </>
                ) : (
                  <>
                    <Scissors size={16} />
                    Apply Trim & Re-analyze
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

