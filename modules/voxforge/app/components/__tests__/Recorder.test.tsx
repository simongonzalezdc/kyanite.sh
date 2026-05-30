import React from 'react'
import { fireEvent, screen, waitFor } from '@testing-library/react'
import { Recorder } from '../Recorder'
import { renderWithProviders, simulateMediaPermission } from '@/__tests__/utils/testHelpers'
import { useAudioStore } from '@/lib/store/audioStore'

const renderRecorder = () => renderWithProviders(<Recorder />)

const recordOneClip = async () => {
  renderRecorder()

  fireEvent.click(screen.getByTestId('record-button'))

  const stopButton = await screen.findByTestId('stop-button')
  fireEvent.click(stopButton)

  await screen.findByTestId('play-button')
}

describe('Recorder Component', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    simulateMediaPermission(true)
    useAudioStore.getState().resetAll()
  })

  it('renders the recorder controls and waveform surface', () => {
    renderRecorder()

    expect(screen.getByTestId('recorder-container')).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: /voice recorder/i })).toBeInTheDocument()
    expect(screen.getByTestId('record-button')).toHaveTextContent('Start Recording')
    expect(screen.getByTestId('waveform-visualization')).toBeInTheDocument()
  })

  it('starts recording from a single click after microphone permission is granted', async () => {
    renderRecorder()

    fireEvent.click(screen.getByTestId('record-button'))

    const stopButton = await screen.findByTestId('stop-button')
    expect(stopButton).toHaveTextContent('Stop Recording')
    expect(screen.getByRole('status')).toHaveTextContent('Recording in progress')
  })

  it('stores a recorded clip and reveals playback and trim controls', async () => {
    await recordOneClip()

    expect(screen.getByTestId('play-button')).toHaveTextContent('Play Back')
    expect(screen.getByTestId('audio-duration')).toHaveTextContent('Duration:')
    expect(screen.getByTestId('trim-controls')).toBeInTheDocument()
    expect(screen.getByTestId('trim-start')).toBeInTheDocument()
    expect(screen.getByTestId('trim-end')).toBeInTheDocument()
  })

  it('plays a recorded clip and announces playback state', async () => {
    await recordOneClip()

    fireEvent.click(screen.getByTestId('play-button'))

    await waitFor(() => {
      expect(screen.getByTestId('play-button')).toBeDisabled()
      expect(screen.getByText('Playing recording')).toBeInTheDocument()
    })
  })

  it('keeps the recorder stable when microphone permission is denied', async () => {
    simulateMediaPermission(false)
    renderRecorder()

    fireEvent.click(screen.getByTestId('record-button'))

    await waitFor(() => {
      expect(window.alert).toHaveBeenCalledWith('Microphone permission is required to record audio.')
    })
    expect(screen.queryByTestId('stop-button')).not.toBeInTheDocument()
    expect(screen.getByTestId('record-button')).toBeInTheDocument()
  })

  it('exposes accessible labels for record, stop, and playback states', async () => {
    renderRecorder()

    expect(screen.getByRole('button', { name: /start recording/i })).toBeInTheDocument()

    fireEvent.click(screen.getByTestId('record-button'))
    expect(await screen.findByRole('button', { name: /stop recording/i })).toBeInTheDocument()

    fireEvent.click(screen.getByTestId('stop-button'))
    expect(await screen.findByRole('button', { name: /play recording/i })).toBeInTheDocument()
  })
})
