import { Position, type Node } from '@xyflow/react'
import { nanoid } from 'nanoid'
import type { NodeComponentType } from './constants'

export interface BrowserConfig {
  type?: 'local' | 'remote'
  withProxy?: boolean
  address?: string
}

const defaultBrowserConfig: BrowserConfig = {
  type: 'local',
  withProxy: false
}

export function nodeFactory(component: NodeComponentType) {
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

  const randPos = Math.random() * 200
  const newNode: Node = {
    id: nanoid(),
    data,
    position: { x: randPos, y: randPos },
    sourcePosition: Position.Bottom,
    type: component
  }
  return newNode
}
