'use client';

import { useState, useRef } from 'react';
import { AudioRecorder } from '@/lib/audio/recorder';
import Waveform from './Waveform';
import { Mic, Square, Play } from 'lucide-react';

interface RecorderProps {
  onRecordingComplete: (audioBuffer: AudioBuffer) => void;
}

export default function Recorder({ onRecordingComplete }: RecorderProps) {
  const [isRecording, setIsRecording] = useState(false);
  const [hasPermission, setHasPermission] = useState(false);
  const [recordedAudio, setRecordedAudio] = useState<AudioBuffer | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  
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

  return (
    <div className="space-y-4">
      <Waveform
        isRecording={isRecording}
        getWaveformData={isRecording ? () => recorderRef.current!.getWaveformData() : undefined}
        audioBuffer={recordedAudio}
      />

      <div className="flex gap-4 justify-center">
        {!isRecording ? (
          <button
            onClick={startRecording}
            className="flex items-center gap-2 px-6 py-3 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors"
          >
            <Mic size={20} />
            Start Recording
          </button>
        ) : (
          <button
            onClick={stopRecording}
            className="flex items-center gap-2 px-6 py-3 bg-red-500 hover:bg-red-600 rounded-lg font-medium transition-colors"
          >
            <Square size={20} />
            Stop Recording
          </button>
        )}

        {recordedAudio && !isRecording && (
          <button
            onClick={playRecording}
            disabled={isPlaying}
            className="flex items-center gap-2 px-6 py-3 bg-secondary-500 hover:bg-secondary-600 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            <Play size={20} />
            Play Back
          </button>
        )}
      </div>

      {recordedAudio && (
        <div className="text-center text-sm text-gray-400">
          Duration: {recordedAudio.duration.toFixed(2)}s
        </div>
      )}
    </div>
  );
}

