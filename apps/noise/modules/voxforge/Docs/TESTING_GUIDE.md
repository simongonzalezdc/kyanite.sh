# VoxForge Testing Guide

This guide provides comprehensive information about the testing infrastructure implemented for VoxForge.

## Table of Contents

1. [Testing Configuration](#testing-configuration)
2. [Test Utilities](#test-utilities)
3. [Unit Tests](#unit-tests)
4. [Component Tests](#component-tests)
5. [Integration Tests](#integration-tests)
6. [E2E Tests](#e2e-tests)
7. [Performance Tests](#performance-tests)
8. [Test Coverage](#test-coverage)
9. [Pre-commit Hooks](#pre-commit-hooks)
10. [CI/CD Integration](#cicd-integration)

## Testing Configuration

### Jest Configuration

The project uses Jest as the primary testing framework with the following configuration:

- **TypeScript Support**: Full TypeScript support with ts-jest
- **React Testing Library**: For component testing
- **Custom Matchers**: Audio-specific matchers for testing
- **Test Environment**: jsdom for DOM testing
- **Coverage Reporting**: Istanbul integration with coverage thresholds

### Cypress Configuration

Cypress is configured for E2E testing with:

- **Base URL**: http://localhost:3000
- **Browser Support**: Chrome, Firefox, Safari
- **Mobile Testing**: Responsive viewport testing
- **Video Recording**: For debugging test failures
- **Screenshot Capture**: On test failures

## Test Utilities

### Audio Mocks

The project includes comprehensive mocks for Web Audio API:

```typescript
import { setupAudioMocks, createAudioBuffer } from '@/__tests__/utils/audioMocks'

// Setup mocks before tests
beforeAll(() => {
  setupAudioMocks()
})

// Create test audio buffer
const audioBuffer = createAudioBuffer(2, 44100)
```

### Test Data Factories

Factories for creating test data:

```typescript
import { 
  createPitchPoints, 
  createBPMAnalysis, 
  createKeyAnalysis, 
  createSection 
} from '@/__tests__/utils/testDataFactories'

// Create test pitch points
const pitches = createPitchPoints(5, 440)
```

### Test Helpers

Common testing utilities:

```typescript
import { 
  renderWithProviders, 
  simulateMediaPermission, 
  measurePerformance 
} from '@/__tests__/utils/testHelpers'

// Render with providers
const { container } = renderWithProviders(<Component />)

// Simulate media permission
simulateMediaPermission(true)

// Measure performance
const duration = await measurePerformance(async () => {
  return someAsyncFunction()
})
```

## Unit Tests

### Audio Processing Tests

Tests for audio processing functions:

- **Pitch Detection**: Tests for pitch detection algorithms
- **BPM Detection**: Tests for tempo detection
- **Key Detection**: Tests for musical key detection
- **Audio Analysis**: Tests for complete audio analysis pipeline

### Utility Function Tests

Tests for utility functions:

- **Audio Utils**: Tests for audio manipulation functions
- **Music Theory**: Tests for music theory calculations
- **Format Utils**: Tests for data formatting functions

### Store Tests

Tests for state management:

- **Store Actions**: Tests for all store actions
- **Store Selectors**: Tests for optimized selectors
- **State Persistence**: Tests for state persistence

## Component Tests

### React Component Testing

Components are tested with React Testing Library:

- **User Interactions**: Click, type, drag events
- **State Changes**: Component state updates
- **Responsive Behavior**: Mobile, tablet, desktop views
- **Accessibility**: ARIA labels, keyboard navigation

### Example Component Test

```typescript
import { render, screen, fireEvent } from '@testing-library/react'
import { renderWithProviders } from '@/__tests__/utils/testHelpers'

test('should handle user interactions', () => {
  renderWithProviders(<Recorder />)
  
  const recordButton = screen.getByTestId('record-button')
  fireEvent.click(recordButton)
  
  expect(screen.getByTestId('recording-indicator')).toBeInTheDocument()
})
```

## Integration Tests

### Audio Pipeline Tests

Tests for the complete audio processing pipeline:

- **Recording → Analysis**: End-to-end audio processing
- **Analysis → Generation**: Audio generation from analysis
- **State Management**: Integration between store modules
- **API Routes**: Server-side API integration

### Example Integration Test

```typescript
test('should process audio from recording to analysis', async () => {
  // Mock recording
  const audioBuffer = createAudioBuffer(2, 44100)
  
  // Process through pipeline
  const analysis = await processAudio(audioBuffer)
  
  // Verify results
  expect(analysis.pitches).toHaveLength(greaterThan(0))
  expect(analysis.bpm).toBeBPMAnalysis()
})
```

## E2E Tests

### User Workflow Tests

Complete user workflows are tested with Cypress:

- **Audio Recording**: Full recording workflow
- **Audio Analysis**: Analysis of recorded audio
- **Music Generation**: Generated music playback
- **Export Functionality**: Export in various formats
- **Responsive Design**: Mobile and desktop experiences

### Example E2E Test

```javascript
describe('Audio Recording Workflow', () => {
  it('should complete full recording workflow', () => {
    cy.visit('/')
    cy.mockMicrophonePermission(true)
    
    // Start recording
    cy.get('[data-testid="record-button"]').click()
    cy.get('[data-testid="recording-indicator"]').should('be.visible')
    
    // Stop recording
    cy.get('[data-testid="stop-button"]').click()
    cy.get('[data-testid="play-button"]').should('be.visible')
  })
})
```

## Performance Tests

### Audio Processing Performance

Performance tests ensure efficient audio processing:

- **Processing Time**: Time limits for audio analysis
- **Memory Usage**: Memory efficiency for large audio files
- **CPU Usage**: Reasonable CPU utilization
- **Scalability**: Performance with different audio sizes

### Example Performance Test

```typescript
test('should process audio within time limits', async () => {
  const audioBuffer = createAudioBuffer(5, 44100)
  
  const duration = await measurePerformance(async () => {
    return pitchDetector.analyze(audioBuffer)
  })
  
  expect(duration).toBeLessThan(500) // 500ms for 5 seconds
})
```

## Test Coverage

### Coverage Configuration

Test coverage is configured with:

- **Thresholds**: 70% for branches, functions, lines, statements
- **Reporting**: HTML and LCOV formats
- **Exclusions**: Test files, configuration files
- **Integration**: Coverage reporting in CI/CD

### Coverage Reports

Coverage reports are generated in:

- **HTML Report**: `coverage/lcov-report/index.html`
- **LCOV Report**: `coverage/lcov.info`
- **JSON Report**: `coverage/coverage-final.json`
- **Text Summary**: Console output

## Pre-commit Hooks

### Hook Configuration

Pre-commit hooks ensure code quality:

- **Linting**: ESLint with project configuration
- **Testing**: All tests must pass
- **Type Checking**: TypeScript compilation
- **Formatting**: Code formatting with Prettier

### Hook Implementation

```bash
#!/usr/bin/env sh
. "$(dirname "$0")/_/husky.sh"

# Run linting
npm run lint

# Run tests
npm run test

# Check if tests passed
if [ $? -ne 0 ]; then
  echo "Tests failed. Please fix them before committing."
  exit 1
fi

# Run type checking
npm run type-check

# Check if type checking passed
if [ $? -ne 0 ]; then
  echo "Type checking failed. Please fix type errors before committing."
  exit 1
fi

echo "Pre-commit checks passed!"
```

## CI/CD Integration

### GitHub Actions

The project uses GitHub Actions for CI/CD:

- **Multiple Test Types**: Unit, integration, E2E, performance
- **Parallel Execution**: Tests run in parallel for speed
- **Artifact Upload**: Test results and coverage reports
- **Test Summary**: Consolidated test status reporting

### Workflow Triggers

CI/CD is triggered by:

- **Push Events**: On push to main/develop branches
- **Pull Requests**: On PR creation/update
- **Scheduled Runs**: Daily for stability testing

### Test Environments

Tests run in:

- **Unit Tests**: Ubuntu with Node.js
- **E2E Tests**: Ubuntu with Chrome browser
- **Performance Tests**: Ubuntu with performance monitoring
- **Accessibility Tests**: Ubuntu with accessibility tools

## Running Tests

### Local Development

Run tests locally with:

```bash
# Run all tests
npm test

# Run with coverage
npm run test:coverage

# Run watch mode
npm run test:watch

# Run specific test types
npm run test:unit
npm run test:integration
npm run test:component
npm run test:e2e
npm run test:performance
```

### Debugging Tests

Debug failing tests with:

```bash
# Run tests in debug mode
node --inspect-brk node_modules/.bin/jest --runInBand

# Run specific test file
npm test -- --testPathPattern=specific.test.ts

# Run tests with verbose output
npm test -- --verbose
```

## Best Practices

### Test Organization

- **Descriptive Names**: Clear, descriptive test names
- **Test Isolation**: Tests should not depend on each other
- **Mocking**: Mock external dependencies appropriately
- **Assertions**: Comprehensive assertions for all outcomes

### Audio Testing

- **Mock Audio Context**: Use provided audio context mocks
- **Test Data**: Use realistic audio test data
- **Performance**: Consider performance in audio tests
- **Edge Cases**: Test silence, noise, clipping

### Component Testing

- **User Interactions**: Test all user interactions
- **State Changes**: Verify state updates
- **Accessibility**: Test with screen readers
- **Responsive**: Test different viewport sizes

### E2E Testing

- **User Workflows**: Test complete user journeys
- **Mobile Testing**: Test on mobile devices
- **Error Handling**: Test error states and recovery
- **Network Conditions**: Test slow/fast networks

## Troubleshooting

### Common Issues

1. **Audio Context Errors**: Ensure proper mocking
2. **Async Test Timeouts**: Use proper async/await
3. **Test Isolation**: Clean up between tests
4. **Mock Failures**: Verify mock implementations

### Debugging Tips

1. **Use Console Logs**: Add debugging logs to tests
2. **Check Mocks**: Verify mock implementations
3. **Test Data**: Validate test data creation
4. **Performance**: Profile slow tests

## Future Enhancements

### Planned Improvements

1. **Visual Regression**: Add visual regression testing
2. **API Testing**: Expand API route testing
3. **Load Testing**: Add performance load testing
4. **Mobile Testing**: Expand mobile device testing
5. **Accessibility**: Enhance accessibility testing

### Testing Tools

Consider adding:

1. **Storybook**: For component documentation
2. **MSW**: For API mocking
3. **Playwright**: For cross-browser E2E testing
4. **Axe**: For automated accessibility testing
5. **Lighthouse**: For performance auditing