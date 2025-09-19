import { memo } from 'react'
import {
  BaseNode,
  BaseNodeHeader,
  BaseNodeHeaderTitle
} from '@/components/base-node'
import { AppWindow } from 'lucide-react'
import { Handle, Position } from '@xyflow/react'
import { nanoid } from 'nanoid'

export const PageNode = memo(() => {
  return (
    <BaseNode className="">
      <BaseNodeHeader className="border-b">
        <AppWindow className="size-4" />
        <BaseNodeHeaderTitle>Page</BaseNodeHeaderTitle>
      </BaseNodeHeader>

      {/* Target handle to receive edges */}
      <Handle
        type="target"
        position={Position.Top}
        id={nanoid()}
      />

      {/* Source handle to send edges */}
      <Handle
        type="source"
        position={Position.Bottom}
       id={nanoid()}
      />
    </BaseNode>
  )
})

PageNode.displayName = 'PageNode'
