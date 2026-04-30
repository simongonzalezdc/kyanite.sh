import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { Recorder } from '../Recorder'
import { renderWithProviders, simulateMediaPermission, createAudioBuffer } from '@/__tests__/utils/testHelpers'

// Mock Web Audio API
beforeAll(() => {
  // Mock getUserMedia to grant permission
  simulateMediaPermission(true)
})

describe('Recorder Component', () => {
  it('should render recorder component', () => {
    renderWithProviders(<Recorder />)
    
    expect(screen.getByTestId('recorder-container')).toBeInTheDocument()
    expect(screen.getByTestId('record-button')).toBeInTheDocument()
    expect(screen.getByTestId('stop-button')).toBeInTheDocument()
    expect(screen.getByTestId('play-button')).toBeInTheDocument()
  })

  it('should show record button initially', () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    expect(recordButton).toBeInTheDocument()
    expect(recordButton).toHaveTextContent('Record')
  })

  it('should show stop button when recording', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      expect(screen.getByTestId('stop-button')).toBeInTheDocument()
      expect(screen.getByTestId('stop-button')).toHaveTextContent('Stop')
    })
  })

  it('should show play button when audio is recorded', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('play-button')).toBeInTheDocument()
      expect(screen.getByTestId('play-button')).toHaveTextContent('Play')
    })
  })

  it('should show pause button when playing', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording, then play
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      const playButton = screen.getByTestId('play-button')
      fireEvent.click(playButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('pause-button')).toBeInTheDocument()
      expect(screen.getByTestId('pause-button')).toHaveTextContent('Pause')
    })
  })

  it('should show recording indicator when recording', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      expect(screen.getByTestId('recording-indicator')).toBeInTheDocument()
      expect(screen.getByTestId('recording-indicator')).toHaveTextContent('Recording...')
    })
  })

  it('should show playing indicator when playing', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording, then play
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      const playButton = screen.getByTestId('play-button')
      fireEvent.click(playButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('playing-indicator')).toBeInTheDocument()
      expect(screen.getByTestId('playing-indicator')).toHaveTextContent('Playing...')
    })
  })

  it('should show audio level meter', () => {
    renderWithProviders(<Recorder />)
    
    expect(screen.getByTestId('audio-level-meter')).toBeInTheDocument()
  })

  it('should show timer when recording', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      expect(screen.getByTestId('recording-timer')).toBeInTheDocument()
    })
  })

  it('should show audio duration when recorded', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('audio-duration')).toBeInTheDocument()
    })
  })

  it('should handle microphone permission denied', async () => {
    // Mock denied permission
    simulateMediaPermission(false)
    
    renderWithProviders(<Recorder />)
    
    await waitFor(() => {
      expect(screen.getByTestId('permission-denied')).toBeInTheDocument()
      expect(screen.getByTestId('permission-denied')).toHaveTextContent('Microphone access denied')
    })
  })

  it('should show settings button', () => {
    renderWithProviders(<Recorder />)
    
    expect(screen.getByTestId('settings-button')).toBeInTheDocument()
  })

  it('should show clear button when audio is recorded', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('clear-button')).toBeInTheDocument()
      expect(screen.getByTestId('clear-button')).toHaveTextContent('Clear')
    })
  })

  it('should handle keyboard shortcuts', () => {
    renderWithProviders(<Recorder />)
    
    // Test spacebar for record/stop
    fireEvent.keyDown(document, { key: ' ' })
    
    // Should trigger record button
    expect(screen.getByTestId('record-button')).toHaveFocus()
  })

  it('should be accessible', () => {
    const { container } = renderWithProviders(<Recorder />)
    
    // Check for proper ARIA labels
    expect(screen.getByLabelText('Record audio')).toBeInTheDocument()
    expect(screen.getByLabelText('Stop recording')).toBeInTheDocument()
    expect(screen.getByLabelText('Play audio')).toBeInTheDocument()
    
    // Check for proper roles
    expect(screen.getByRole('button', { name: 'Record audio' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Stop recording' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Play audio' })).toBeInTheDocument()
  })

  it('should handle responsive design', () => {
    const { container } = renderWithProviders(<Recorder />)
    
    // Should be responsive
    expect(container.firstChild).toHaveClass('responsive')
  })

  it('should show loading state during processing', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      expect(screen.getByTestId('processing-indicator')).toBeInTheDocument()
      expect(screen.getByTestId('processing-indicator')).toHaveTextContent('Processing...')
    })
  })

  it('should show error message on recording error', async () => {
    renderWithProviders(<Recorder />)
    
    // Mock recording error
    const recordButton = screen.getByTestId('record-button')
    fireEvent.error(recordButton)
    
    await waitFor(() => {
      expect(screen.getByTestId('error-message')).toBeInTheDocument()
      expect(screen.getByTestId('error-message')).toHaveTextContent('Recording failed')
    })
  })

  it('should show waveform visualization', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('waveform-visualization')).toBeInTheDocument()
    })
  })

  it('should handle audio export', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('export-button')).toBeInTheDocument()
      expect(screen.getByTestId('export-button')).toHaveTextContent('Export')
    })
  })

  it('should handle audio trim controls', async () => {
    renderWithProviders(<Recorder />)
    
    const recordButton = screen.getByTestId('record-button')
    
    // Start and stop recording
    fireEvent.click(recordButton)
    
    await waitFor(() => {
      const stopButton = screen.getByTestId('stop-button')
      fireEvent.click(stopButton)
    })
    
    await waitFor(() => {
      expect(screen.getByTestId('trim-controls')).toBeInTheDocument()
      expect(screen.getByTestId('trim-start')).toBeInTheDocument()
      expect(screen.getByTestId('trim-end')).toBeInTheDocument()
    })
  })
})