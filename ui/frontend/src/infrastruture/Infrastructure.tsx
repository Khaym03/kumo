import { useCallback, useEffect, useState } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  useNodesState,
  useEdgesState,
  addEdge,
  type Node
} from '@xyflow/react'

import '@xyflow/react/dist/style.css'
import { getLayoutedElements } from '@/infrastruture/sort'
import { DevTools } from './devtools'
import Cmd from './cmd'
import { EdgeTypes, initialEdges, initialNodes, NodeTypes } from './constants'
import { useAppContext } from '@/context/app_ctx'
import type { BrowserConfig } from './node_factory'
import { useBuildConfig } from './useBuildConfig'

const { nodes: layoutedNodes, edges: layoutedEdges } = getLayoutedElements(
  initialNodes,
  initialEdges
)

const Flow = () => {
  const [nodes, setNodes, onNodesChange] = useNodesState(layoutedNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(layoutedEdges)
  // const [config, setConfig] = useState<{ root?: string; browsers: any[] }>({ browsers: [] });

  const { isDarkMode } = useAppContext()

  const onConnect = useCallback(
    (params: any) => setEdges(els => addEdge(params, els)),
    [setEdges]
  )

  // const {config} = useBuildConfig(nodes, edges)

  // useEffect(() => {
  //   console.log('Built scrapper config:', config)
  // }, [config])

  return (
    <ReactFlow
      nodes={nodes}
      edges={edges}
      nodeTypes={NodeTypes}
      edgeTypes={EdgeTypes}
      onNodesChange={onNodesChange}
      onEdgesChange={onEdgesChange}
      onConnect={onConnect}
      colorMode={isDarkMode ? 'dark' : 'light'}
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
