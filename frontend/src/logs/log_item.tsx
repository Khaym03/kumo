import React from 'react';
import type { LogEntry } from "./types"
import { getLevelColor, getLevelBg, formatTime, truncateMessage } from './utils';
import { type VirtualItem } from '@tanstack/react-virtual';
import { cn } from '@/lib/utils';

interface LogItemProps {
  log: LogEntry;
  isExpanded: boolean;
  toggleLogExpansion: () => void;
  virtualItem: VirtualItem;
  measureElement: (element: HTMLElement | null) => void;
}

const LogItem: React.FC<LogItemProps> = ({ log, isExpanded, toggleLogExpansion, virtualItem, measureElement }) => (
  <div
    key={virtualItem.key}
    data-index={virtualItem.index}
    ref={measureElement}
    style={{
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      transform: `translateY(${virtualItem.start}px)`,
    }}
    className="group flex items-start gap-3 px-4 py-2 hover:bg-accent border-b border-border/50"
  >
    <div className="max-w-[65px] w-full">
      <div className={`inline-flex items-center px-2 py-1 rounded text-xs font-medium border ${getLevelBg(log.level)} `}>
        <span className={cn(getLevelColor(log.level))}>{log.level.toUpperCase()}</span>
      </div>
    </div>
     <div className="flex-shrink-0 text-xs font-mono text-muted-foreground bg-muted/30 px-2 py-1 rounded border min-w-[80px] text-center">
        {formatTime(log.time)}
      </div>
    <div className="flex-1 min-w-0 cursor-pointer" onClick={toggleLogExpansion}>
      {isExpanded ? (
        <pre className="text-sm font-mono whitespace-pre-wrap break-words">
          {log.message}
        </pre>
      ) : (
        <div className="text-sm font-mono">
          <span className="break-words">{truncateMessage(log.message)}</span>
          {log.message.length > 100 && (
            <span className="text-primary hover:text-primary/50 ml-1 text-xs">(click to expand)</span>
          )}
        </div>
      )}
    </div>
  </div>
);

export default LogItem;