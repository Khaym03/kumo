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
import { Button } from '@/components/ui/button'
import { CancelKumo, RunKumo } from '@wailsjs/go/main/App'
import { main } from '@wailsjs/go/models'
import { Loader2Icon, LoaderCircle } from 'lucide-react'
import RunButton from './run_kumo_btn'

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
  const { config } = useBuildConfig(nodes, edges)

  const [isBuilding, setIsBuilding] = useState(false)

  // useEffect(() => {
  //   console.log('Built scrapper config:', config)
  // }, [config])

  return (
    <div className="relative h-full w-full">
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

      <RunButton
        isBuilding={isBuilding}
        setIsBuilding={setIsBuilding}
        config={config}
        className="absolute top-4 right-4 cursor-pointer w-28"
      />
    </div>
  )
}

export default Flow
