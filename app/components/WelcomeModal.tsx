'use client';

import { useState, useEffect } from 'react';
import { X, Play, Sparkles, Music, Zap } from 'lucide-react';

// Import store hooks
import { 
  useModalState, 
  useUIActions,
  useOnboardingAndHelp 
} from '@/lib/store/hooks';

export default function WelcomeModal() {
  // Use store hooks instead of props
  const { welcomeModal } = useModalState();
  const { setWelcomeModal, setOnboardingTour } = useUIActions();
  const { actions: onboardingActions } = useOnboardingAndHelp();
  
  const [isVisible, setIsVisible] = useState(false);

  const handleStartTour = () => {
    setWelcomeModal(false);
    setOnboardingTour(true);
    onboardingActions.setModalState('onboarding', true);
  };

  const handleJumpIn = () => {
    setWelcomeModal(false);
    onboardingActions.setModalState('welcome', false);
  };

  useEffect(() => {
    if (welcomeModal) {
      // Small delay for smooth entrance animation
      const timer = setTimeout(() => setIsVisible(true), 100);
      return () => clearTimeout(timer);
    } else {
      setIsVisible(false);
    }
  }, [welcomeModal]);

  if (!welcomeModal) return null;

  return (
    <div className="fixed inset-0 z-[70] flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className={`absolute inset-0 bg-black/80 transition-opacity duration-300 ${
          isVisible ? 'opacity-100' : 'opacity-0'
        }`}
      />

      {/* Modal Content */}
      <div 
        className={`relative bg-gray-900 border border-gray-700 rounded-2xl shadow-2xl max-w-lg w-full p-8 transition-all duration-300 transform ${
          isVisible ? 'scale-100 opacity-100' : 'scale-95 opacity-0'
        }`}
      >
        {/* Close Button */}
        <button
          onClick={handleJumpIn}
          className="absolute top-4 right-4 p-2 text-gray-400 hover:text-white transition-colors rounded-lg hover:bg-gray-800"
          aria-label="Close welcome modal"
        >
          <X size={20} />
        </button>

        {/* Header */}
        <div className="text-center space-y-4 mb-8">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-primary-500 to-secondary-500 rounded-full mb-4">
            <Sparkles className="text-white" size={32} />
          </div>
          
          <h1 className="text-3xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            Welcome to VoxForge
          </h1>
          
          <p className="text-gray-300 text-lg">
            Transform your voice into beautiful music with AI-powered technology
          </p>
        </div>

        {/* Features */}
        <div className="space-y-4 mb-8">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-10 h-10 bg-primary-500/20 rounded-lg flex items-center justify-center">
              <Music className="text-primary-500" size={20} />
            </div>
            <div>
              <h3 className="font-semibold text-white mb-1">Voice to Music</h3>
              <p className="text-gray-400 text-sm">
                Record your voice and watch it transform into a full musical arrangement
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-10 h-10 bg-secondary-500/20 rounded-lg flex items-center justify-center">
              <Zap className="text-secondary-500" size={20} />
            </div>
            <div>
              <h3 className="font-semibold text-white mb-1">Smart Analysis</h3>
              <p className="text-gray-400 text-sm">
                AI detects pitch, tempo, and key automatically to create perfect accompaniment
              </p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-10 h-10 bg-gradient-to-r from-primary-500/20 to-secondary-500/20 rounded-lg flex items-center justify-center">
              <Play className="text-primary-500" size={20} />
            </div>
            <div>
              <h3 className="font-semibold text-white mb-1">Easy Export</h3>
              <p className="text-gray-400 text-sm">
                Download your creation as MIDI or individual instrument stems
              </p>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="space-y-3">
          <button
            onClick={handleStartTour}
            className="w-full flex items-center justify-center gap-2 px-6 py-3 bg-gradient-to-r from-primary-500 to-secondary-500 hover:from-primary-600 hover:to-secondary-600 text-white font-medium rounded-lg transition-all duration-200 transform hover:scale-[1.02]"
          >
            <Play size={20} />
            Take a Quick Tour
          </button>
          
          <button
            onClick={handleJumpIn}
            className="w-full px-6 py-3 bg-gray-800 hover:bg-gray-700 text-white font-medium rounded-lg transition-colors"
          >
            Jump Right In
          </button>
        </div>

        {/* Footer */}
        <div className="mt-6 text-center">
          <p className="text-xs text-gray-500">
            You can always access the tour later from the help menu
          </p>
        </div>
      </div>
    </div>
  );
}