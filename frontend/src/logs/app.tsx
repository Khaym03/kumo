import React, { useState } from 'react'
import LogsSection from './logs_section' // Or the correct path to your LogsSection component
import { useAppContext } from '@/context/app_ctx'

// Define a type for your log entries
interface LogEntry {
  level: string
  message: string
  time: string
}

const App: React.FC = () => {
  const {logs} = useAppContext()

  return (
    <div className="h-full w-full flex flex-col items-center justify-center p-6">
      {/* Other components */}
      <LogsSection logs={logs} />
    </div>
  )
}

export default App
