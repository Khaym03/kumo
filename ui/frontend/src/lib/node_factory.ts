import { type NodeComponentTypes } from '@/components/node/types'
import { Position, type Node } from '@xyflow/react'
import { nanoid } from 'nanoid'

export interface BrowserConfig {
  type?: 'local' | 'remote'
  withProxy?: boolean
  address?: string
}

const defaultBrowserConfig: BrowserConfig = {
  type: 'local',
  withProxy: false
}

export function nodeFactory(component: NodeComponentTypes) {
  let data = {}
  switch (component) {
    case 'browser':
      data = defaultBrowserConfig
      break
    case 'root':
      data = {}
      break
    case 'page':
      data = {}
      break
  }

  const randPos = Math.random() * 300
  const newNode: Node = {
    id: nanoid(),
    data,
    position: { x: randPos, y: randPos },
    sourcePosition: Position.Bottom,
    type: component
  }
  return newNode
}
