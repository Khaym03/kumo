import React from 'react'
import { Search, Download, Trash2, Filter } from 'lucide-react'
import { getLevelColor, getLevelBg } from './utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

interface LogsHeaderProps {
  searchTerm: string
  setSearchTerm: (term: string) => void
  showFilters: boolean
  setShowFilters: (show: boolean) => void
  selectedLevels: Set<string>
  toggleLevel: (level: string) => void
  handleExport: () => void
  onClear?: () => void
  filteredLogCount: number
  levelCounts: Record<string, number>
}

const LogsHeader: React.FC<LogsHeaderProps> = ({
  searchTerm,
  setSearchTerm,
  showFilters,
  setShowFilters,
  selectedLevels,
  toggleLevel,
  handleExport,
  onClear,
  filteredLogCount,
  levelCounts
}) => (
  <CardHeader className="pb-4 border-b border-border/50">
    <div className="flex items-center justify-between gap-4">
      <div className="flex items-center gap-3">
         <CardTitle className="text-xl font-semibold">Logs</CardTitle>
        <CardDescription className="text-sm font-medium px-2 py-1 bg-muted rounded-md">
          {filteredLogCount.toLocaleString()} entries
        </CardDescription>
      </div>
      <div className="flex items-center gap-2">
        <Button
          variant={showFilters ? 'default' : 'outline'}
          size="sm"
          onClick={() => setShowFilters(!showFilters)}
          className="h-8"
        >
          <Filter size={14} />
          <span className="ml-1.5 hidden sm:inline">Filter</span>
        </Button>
        <Button
          variant="outline"
          size="sm"
          onClick={handleExport}
          className="h-8"
        >
          <Download size={14} />
          <span className="ml-1.5 hidden sm:inline">Export</span>
        </Button>
        {onClear && (
          <Button
            onClick={onClear}
            variant="destructive"
            size="sm"
            className="h-8"
          >
            <Trash2 size={14} />
            <span className="ml-1.5 hidden sm:inline">Clear</span>
          </Button>
        )}
      </div>
    </div>
    <div className="mt-3 relative">
      <Search
        className="absolute left-3 top-1/2 transform -translate-y-1/2"
        size={16}
      />
      <Input
        type="text"
        placeholder="Search logs..."
        value={searchTerm}
        onChange={e => setSearchTerm(e.target.value)}
        className="w-full ps-8"
      />
    </div>
    {showFilters && (
      <div className="mt-3 flex flex-wrap gap-2">
        {['info', 'warn', 'error', 'debug'].map(level => (
          <button
            key={level}
            onClick={() => toggleLevel(level)}
            className={`px-3 py-1 rounded text-xs font-medium border transition-colors ${
              selectedLevels.has(level)
                ? `${getLevelBg(level)} ${getLevelColor(level)}`
                : 'bg-gray-700 border-gray-600 text-gray-400 hover:bg-gray-600'
            }`}
          >
            {level.toUpperCase()} ({levelCounts[level] || 0})
          </button>
        ))}
      </div>
    )}
  </CardHeader>
)

export default LogsHeader
