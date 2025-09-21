import React from 'react'
import { WindowMinimise, WindowToggleMaximise, Quit } from '@wailsjs/runtime'
import { X, Minus, PictureInPicture } from 'lucide-react'
import { AnimatedThemeToggler } from './ui/animated-theme-toggler'
import { Button } from './ui/button'
import NavBar from './navbar'
import { useAppContext } from '@/context/app_ctx'
import RunButton from '@/infrastruture/run_kumo_btn'

const TitleBar: React.FC = () => {
  const { isBuilding, setIsBuilding, config } = useAppContext()
  return (
    <div
      className="flex justify-between h-10 items-center text-foreground border-b"
      style={{
        // Make the entire title bar draggable
        //@ts-expect-error linter doesn't know about wails
        '--wails-draggable': 'drag'
      }}
    >
      <nav className="flex gap-4 pl-3">
        <NavBar />
      </nav>

      <RunButton
        isBuilding={isBuilding}
        setIsBuilding={setIsBuilding}
        config={config}
        className="cursor-pointer w-28"
      />

      <div className="grid grid-flow-col items-center h-10">
        <AnimatedThemeToggler className="w-6 px-2 rounded-none" />

        <Button
          variant="ghost"
          size="icon"
          className="size-8 rounded-none h-full"
          onClick={() => WindowMinimise()}
          //@ts-expect-error linter doesn't know about wails
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <Minus className="h-6 w-6" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="size-8 rounded-none h-full"
          onClick={() => WindowToggleMaximise()}
          //@ts-expect-error linter doesn't know about wails
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <PictureInPicture className="h-6 w-6" />
        </Button>
        <Button
          onClick={() => Quit()}
          variant="default"
          size="icon"
          className="size-8 bg-background text-foreground hover:bg-destructive hover:text-foreground rounded-none aspect-square h-full"
          //@ts-expect-error linter doesn't know about wails
          style={{ '--wails-draggable': 'no-drag' }}
        >
          <X className="h-6 w-6 " />
        </Button>
      </div>
    </div>
  )
}

export default TitleBar
