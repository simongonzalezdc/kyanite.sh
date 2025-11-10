'use client';

import { useState } from 'react';
import { ChevronDown, ChevronUp, Mic, Activity, Music, Download, Play } from 'lucide-react';

export default function QuickStartGuide() {
  const [isExpanded, setIsExpanded] = useState(true);

  const steps = [
    {
      id: 'record',
      icon: Mic,
      title: 'Record Your Voice',
      description: 'Click the microphone button to record your melody, vocals, or any sound you want to transform into music.',
      color: 'primary'
    },
    {
      id: 'analyze',
      icon: Activity,
      title: 'Analyze & Generate',
      description: 'VoxForge automatically analyzes your recording and generates accompaniment with your selected instruments.',
      color: 'secondary'
    },
    {
      id: 'export',
      icon: Download,
      title: 'Export Your Music',
      description: 'Download your creation as MIDI files or individual instrument stems for further production.',
      color: 'accent'
    }
  ];

  const colorClasses = {
    primary: {
      bg: 'bg-primary-500/20',
      border: 'border-primary-500/50',
      text: 'text-primary-500',
      iconBg: 'bg-primary-500'
    },
    secondary: {
      bg: 'bg-secondary-500/20',
      border: 'border-secondary-500/50',
      text: 'text-secondary-500',
      iconBg: 'bg-secondary-500'
    },
    accent: {
      bg: 'bg-gradient-to-r from-primary-500/20 to-secondary-500/20',
      border: 'border-primary-500/50',
      text: 'text-primary-500',
      iconBg: 'bg-gradient-to-r from-primary-500 to-secondary-500'
    }
  };

  return (
    <div className="bg-gray-900 rounded-xl border border-gray-800 overflow-hidden">
      {/* Header */}
      <div 
        className="flex items-center justify-between p-6 cursor-pointer hover:bg-gray-800/50 transition-colors"
        onClick={() => setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-gradient-to-r from-primary-500 to-secondary-500 rounded-lg flex items-center justify-center">
            <Play className="text-white" size={20} />
          </div>
          <div>
            <h2 className="text-xl font-semibold">Quick Start Guide</h2>
            <p className="text-sm text-gray-400">Get started in 3 simple steps</p>
          </div>
        </div>
        
        <button 
          className="p-2 text-gray-400 hover:text-white transition-colors rounded-lg hover:bg-gray-700"
          aria-label={isExpanded ? 'Collapse guide' : 'Expand guide'}
        >
          {isExpanded ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
        </button>
      </div>

      {/* Content */}
      {isExpanded && (
        <div className="px-6 pb-6">
          <div className="space-y-4">
            {steps.map((step, index) => {
              const Icon = step.icon;
              const colors = colorClasses[step.color as keyof typeof colorClasses];
              
              return (
                <div 
                  key={step.id}
                  className="flex gap-4 p-4 bg-black/30 rounded-lg border border-gray-800 hover:border-gray-700 transition-colors"
                >
                  {/* Step Number */}
                  <div className="flex-shrink-0">
                    <div className={`w-8 h-8 ${colors.iconBg} rounded-full flex items-center justify-center text-white font-bold text-sm`}>
                      {index + 1}
                    </div>
                  </div>

                  {/* Icon */}
                  <div className={`flex-shrink-0 w-12 h-12 ${colors.bg} ${colors.border} border rounded-lg flex items-center justify-center`}>
                    <Icon className={colors.text} size={24} />
                  </div>

                  {/* Content */}
                  <div className="flex-1">
                    <h3 className="font-semibold text-white mb-1">
                      {step.title}
                    </h3>
                    <p className="text-sm text-gray-400 leading-relaxed">
                      {step.description}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>

          {/* Tips Section */}
          <div className="mt-6 p-4 bg-gradient-to-r from-primary-500/10 to-secondary-500/10 rounded-lg border border-primary-500/20">
            <h4 className="font-medium text-white mb-2 flex items-center gap-2">
              <Play size={16} className="text-primary-500" />
              Pro Tips
            </h4>
            <ul className="space-y-1 text-sm text-gray-300">
              <li>• Record in a quiet environment for best results</li>
              <li>• Try humming or singing a simple melody</li>
              <li>• Experiment with different instrument combinations</li>
              <li>• Use the trim feature to remove unwanted parts</li>
            </ul>
          </div>
        </div>
      )}
    </div>
  );
}