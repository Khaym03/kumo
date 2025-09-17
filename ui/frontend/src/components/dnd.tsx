import React, { useState, useRef } from 'react'
import { DndContext, useDraggable } from '@dnd-kit/core'
import type { DragEndEvent, UniqueIdentifier } from '@dnd-kit/core'
import { nanoid } from 'nanoid'
import { restrictToParentElement } from '@dnd-kit/modifiers'
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuTrigger
} from './ui/context-menu'
import { Card } from './ui/card'

// Define the component types and data structure
type ComponentType = 'input' | 'collector' | 'filter'

interface ComponentData {
  id: UniqueIdentifier
  type: ComponentType
  position: { x: number; y: number }
}

interface CanvasItemProps {
  id: UniqueIdentifier
  type: ComponentType
  initialPosition?: { x: number; y: number }
}

// Draggable component logic (used for canvas items)
const CanvasItem: React.FC<CanvasItemProps> = ({
  id,
  type,
  initialPosition = { x: 0, y: 0 }
}) => {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({ id })

  const finalTransform = {
    x: initialPosition.x + (transform?.x || 0),
    y: initialPosition.y + (transform?.y || 0)
  }

  const style = {
    left: finalTransform.x,
    top: finalTransform.y
  }

  return (
    <Card
      ref={setNodeRef}
      style={style}
      {...listeners}
      {...attributes}
      className="absolute p-3 cursor-grab border-2 border-primary"
    >
      {type}
    </Card>
  )
}

// Main component with drag-and-drop logic
const DragAndDropCanvas: React.FC = () => {
  const [items, setItems] = useState<ComponentData[]>([])
  const canvasRef = useRef<HTMLDivElement>(null)
  let menuCoords = { x: 0, y: 0 }

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, delta } = event
    const existingItem = items.find(item => item.id === active.id)

    if (existingItem) {
      setItems(prevItems =>
        prevItems.map(item =>
          item.id === existingItem.id
            ? {
                ...item,
                position: {
                  x: item.position.x + delta.x,
                  y: item.position.y + delta.y
                }
              }
            : item
        )
      )
    }
  }

  const addComponent = (type: ComponentType) => {
    const newComponent: ComponentData = {
      id: nanoid(),
      type: type,
      position: { x: menuCoords.x, y: menuCoords.y }
    }
    setItems(prevItems => [...prevItems, newComponent])
  }

  return (
    <DndContext onDragEnd={handleDragEnd} modifiers={[restrictToParentElement]}>
      <ContextMenu
        onOpenChange={isOpen => {
          if (isOpen && canvasRef.current) {
            const rect = canvasRef.current.getBoundingClientRect()
            menuCoords = { x: rect.left, y: rect.top }
          }
        }}
      >
        <ContextMenuTrigger asChild>
          <div
            ref={canvasRef}
            id="canvas-area"
            className="w-full h-full p-4 relative border-2 border-dashed bg-muted"
          >
            {items.map(item => (
              <CanvasItem
                key={item.id}
                id={item.id}
                type={item.type}
                initialPosition={item.position}
              />
            ))}
          </div>
        </ContextMenuTrigger>
        <ContextMenuContent className="w-32">
          <ContextMenuItem onClick={() => addComponent('input')}>
            Input
          </ContextMenuItem>
          <ContextMenuItem onClick={() => addComponent('collector')}>
            Collector
          </ContextMenuItem>
          <ContextMenuItem onClick={() => addComponent('filter')}>
            Filter
          </ContextMenuItem>
        </ContextMenuContent>
      </ContextMenu>
    </DndContext>
  )
}

export default DragAndDropCanvas
