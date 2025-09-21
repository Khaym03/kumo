import React, { useState, useEffect, useRef, useCallback } from 'react'
import { useVirtualizer } from '@tanstack/react-virtual'
import type { LogEntry } from './types'
import LogsHeader from './logs_header'
import LogList from './log_list'
import LogsFooter from './logs_footer'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'

interface LogsSectionProps {
  logs: LogEntry[]
  onClear?: () => void
}

const LogsSection: React.FC<LogsSectionProps> = ({ logs, onClear }) => {
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedLevels, setSelectedLevels] = useState<Set<string>>(
    new Set(['info', 'warn', 'error', 'debug'])
  )
  const [autoScroll, setAutoScroll] = useState(true)
  const [showFilters, setShowFilters] = useState(false)
  const [expandedLogs, setExpandedLogs] = useState<Set<number>>(new Set())
  const [containerHeight, setContainerHeight] = useState(400)

  const parentRef = useRef<HTMLDivElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  // --- Logic remains the same, just the rendering is delegated ---
  const filteredLogs = logs.filter(log => {
    const matchesSearch =
      log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
      log.level.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesLevel = selectedLevels.has(log.level.toLowerCase())
    return matchesSearch && matchesLevel
  })

  // Dynamic height calculation (remains in main component)
  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const updateHeight = () => {
      const rect = container.getBoundingClientRect()
      const headerHeight =
        container.querySelector('.bg-gray-800')?.getBoundingClientRect()
          .height || 0
      const footerHeight =
        container.querySelector('.border-t')?.getBoundingClientRect().height ||
        0
      const availableHeight = rect.height - headerHeight - footerHeight
      setContainerHeight(Math.max(200, availableHeight)) // Minimum 200px
    }

    // Initial calculation
    updateHeight()

    // Use ResizeObserver for better performance
    const resizeObserver = new ResizeObserver(() => {
      updateHeight()
    })

    resizeObserver.observe(container)

    console.log('resizeObserver', resizeObserver)

    return () => {
      resizeObserver.disconnect()
    }
  }, [])

  // Virtualizer setup (remains in main component)
  const rowVirtualizer = useVirtualizer({
    count: filteredLogs.length,
    getScrollElement: () => parentRef.current,
    estimateSize: useCallback(
      (index: number) => {
        const isExpanded = expandedLogs.has(index)
        const log = filteredLogs[index]
        if (!log) return 60
        const baseHeight = 60
        if (isExpanded) {
          const lines = Math.ceil(log.message.length / 80)
          return Math.max(baseHeight, lines * 20 + 40)
        }
        return baseHeight
      },
      [expandedLogs, filteredLogs]
    ),
    overscan: 5
  })

  // Auto-scroll logic (remains in main component)
  useEffect(() => {
    if (autoScroll && filteredLogs.length > 0) {
      rowVirtualizer.scrollToIndex(filteredLogs.length - 1, { align: 'end' })
    }
  }, [filteredLogs.length, autoScroll, rowVirtualizer])

  useEffect(() => {
    const element = parentRef.current
    if (!element) return
    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = element
      const isNearBottom = scrollTop + clientHeight >= scrollHeight - 100
      setAutoScroll(isNearBottom)
    }
    element.addEventListener('scroll', handleScroll)
    return () => element.removeEventListener('scroll', handleScroll)
  }, [])

  const toggleLevel = (level: string) => {
    setSelectedLevels(prev => {
      const newSet = new Set(prev)
      if (newSet.has(level)) {
        newSet.delete(level)
      } else {
        newSet.add(level)
      }
      return newSet
    })
  }

  const handleExport = () => {
    /* ... export logic ... */
  }
  const toggleLogExpansion = (index: number) => {
    setExpandedLogs(prev => {
      const newSet = new Set(prev)
      if (newSet.has(index)) newSet.delete(index)
      else newSet.add(index)
      return newSet
    })
  }
  const handleAutoScrollChange = (checked: boolean) => {
    setAutoScroll(checked)
    if (checked && filteredLogs.length > 0) {
      rowVirtualizer.scrollToIndex(filteredLogs.length - 1, { align: 'end' })
    }
  }
  const levelCounts = logs.reduce((acc, log) => {
    acc[log.level.toLowerCase()] = (acc[log.level.toLowerCase()] || 0) + 1
    return acc
  }, {} as Record<string, number>)

  return (
    <Card
      ref={containerRef}
      className="rounded-lg overflow-hidden h-[90dvh] flex flex-col w-full gap-0 border-0 bg-transparent"
    >
      <LogsHeader
        searchTerm={searchTerm}
        setSearchTerm={setSearchTerm}
        showFilters={showFilters}
        setShowFilters={setShowFilters}
        selectedLevels={selectedLevels}
        toggleLevel={toggleLevel}
        handleExport={handleExport}
        onClear={onClear}
        filteredLogCount={filteredLogs.length}
        levelCounts={levelCounts}
      />

      <LogList
        parentRef={parentRef}
        filteredLogs={filteredLogs}
        virtualizer={rowVirtualizer}
        containerHeight={containerHeight}
        expandedLogs={expandedLogs}
        toggleLogExpansion={toggleLogExpansion}
      />

      <CardFooter className="pt-4 flex items-center justify-between">
        <LogsFooter
          autoScroll={autoScroll}
          handleAutoScrollChange={handleAutoScrollChange}
          filteredCount={filteredLogs.length}
          totalCount={logs.length}
        />
      </CardFooter>
    </Card>
  )
}

export default LogsSection
