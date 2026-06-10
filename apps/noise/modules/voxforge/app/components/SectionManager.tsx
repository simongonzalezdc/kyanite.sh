'use client';

import { Section } from '@/lib/types';
import { Trash2, Edit2 } from 'lucide-react';
import { useState } from 'react';

interface SectionManagerProps {
  sections: Section[];
  onDelete: (id: string) => void;
  onRename: (id: string, name: string) => void;
}

export default function SectionManager({ sections, onDelete, onRename }: SectionManagerProps) {
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState('');

  const startEdit = (section: Section) => {
    setEditingId(section.id);
    setEditName(section.name);
  };

  const saveEdit = (id: string) => {
    if (editName.trim()) {
      onRename(id, editName.trim());
    }
    setEditingId(null);
    setEditName('');
  };

  if (sections.length === 0) {
    return (
      <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
        <p className="text-gray-400 text-center">No sections yet. Record your first section above.</p>
      </div>
    );
  }

  return (
    <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-4">
      <h2 className="text-xl font-semibold">Sections</h2>
      <div className="space-y-3">
        {sections.map((section) => (
          <div
            key={section.id}
            className="bg-black/30 rounded-lg p-4 border border-gray-800 flex items-center justify-between"
          >
            <div className="flex-1">
              {editingId === section.id ? (
                <input
                  type="text"
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                  onBlur={() => saveEdit(section.id)}
                  onKeyDown={(e) => {
                    if (e.key === 'Enter') saveEdit(section.id);
                    if (e.key === 'Escape') setEditingId(null);
                  }}
                  className="bg-gray-800 border border-gray-700 rounded px-2 py-1 text-white"
                  autoFocus
                />
              ) : (
                <div>
                  <h3 className="font-medium">{section.name}</h3>
                  <div className="text-sm text-gray-400 mt-1">
                    Duration: {section.duration.toFixed(2)}s | BPM: {section.bpm?.bpm || 'N/A'} | Key: {section.key?.key || 'N/A'}
                  </div>
                </div>
              )}
            </div>
            <div className="flex gap-2">
              <button
                onClick={() => startEdit(section)}
                className="p-2 hover:bg-gray-700 rounded transition-colors"
              >
                <Edit2 size={16} className="text-gray-400" />
              </button>
              <button
                onClick={() => onDelete(section.id)}
                className="p-2 hover:bg-red-500/20 rounded transition-colors"
              >
                <Trash2 size={16} className="text-red-400" />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

