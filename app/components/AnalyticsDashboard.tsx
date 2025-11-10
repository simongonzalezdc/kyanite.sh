'use client';

import React, { useState, useEffect, useRef } from 'react';
import { useAnalytics } from './AnalyticsProvider';
import { AnalyticsEvent, AnalyticsSession, PerformanceMetrics } from '@/lib/types/analytics';

interface AnalyticsDashboardProps {
  isOpen: boolean;
  onClose: () => void;
}

interface EventSummary {
  type: string;
  count: number;
  lastSeen: number;
}

interface PerformanceSummary {
  metric: string;
  value: number;
  unit: string;
  trend: 'up' | 'down' | 'stable';
}

export function AnalyticsDashboard({ isOpen, onClose }: AnalyticsDashboardProps) {
  const { getSession, exportData, getCurrentMetrics, getPerformanceReport } = useAnalytics();
  const [events, setEvents] = useState<AnalyticsEvent[]>([]);
  const [session, setSession] = useState<AnalyticsSession | null>(null);
  const [metrics, setMetrics] = useState<PerformanceMetrics>({});
  const [performanceReport, setPerformanceReport] = useState<string>('');
  const [activeTab, setActiveTab] = useState<'events' | 'performance' | 'session' | 'export'>('events');
  const [eventFilter, setEventFilter] = useState<string>('all');
  const [autoRefresh, setAutoRefresh] = useState(true);
  const refreshInterval = useRef<NodeJS.Timeout | null>(null);

  // Load data when dashboard opens
  useEffect(() => {
    if (isOpen) {
      refreshData();
      
      if (autoRefresh) {
        refreshInterval.current = setInterval(refreshData, 2000);
      }
    } else {
      if (refreshInterval.current) {
        clearInterval(refreshInterval.current);
        refreshInterval.current = null;
      }
    }

    return () => {
      if (refreshInterval.current) {
        clearInterval(refreshInterval.current);
      }
    };
  }, [isOpen, autoRefresh]);

  const refreshData = () => {
    // Get current session
    const currentSession = getSession();
    setSession(currentSession);

    // Get current metrics
    const currentMetrics = getCurrentMetrics();
    setMetrics(currentMetrics);

    // Get performance report
    const report = getPerformanceReport();
    setPerformanceReport(report);

    // Get stored events (in a real implementation, you'd fetch from your analytics service)
    // For now, we'll use a mock implementation
    setEvents([]);
  };

  const handleExport = () => {
    const data = exportData();
    const blob = new Blob([data], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `voxforge-analytics-${new Date().toISOString().split('T')[0]}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  const filteredEvents = events.filter(event => 
    eventFilter === 'all' || event.type === eventFilter
  );

  const getEventSummary = (): EventSummary[] => {
    const summary: Record<string, EventSummary> = {};
    
    events.forEach(event => {
      if (!summary[event.type]) {
        summary[event.type] = {
          type: event.type,
          count: 0,
          lastSeen: 0
        };
      }
      summary[event.type].count++;
      summary[event.type].lastSeen = Math.max(summary[event.type].lastSeen, event.timestamp);
    });

    return Object.values(summary).sort((a, b) => b.count - a.count);
  };

  const getPerformanceSummary = (): PerformanceSummary[] => {
    const summary: PerformanceSummary[] = [];
    
    if (metrics.fcp !== undefined) {
      summary.push({
        metric: 'First Contentful Paint',
        value: metrics.fcp,
        unit: 'ms',
        trend: 'stable'
      });
    }
    
    if (metrics.lcp !== undefined) {
      summary.push({
        metric: 'Largest Contentful Paint',
        value: metrics.lcp,
        unit: 'ms',
        trend: 'stable'
      });
    }
    
    if (metrics.fid !== undefined) {
      summary.push({
        metric: 'First Input Delay',
        value: metrics.fid,
        unit: 'ms',
        trend: 'stable'
      });
    }
    
    if (metrics.cls !== undefined) {
      summary.push({
        metric: 'Cumulative Layout Shift',
        value: metrics.cls,
        unit: 'score',
        trend: 'stable'
      });
    }
    
    if (metrics.ttfb !== undefined) {
      summary.push({
        metric: 'Time to First Byte',
        value: metrics.ttfb,
        unit: 'ms',
        trend: 'stable'
      });
    }
    
    if (metrics.memoryUsage !== undefined) {
      summary.push({
        metric: 'Memory Usage',
        value: metrics.memoryUsage / 1024 / 1024,
        unit: 'MB',
        trend: 'stable'
      });
    }
    
    return summary;
  };

  const formatTimestamp = (timestamp: number) => {
    return new Date(timestamp).toLocaleTimeString();
  };

  const formatDuration = (ms: number) => {
    if (ms < 1000) return `${ms}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-900 text-white rounded-lg shadow-xl w-full max-w-4xl max-h-[80vh] overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-700">
          <h2 className="text-xl font-semibold">Analytics Dashboard</h2>
          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2 text-sm">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                className="rounded"
              />
              Auto Refresh
            </label>
            <button
              onClick={refreshData}
              className="px-3 py-1 bg-blue-600 hover:bg-blue-700 rounded text-sm"
            >
              Refresh
            </button>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-white"
            >
              ✕
            </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="flex border-b border-gray-700">
          {(['events', 'performance', 'session', 'export'] as const).map(tab => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 capitalize ${
                activeTab === tab
                  ? 'bg-gray-800 text-white border-b-2 border-blue-500'
                  : 'text-gray-400 hover:text-white'
              }`}
            >
              {tab}
            </button>
          ))}
        </div>

        {/* Content */}
        <div className="p-4 overflow-y-auto max-h-[60vh]">
          {/* Events Tab */}
          {activeTab === 'events' && (
            <div className="space-y-4">
              <div className="flex items-center gap-4">
                <label className="text-sm">Filter:</label>
                <select
                  value={eventFilter}
                  onChange={(e) => setEventFilter(e.target.value)}
                  className="bg-gray-800 rounded px-2 py-1 text-sm"
                >
                  <option value="all">All Events</option>
                  <option value="interaction">Interactions</option>
                  <option value="performance">Performance</option>
                  <option value="audio">Audio</option>
                  <option value="feature">Features</option>
                  <option value="error">Errors</option>
                  <option value="pageview">Page Views</option>
                </select>
              </div>

              {/* Event Summary */}
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
                {getEventSummary().map(summary => (
                  <div key={summary.type} className="bg-gray-800 rounded p-3">
                    <div className="text-xs text-gray-400 capitalize">{summary.type}</div>
                    <div className="text-lg font-semibold">{summary.count}</div>
                    <div className="text-xs text-gray-500">
                      Last: {formatTimestamp(summary.lastSeen)}
                    </div>
                  </div>
                ))}
              </div>

              {/* Event List */}
              <div className="space-y-2">
                <h3 className="text-lg font-semibold">Recent Events</h3>
                <div className="bg-gray-800 rounded overflow-hidden">
                  {filteredEvents.length === 0 ? (
                    <div className="p-4 text-center text-gray-500">No events to display</div>
                  ) : (
                    <div className="max-h-64 overflow-y-auto">
                      {filteredEvents.slice(-20).reverse().map((event, index) => (
                        <div key={index} className="border-b border-gray-700 p-3 last:border-b-0">
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                              <span className="text-xs bg-gray-700 px-2 py-1 rounded capitalize">
                                {event.type}
                              </span>
                              <span className="text-sm">
                                {event.type === 'interaction' && (event as any).element}
                                {event.type === 'performance' && (event as any).metric}
                                {event.type === 'audio' && (event as any).action}
                                {event.type === 'feature' && `${(event as any).feature}:${(event as any).action}`}
                                {event.type === 'error' && (event as any).error.message}
                                {event.type === 'pageview' && (event as any).page}
                              </span>
                            </div>
                            <span className="text-xs text-gray-500">
                              {formatTimestamp(event.timestamp)}
                            </span>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Performance Tab */}
          {activeTab === 'performance' && (
            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Performance Metrics</h3>
              
              {/* Performance Summary */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {getPerformanceSummary().map(summary => (
                  <div key={summary.metric} className="bg-gray-800 rounded p-4">
                    <div className="text-sm text-gray-400">{summary.metric}</div>
                    <div className="text-2xl font-semibold">
                      {summary.value.toFixed(summary.unit === 'score' ? 3 : 1)}
                      <span className="text-sm text-gray-400 ml-1">{summary.unit}</span>
                    </div>
                    <div className="text-xs text-gray-500">
                      {summary.trend === 'up' && '↑'}
                      {summary.trend === 'down' && '↓'}
                      {summary.trend === 'stable' && '→'}
                    </div>
                  </div>
                ))}
              </div>

              {/* Performance Report */}
              <div className="space-y-2">
                <h3 className="text-lg font-semibold">Performance Report</h3>
                <div className="bg-gray-800 rounded p-4">
                  <pre className="text-xs text-green-400 overflow-x-auto">
                    {performanceReport || 'No performance data available'}
                  </pre>
                </div>
              </div>
            </div>
          )}

          {/* Session Tab */}
          {activeTab === 'session' && (
            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Session Information</h3>
              
              {session ? (
                <div className="bg-gray-800 rounded p-4 space-y-3">
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <div className="text-sm text-gray-400">Session ID</div>
                      <div className="font-mono text-xs">{session.id}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-400">Duration</div>
                      <div>{formatDuration(Date.now() - session.startTime)}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-400">Page Views</div>
                      <div>{session.pageViews}</div>
                    </div>
                    <div>
                      <div className="text-sm text-gray-400">Events</div>
                      <div>{session.events}</div>
                    </div>
                  </div>
                  
                  <div>
                    <div className="text-sm text-gray-400">User Agent</div>
                    <div className="text-xs truncate">{session.userAgent}</div>
                  </div>
                  
                  <div>
                    <div className="text-sm text-gray-400">Viewport</div>
                    <div>{session.viewport.width} × {session.viewport.height}</div>
                  </div>
                  
                  <div>
                    <div className="text-sm text-gray-400">Start Time</div>
                    <div>{new Date(session.startTime).toLocaleString()}</div>
                  </div>
                </div>
              ) : (
                <div className="bg-gray-800 rounded p-4 text-center text-gray-500">
                  No session data available
                </div>
              )}
            </div>
          )}

          {/* Export Tab */}
          {activeTab === 'export' && (
            <div className="space-y-4">
              <h3 className="text-lg font-semibold">Export Analytics Data</h3>
              
              <div className="bg-gray-800 rounded p-4">
                <p className="text-sm text-gray-300 mb-4">
                  Export all analytics data for this session. The data includes all events, 
                  performance metrics, and session information in JSON format.
                </p>
                
                <button
                  onClick={handleExport}
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded"
                >
                  Export Data
                </button>
              </div>
              
              <div className="bg-gray-800 rounded p-4">
                <h4 className="font-semibold mb-2">Data Privacy</h4>
                <ul className="text-sm text-gray-300 space-y-1">
                  <li>• No personal data is collected</li>
                  <li>• All data is anonymized</li>
                  <li>• IP addresses are not stored</li>
                  <li>• You can opt out at any time</li>
                  <li>• Data is stored locally in your browser</li>
                </ul>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// Development-only wrapper component
export function AnalyticsDashboardDev() {
  const [isOpen, setIsOpen] = useState(false);
  
  // Only show in development
  if (process.env.NODE_ENV !== 'development') {
    return null;
  }

  return (
    <>
      {/* Floating button to open dashboard */}
      <button
        onClick={() => setIsOpen(true)}
        className="fixed bottom-4 right-4 bg-blue-600 hover:bg-blue-700 text-white rounded-full p-3 shadow-lg z-40"
        title="Open Analytics Dashboard"
      >
        📊
      </button>
      
      <AnalyticsDashboard isOpen={isOpen} onClose={() => setIsOpen(false)} />
    </>
  );
}