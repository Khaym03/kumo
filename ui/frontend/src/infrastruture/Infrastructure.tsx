import { useCallback} from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  useNodesState,
  useEdgesState,
  addEdge,
} from '@xyflow/react'

import '@xyflow/react/dist/style.css'
import { getLayoutedElements } from '@/infrastruture/sort'
import { DevTools } from './devtools'
import Cmd from './cmd'
import { EdgeTypes, initialEdges, initialNodes, NodeTypes } from './constants'



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
      nodeTypes={NodeTypes}
      edgeTypes={EdgeTypes}
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
