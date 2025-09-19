// src/App.tsx

import React from 'react';
import TitleBar from './components/tittle_bar';
import Infrastructure from './components/Infrastructure';

const App: React.FC = () => {
  return (

    <div className="flex flex-col h-screen bg-neutral-900 text-white">
      {/* The custom title bar goes at the top */}
      <TitleBar />
      
      {/* The main content area contains your new component */}
      <main className="flex-1 overflow-auto p-4">
       <Infrastructure/>
      </main>
    </div>

  );
};

export default App;