// hooks/useBuildConfig.ts
import { useEffect, useState } from 'react'
import { buildConfigFromGraph } from '@/infrastruture/utils/graph'
import type { Node, Edge } from '@xyflow/react'

export function useBuildConfig(nodes: Node[], edges: Edge[]) {
  const [config, setConfig] = useState<{ root?: string; browsers: any[] }>({
    browsers: []
  })
 
  useEffect(() => {
    setConfig(buildConfigFromGraph(nodes, edges))
  }, [nodes, edges])

  return { config }
}
