describe('Audio Recording Workflow', () => {
  beforeEach(() => {
    cy.visit('/')
    cy.mockMicrophonePermission(true)
  })

  it('should complete full recording workflow', () => {
    // Start recording
    cy.get('[data-testid="record-button"]').click()
    
    // Check recording state
    cy.get('[data-testid="recording-indicator"]').should('be.visible')
    cy.get('[data-testid="recording-indicator"]').should('contain', 'Recording...')
    
    // Wait for recording to simulate some audio
    cy.wait(2000)
    
    // Stop recording
    cy.get('[data-testid="stop-button"]').click()
    
    // Check recorded state
    cy.get('[data-testid="play-button"]').should('be.visible')
    cy.get('[data-testid="audio-duration"]').should('be.visible')
    cy.get('[data-testid="waveform-visualization"]').should('be.visible')
  })

  it('should handle recording permission denied', () => {
    cy.mockMicrophonePermission(false)
    cy.visit('/')
    
    // Should show permission denied message
    cy.get('[data-testid="permission-denied"]').should('be.visible')
    cy.get('[data-testid="permission-denied"]').should('contain', 'Microphone access denied')
  })

  it('should play recorded audio', () => {
    // Record some audio first
    cy.get('[data-testid="record-button"]').click()
    cy.wait(2000)
    cy.get('[data-testid="stop-button"]').click()
    
    // Play the recorded audio
    cy.get('[data-testid="play-button"]').click()
    
    // Check playing state
    cy.get('[data-testid="playing-indicator"]').should('be.visible')
    cy.get('[data-testid="playing-indicator"]').should('contain', 'Playing...')
    cy.get('[data-testid="pause-button"]').should('be.visible')
  })

  it('should pause and resume playback', () => {
    // Record and play audio first
    cy.get('[data-testid="record-button"]').click()
    cy.wait(2000)
    cy.get('[data-testid="stop-button"]').click()
    cy.get('[data-testid="play-button"]').click()
    
    // Pause playback
    cy.get('[data-testid="pause-button"]').click()
    
    // Check paused state
    cy.get('[data-testid="play-button"]').should('be.visible')
    cy.get('[data-testid="play-button"]').should('contain', 'Resume')
    
    // Resume playback
    cy.get('[data-testid="play-button"]').click()
    
    // Check playing state again
    cy.get('[data-testid="pause-button"]').should('be.visible')
    cy.get('[data-testid="playing-indicator"]').should('be.visible')
  })

  it('should clear recorded audio', () => {
    // Record some audio first
    cy.get('[data-testid="record-button"]').click()
    cy.wait(2000)
    cy.get('[data-testid="stop-button"]').click()
    
    // Clear the recording
    cy.get('[data-testid="clear-button"]').click()
    
    // Check that recording is cleared
    cy.get('[data-testid="record-button"]').should('be.visible')
    cy.get('[data-testid="play-button"]').should('not.exist')
    cy.get('[data-testid="audio-duration"]').should('not.exist')
    cy.get('[data-testid="waveform-visualization"]').should('not.exist')
  })

  it('should show audio level meter during recording', () => {
    cy.get('[data-testid="record-button"]').click()
    
    // Check that audio level meter is visible
    cy.get('[data-testid="audio-level-meter"]').should('be.visible')
    
    // Stop recording
    cy.get('[data-testid="stop-button"]').click()
    
    // Audio level meter should still be visible for playback
    cy.get('[data-testid="audio-level-meter"]').should('be.visible')
  })

  it('should show recording timer', () => {
    cy.get('[data-testid="record-button"]').click()
    
    // Check that timer is visible and counting
    cy.get('[data-testid="recording-timer"]').should('be.visible')
    
    // Wait a bit and check timer updates
    cy.wait(1000)
    cy.get('[data-testid="recording-timer"]').should('not.contain', '00:00')
    
    // Stop recording
    cy.get('[data-testid="stop-button"]').click()
  })

  it('should handle keyboard shortcuts', () => {
    // Test spacebar for recording
    cy.get('body').type(' ')
    
    // Should start recording
    cy.get('[data-testid="recording-indicator"]').should('be.visible')
    
    // Test spacebar again to stop
    cy.get('body').type(' ')
    
    // Should stop recording
    cy.get('[data-testid="play-button"]').should('be.visible')
  })

  it('should be accessible', () => {
    // Check ARIA labels
    cy.get('[aria-label="Record audio"]').should('exist')
    cy.get('[aria-label="Stop recording"]').should('exist')
    cy.get('[aria-label="Play audio"]').should('exist')
    
    // Check roles
    cy.get('[role="button"][aria-label="Record audio"]').should('exist')
    cy.get('[role="button"][aria-label="Stop recording"]').should('exist')
    cy.get('[role="button"][aria-label="Play audio"]').should('exist')
  })

  it('should handle trim controls', () => {
    // Record some audio first
    cy.get('[data-testid="record-button"]').click()
    cy.wait(2000)
    cy.get('[data-testid="stop-button"]').click()
    
    // Check trim controls are visible
    cy.get('[data-testid="trim-controls"]').should('be.visible')
    cy.get('[data-testid="trim-start"]').should('be.visible')
    cy.get('[data-testid="trim-end"]').should('be.visible')
    
    // Test trim start
    cy.get('[data-testid="trim-start"]').type('0.5')
    
    // Test trim end
    cy.get('[data-testid="trim-end"]').type('1.5')
    
    // Apply trim
    cy.get('[data-testid="apply-trim"]').click()
    
    // Check that duration is updated
    cy.get('[data-testid="audio-duration"]').should('contain', '01:00')
  })

  it('should handle export functionality', () => {
    // Record some audio first
    cy.get('[data-testid="record-button"]').click()
    cy.wait(2000)
    cy.get('[data-testid="stop-button"]').click()
    
    // Open export panel
    cy.get('[data-testid="export-button"]').click()
    
    // Check export options
    cy.get('[data-testid="export-panel"]').should('be.visible')
    cy.get('[data-testid="export-format"]').should('be.visible')
    cy.get('[data-testid="export-quality"]').should('be.visible')
    
    // Select export format
    cy.get('[data-testid="export-format"]').select('wav')
    
    // Export the audio
    cy.get('[data-testid="confirm-export"]').click()
    
    // Check success message
    cy.get('[data-testid="export-success"]').should('be.visible')
    cy.get('[data-testid="export-success"]').should('contain', 'Audio exported successfully')
  })

  it('should handle mobile view', () => {
    // Test mobile viewport
    cy.viewport(375, 667) // iPhone SE
    
    cy.visit('/')
    cy.mockMicrophonePermission(true)
    
    // Check mobile-specific elements
    cy.get('[data-testid="mobile-navigation"]').should('be.visible')
    cy.get('[data-testid="record-button"]').should('be.visible')
    
    // Test mobile touch interactions
    cy.get('[data-testid="record-button"]').tap()
    
    // Should start recording
    cy.get('[data-testid="recording-indicator"]').should('be.visible')
  })

  it('should handle tablet view', () => {
    // Test tablet viewport
    cy.viewport(768, 1024) // iPad
    
    cy.visit('/')
    cy.mockMicrophonePermission(true)
    
    // Check tablet layout
    cy.get('[data-testid="tablet-layout"]').should('be.visible')
    cy.get('[data-testid="record-button"]').should('be.visible')
  })

  it('should handle desktop view', () => {
    // Test desktop viewport
    cy.viewport(1280, 720) // Desktop
    
    cy.visit('/')
    cy.mockMicrophonePermission(true)
    
    // Check desktop layout
    cy.get('[data-testid="desktop-layout"]').should('be.visible')
    cy.get('[data-testid="record-button"]').should('be.visible')
  })

  it('should handle error states', () => {
    // Mock recording error
    cy.window().then((win) => {
      win.navigator.mediaDevices.getUserMedia = () => Promise.reject(
        new Error('Recording device not found')
      )
    })
    
    cy.visit('/')
    
    // Check error message
    cy.get('[data-testid="error-message"]').should('be.visible')
    cy.get('[data-testid="error-message"]').should('contain', 'Recording device not found')
  })

  it('should handle network conditions', () => {
    // Test slow network
    cy.window().then((win) => {
      Object.defineProperty(win.navigator, 'connection', {
        value: {
          effectiveType: 'slow-2g',
          downlink: 0.1,
          rtt: 2000,
        },
        writable: true,
      })
    })
    
    cy.visit('/')
    cy.mockMicrophonePermission(true)
    
    // Should show network warning
    cy.get('[data-testid="network-warning"]').should('be.visible')
    cy.get('[data-testid="network-warning"]').should('contain', 'Slow network detected')
  })
})