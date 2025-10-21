import { Button } from '@/components/ui/button'
import { useAppContext } from '@/context/app_ctx'
import { Panel } from '@xyflow/react'
import { SaveWorkFlow, LoadWorkFlows } from '@wailsjs/go/main/App'
import type { Node, Edge } from '@xyflow/react'

interface workflow {
  nodes: Node[]
  edges: Edge[]
}

export default function FlowBar() {
  const { nodes, edges, setNodes, setEdges } = useAppContext()

  const handleSave = async () => {
    await SaveWorkFlow('dummy', { nodes, edges })
  }

  const handleLoad = async () => {
    const workflows = (await LoadWorkFlows()) as workflow[]
    if (workflows.length > 0) {
      setNodes(workflows[0].nodes)
      setEdges(workflows[0].edges)
    }
  }
  return (
    <Panel position="top-right" className="mt-20 flex gap-4">
      <Button onClick={handleLoad}>Load</Button>
      <Button onClick={handleSave}>Save</Button>
    </Panel>
  )
}
