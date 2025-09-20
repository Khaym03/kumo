// src/App.tsx

import React from 'react'
import TitleBar from './components/tittle_bar'
import Infrastructure from './infrastruture/Infrastructure'
import { AppProvider } from './context/app_ctx'

const App: React.FC = () => {
  return (
    <AppProvider>
      <div className="flex flex-col h-screen">
        {/* The custom title bar goes at the top */}
        <TitleBar />

        {/* The main content area contains your new component */}
        <main className="flex-1 overflow-auto">
          <Infrastructure />
        </main>
      </div>
    </AppProvider>
  )
}

export default App
