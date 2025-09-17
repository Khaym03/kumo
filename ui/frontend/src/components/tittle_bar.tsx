// src/components/TitleBar.tsx

import React from 'react';
import {WindowMinimise, WindowToggleMaximise, Quit} from '@wailsjs/runtime';
import {X, Minus, PictureInPicture} from 'lucide-react';

const TitleBar: React.FC = () => {
  return (
    <div
      className="flex justify-between items-center text-foreground"
      style={{
        // Make the entire title bar draggable
        //@ts-expect-error linter doesn't know about wails
        '--wails-draggable': 'drag',
      }}
    >
      <div className='pl-3 font-medium'>Kumo</div>
      <div
        className="flex gap-2"
      >
        <button
          onClick={() => WindowMinimise()}
          className="hover:bg-neutral-600 px-2 py-1"
          //@ts-expect-error linter doesn't know about wails
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <Minus size={20} />
        </button>
        <button
          onClick={() => WindowToggleMaximise()}
          className="hover:bg-neutral-600 px-2 py-1"
          //@ts-expect-error linter doesn't know about wails
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <PictureInPicture size={20}/>
        </button>
        <button
          onClick={() => Quit()}
          className="hover:bg-destructive hover:text-white px-2 py-1"
          //@ts-expect-error linter doesn't know about wails
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <X size={20} />
        </button>
      </div>
    </div>
  );
};

export default TitleBar;