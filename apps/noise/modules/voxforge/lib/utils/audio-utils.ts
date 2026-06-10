export function frequencyToMidi(frequency: number): number {
  return Math.round(69 + 12 * Math.log2(frequency / 440));
}

export function midiToFrequency(midi: number): number {
  return 440 * Math.pow(2, (midi - 69) / 12);
}

export function downloadBlob(blob: Blob, filename: string): void {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

export function formatDuration(seconds: number): string {
  const safeSeconds = Math.max(0, seconds);
  const mins = Math.floor(safeSeconds / 60);
  const secs = Math.floor(safeSeconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

/**
 * Trim an AudioBuffer by extracting a portion between start and end times
 * @param audioBuffer The source AudioBuffer
 * @param startTime Start time in seconds (default: 0)
 * @param endTime End time in seconds (default: audioBuffer.duration)
 * @returns A new trimmed AudioBuffer
 */
export async function trimAudioBuffer(
  audioBuffer: AudioBuffer,
  startTime: number,
  endTime: number
): Promise<AudioBuffer> {
  const sampleRate = audioBuffer.sampleRate;
  const numberOfChannels = audioBuffer.numberOfChannels;
  
  // Clamp values to valid range
  const start = Math.max(0, Math.min(startTime, audioBuffer.duration));
  const end = Math.max(start, Math.min(endTime, audioBuffer.duration));
  
  // Calculate sample indices
  const startSample = Math.floor(start * sampleRate);
  const endSample = Math.floor(end * sampleRate);
  const length = endSample - startSample;
  
  if (length <= 0) {
    throw new Error('Invalid trim range: end time must be greater than start time');
  }
  
  // Create new AudioBuffer with trimmed length
  const trimmedBuffer = new AudioContext().createBuffer(
    numberOfChannels,
    length,
    sampleRate
  );
  
  // Copy channel data for each channel
  for (let channel = 0; channel < numberOfChannels; channel++) {
    const sourceData = audioBuffer.getChannelData(channel);
    const targetData = trimmedBuffer.getChannelData(channel);
    const trimmedData = sourceData.subarray(startSample, endSample);
    targetData.set(trimmedData);
  }
  
  return trimmedBuffer;
}
