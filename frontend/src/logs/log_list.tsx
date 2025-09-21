import React from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import type { LogEntry } from './types'
import LogItem from './log_item'
import { CardContent } from '@/components/ui/card'

interface LogListProps {
  parentRef: React.RefObject<HTMLDivElement | null>
  filteredLogs: LogEntry[]
  virtualizer: ReturnType<typeof useVirtualizer<HTMLDivElement, Element>>
  containerHeight: number
  expandedLogs: Set<number>
  toggleLogExpansion: (index: number) => void
}

const LogList: React.FC<LogListProps> = ({
  parentRef,
  filteredLogs,
  virtualizer,
  containerHeight,
  expandedLogs,
  toggleLogExpansion
}) => {
  console.log(containerHeight)
  return (
    <CardContent
      ref={parentRef}
      className="overflow-y-auto scrollbar-thin scrollbar-track-gray-800 scrollbar-thumb-gray-600 px-0"
      style={{ height: `${containerHeight}px` }}
    >
      {filteredLogs.length === 0 ? (
        <div className="flex items-center justify-center h-full">
          <div className="text-center">
            <div className="text-2xl mb-2">📋</div>
            <div>No logs to display</div>
            {/* Note: searchTerm prop passed here for this conditional */}
            {/* <div className="text-sm mt-1">Try adjusting your search or filters</div> */}
          </div>
        </div>
      ) : (
        <div
          style={{
            height: `${virtualizer.getTotalSize()}px`,
            width: '100%',
            position: 'relative'
          }}
        >
          {virtualizer.getVirtualItems().map(virtualItem => {
            const log = filteredLogs[virtualItem.index]
            return (
              <LogItem
                key={virtualItem.key}
                virtualItem={virtualItem}
                log={log}
                isExpanded={expandedLogs.has(virtualItem.index)}
                toggleLogExpansion={() => toggleLogExpansion(virtualItem.index)}
                measureElement={virtualizer.measureElement}
              />
            )
          })}
        </div>
      )}
    </CardContent>
  )
}

export default LogList
