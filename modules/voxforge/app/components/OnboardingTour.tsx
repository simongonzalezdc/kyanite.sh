'use client';

import { useState, useEffect, useRef } from 'react';
import { X, ArrowRight, ArrowLeft, SkipForward } from 'lucide-react';

interface TourStep {
  id: string;
  title: string;
  content: string;
  target: string; // CSS selector for the element to highlight
  position: 'top' | 'bottom' | 'left' | 'right' | 'center';
}

interface OnboardingTourProps {
  isOpen: boolean;
  onComplete: () => void;
  onSkip: () => void;
}

const tourSteps: TourStep[] = [
  {
    id: 'welcome',
    title: 'Welcome to VoxForge!',
    content: 'Let\'s take a quick tour to help you get started with transforming your voice into music.',
    target: 'body',
    position: 'center'
  },
  {
    id: 'recorder',
    title: 'Record Your Voice',
    content: 'Start by recording your voice using the recorder. Click the microphone button to begin recording your melody or vocals.',
    target: '[data-tour="recorder"]',
    position: 'bottom'
  },
  {
    id: 'analysis',
    title: 'Audio Analysis',
    content: 'VoxForge automatically analyzes your recording to detect pitch, tempo, key, and time signature. You can edit these values if needed.',
    target: '[data-tour="analysis"]',
    position: 'top'
  },
  {
    id: 'instruments',
    title: 'Choose Instruments',
    content: 'Select which instruments you want to generate. You can choose from drums, bass, chords, and more.',
    target: '[data-tour="instruments"]',
    position: 'left'
  },
  {
    id: 'generate',
    title: 'Generate Music',
    content: 'Click the Generate Music button to create accompaniment based on your voice recording and selected instruments.',
    target: '[data-tour="generate"]',
    position: 'top'
  },
  {
    id: 'export',
    title: 'Export Your Creation',
    content: 'Export your music as individual stems or MIDI files. You can download each instrument separately or the full mix.',
    target: '[data-tour="export"]',
    position: 'left'
  }
];

export default function OnboardingTour({ isOpen, onComplete, onSkip }: OnboardingTourProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const [highlightedElement, setHighlightedElement] = useState<Element | null>(null);
  const [tooltipPosition, setTooltipPosition] = useState({ top: 0, left: 0 });
  const overlayRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isOpen) {
      cleanup();
      return;
    }

    const step = tourSteps[currentStep];
    if (step.target === 'body') {
      setHighlightedElement(null);
      setTooltipPosition({
        top: window.innerHeight / 2 - 150,
        left: window.innerWidth / 2 - 200
      });
      return;
    }

    const element = document.querySelector(step.target);
    if (element) {
      setHighlightedElement(element);
      updateTooltipPosition(element, step.position);
      
      // Scroll element into view if needed
      element.scrollIntoView({
        behavior: 'smooth',
        block: 'center',
        inline: 'center'
      });
    } else {
      // If element not found, skip to next step
      if (currentStep < tourSteps.length - 1) {
        setCurrentStep(currentStep + 1);
      } else {
        onComplete();
      }
    }

    return cleanup;
  }, [isOpen, currentStep]); // Remove tourSteps.length since it's constant

  const cleanup = () => {
    setHighlightedElement(null);
  };

  const updateTooltipPosition = (element: Element, position: string) => {
    const rect = element.getBoundingClientRect();
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    const scrollLeft = window.pageXOffset || document.documentElement.scrollLeft;
    
    let top = 0;
    let left = 0;

    switch (position) {
      case 'top':
        top = rect.top + scrollTop - 200;
        left = rect.left + scrollLeft + rect.width / 2 - 200;
        break;
      case 'bottom':
        top = rect.bottom + scrollTop + 20;
        left = rect.left + scrollLeft + rect.width / 2 - 200;
        break;
      case 'left':
        top = rect.top + scrollTop + rect.height / 2 - 75;
        left = rect.left + scrollLeft - 420;
        break;
      case 'right':
        top = rect.top + scrollTop + rect.height / 2 - 75;
        left = rect.right + scrollLeft + 20;
        break;
      case 'center':
        top = rect.top + scrollTop + rect.height / 2 - 150;
        left = rect.left + scrollLeft + rect.width / 2 - 200;
        break;
    }

    // Ensure tooltip stays within viewport
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;
    
    if (left < 20) left = 20;
    if (left + 400 > viewportWidth) left = viewportWidth - 420;
    if (top < 20) top = 20;
    if (top + 300 > viewportHeight) top = viewportHeight - 320;

    setTooltipPosition({ top, left });
  };

  const handleNext = () => {
    if (currentStep < tourSteps.length - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      onComplete();
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const handleSkip = () => {
    onSkip();
  };

  if (!isOpen) return null;

  const step = tourSteps[currentStep];
  const isLastStep = currentStep === tourSteps.length - 1;
  const isFirstStep = currentStep === 0;

  return (
    <>
      {/* Overlay */}
      <div 
        ref={overlayRef}
        className="fixed inset-0 z-50 bg-black/60"
        onClick={handleSkip}
      >
        {/* Highlighted element spotlight */}
        {highlightedElement && (
          <div
            className="absolute border-4 border-primary-500 rounded-lg shadow-2xl pointer-events-none"
            style={{
              top: highlightedElement.getBoundingClientRect().top + window.pageYOffset - 4,
              left: highlightedElement.getBoundingClientRect().left + window.pageXOffset - 4,
              width: highlightedElement.getBoundingClientRect().width + 8,
              height: highlightedElement.getBoundingClientRect().height + 8,
              boxShadow: '0 0 0 9999px rgba(0, 0, 0, 0.6)'
            }}
          />
        )}
      </div>

      {/* Tour Tooltip */}
      <div
        className="fixed z-[60] bg-gray-900 border border-gray-700 rounded-xl shadow-2xl p-6 w-96 max-w-[90vw]"
        style={{
          top: tooltipPosition.top,
          left: tooltipPosition.left
        }}
        onClick={(e) => e.stopPropagation()}
      >
        {/* Progress Bar */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex-1 bg-gray-800 rounded-full h-2 mr-4">
            <div 
              className="bg-gradient-to-r from-primary-500 to-secondary-500 h-2 rounded-full transition-all duration-300"
              style={{ width: `${((currentStep + 1) / tourSteps.length) * 100}%` }}
            />
          </div>
          <span className="text-xs text-gray-400">
            {currentStep + 1} / {tourSteps.length}
          </span>
        </div>

        {/* Content */}
        <div className="space-y-4">
          <h3 className="text-xl font-bold text-white">
            {step.title}
          </h3>
          <p className="text-gray-300 leading-relaxed">
            {step.content}
          </p>
        </div>

        {/* Actions */}
        <div className="flex items-center justify-between mt-6">
          <button
            onClick={handleSkip}
            className="flex items-center gap-2 px-4 py-2 text-gray-400 hover:text-white transition-colors"
          >
            <SkipForward size={16} />
            Skip Tour
          </button>

          <div className="flex gap-2">
            {!isFirstStep && (
              <button
                onClick={handlePrevious}
                className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg transition-colors"
              >
                <ArrowLeft size={16} />
                Previous
              </button>
            )}
            <button
              onClick={handleNext}
              className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 text-white rounded-lg transition-colors"
            >
              {isLastStep ? 'Finish' : 'Next'}
              <ArrowRight size={16} />
            </button>
          </div>
        </div>
      </div>
    </>
  );
}