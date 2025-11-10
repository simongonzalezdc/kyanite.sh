'use client';

import { InstrumentType } from '@/lib/types';
import { Drum, Music, Guitar } from 'lucide-react';

interface InstrumentPickerProps {
  selected: InstrumentType[];
  onChange: (instruments: InstrumentType[]) => void;
}

export default function InstrumentPicker({ selected, onChange }: InstrumentPickerProps) {
  const instruments: { type: InstrumentType; label: string; icon: React.ReactNode }[] = [
    { type: 'drums', label: 'Drums', icon: <Drum size={20} /> },
    { type: 'bass', label: 'Bass', icon: <Music size={20} /> },
    { type: 'chords', label: 'Chords', icon: <Guitar size={20} /> },
  ];

  const toggleInstrument = (type: InstrumentType) => {
    if (selected.includes(type)) {
      if (selected.length > 1) {
        onChange(selected.filter(i => i !== type));
      }
    } else {
      onChange([...selected, type]);
    }
  };

  return (
    <div className="space-y-4">
      <h3 className="font-medium">Select Instruments</h3>
      <div className="flex gap-4">
        {instruments.map(({ type, label, icon }) => (
          <button
            key={type}
            onClick={() => toggleInstrument(type)}
            disabled={selected.length === 1 && selected.includes(type)}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg border transition-colors ${
              selected.includes(type)
                ? 'bg-primary-500/20 border-primary-500 text-primary-500'
                : 'bg-gray-800 border-gray-700 text-gray-400 hover:border-gray-600'
            } ${
              selected.length === 1 && selected.includes(type)
                ? 'opacity-50 cursor-not-allowed'
                : 'cursor-pointer'
            }`}
          >
            {icon}
            {label}
          </button>
        ))}
      </div>
      {selected.length === 1 && (
        <p className="text-sm text-gray-400">At least one instrument must be selected</p>
      )}
    </div>
  );
}

