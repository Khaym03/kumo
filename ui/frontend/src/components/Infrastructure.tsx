import { useCallback, useState } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  useNodesState,
  useEdgesState,
  addEdge,
  Position,
  type Node,
  type Edge,
  type XYPosition
} from '@xyflow/react'

import '@xyflow/react/dist/style.css'
import { getLayoutedElements } from '@/lib/utils' // Update the path if necessary
import { DevTools } from './devtools'
import { nodeTypes, edgeTypes } from './node/types'
import { nanoid } from 'nanoid'
import { ContextMenu } from '@radix-ui/react-context-menu'
import { ContextMenuContent, ContextMenuTrigger } from './ui/context-menu'
import Cmd from './cmd'

// const nodeDefaults = {
//   sourcePosition: Position.Bottom,
//   targetPosition: Position.Top,
// };

const initialNodes: Node[] = [
  {
    id: 'root',
    data: { label: 'Root' },
    position: { x: 0, y: 0 },
    sourcePosition: Position.Bottom,
    type: 'root'
  }
]

const initialEdges: Edge[] = []

const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
  initialNodes,
  initialEdges
)

const Flow = () => {
  const [nodes, setNodes, onNodesChange] = useNodesState(layoutedNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(layoutedEdges)

  const onConnect = useCallback(
    (params: any) => setEdges(els => addEdge(params, els)),
    [setEdges]
  )


  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={nodeTypes}
      edgeTypes={edgeTypes}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      colorMode="dark"
      fitView
    >
      <Background />
      <DevTools position="top-left" />
      <Controls />

      <Cmd />
    </ReactFlow>
  )
}

export default Flow
