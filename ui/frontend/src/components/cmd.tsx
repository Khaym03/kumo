import {
  Command,
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList
} from '@/components/ui/command'

import { NodeComponents, type NodeComponentTypes } from '@/components/node/types'
import React from 'react'
import { useReactFlow} from '@xyflow/react'
import { nodeFactory } from '@/lib/node_factory'

export default function Cmd() {
  const [open, setOpen] = React.useState(false)
  const { setNodes } = useReactFlow()

  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen(open => !open)
      }
    }
    document.addEventListener('keydown', down)
    return () => document.removeEventListener('keydown', down)
  }, [])

  const handleAddNode = (component: NodeComponentTypes) => {
    setNodes(nodes => [...nodes, nodeFactory(component)])
    setOpen(false)
  }

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <Command>
        <CommandInput placeholder="Type a command or search..." />
        <CommandList>
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup heading="Components">
            {NodeComponents.map((component, index) => (
              <CommandItem key={index}  onSelect={() => handleAddNode(component)}>
                {component}
              </CommandItem>
            ))}
          </CommandGroup>
        </CommandList>
      </Command>
    </CommandDialog>
  )
}


