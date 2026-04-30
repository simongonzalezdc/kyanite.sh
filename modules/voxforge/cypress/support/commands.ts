// Custom commands for VoxForge testing

declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Mock audio recording with specified duration
       * @param duration - Duration in seconds
       */
      mockAudioRecording(duration: number): Chainable<any>

      /**
       * Mock audio analysis results
       * @param analysis - Mock analysis data
       */
      mockAudioAnalysis(analysis: any): Chainable<any>

      /**
       * Wait for audio processing to complete
       */
      waitForAudioProcessing(): Chainable<any>

      /**
       * Check if audio is playing
       */
      checkAudioPlaying(): Chainable<any>

      /**
       * Mock microphone permission
       * @param granted - Whether permission is granted
       */
      mockMicrophonePermission(granted: boolean): Chainable<any>
    }
  }
}

// Mock audio recording command
Cypress.Commands.add('mockAudioRecording', (duration: number) => {
  cy.window().then((win: any) => {
    // Create a mock audio buffer with the specified duration
    const sampleRate = 44100
    const length = sampleRate * duration
    const audioBuffer = {
      numberOfChannels: 2,
      length,
      sampleRate,
      duration,
      getChannelData: () => new Float32Array(length),
    }

    // Mock the recorder to return this buffer
    win.mockAudioBuffer = audioBuffer
  })
})

// Mock audio analysis command
Cypress.Commands.add('mockAudioAnalysis', (analysis: any) => {
  cy.window().then((win: any) => {
    win.mockAudioAnalysis = analysis
  })
})

// Wait for audio processing command
Cypress.Commands.add('waitForAudioProcessing', () => {
  cy.get('[data-testid="audio-processing"]', { timeout: 10000 }).should('not.exist')
  cy.get('[data-testid="audio-analysis-complete"]', { timeout: 10000 }).should('exist')
})

// Check if audio is playing command
Cypress.Commands.add('checkAudioPlaying', () => {
  cy.get('[data-testid="audio-playing-indicator"]').should('exist')
  cy.get('[data-testid="playback-controls"]').should('contain', 'Pause')
})

// Mock microphone permission command
Cypress.Commands.add('mockMicrophonePermission', (granted: boolean) => {
  cy.window().then((win: any) => {
    if (granted) {
      win.navigator.mediaDevices.getUserMedia = () => Promise.resolve({
        getTracks: () => [{
          stop: () => {},
          getSettings: () => ({
            deviceId: 'default',
            groupId: 'default',
            kind: 'audioinput',
            label: 'Default',
            sampleRate: 44100,
          }),
        }],
      })
    } else {
      win.navigator.mediaDevices.getUserMedia = () => Promise.reject(
        new Error('Permission denied')
      )
    }
  })
})

export {}