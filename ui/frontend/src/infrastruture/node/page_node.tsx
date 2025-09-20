import { memo } from 'react'
import {
  BaseNode,
  BaseNodeHeader,
  BaseNodeHeaderTitle
} from '@/components/base-node'
import { AppWindow } from 'lucide-react'
import { Handle, Position, type NodeProps } from '@xyflow/react'
import { nanoid } from 'nanoid'
import CustomHandle from '../handle/custom_handle'

export const PageNode = memo((props: NodeProps) => {
  return (
    <BaseNode className="">
      <BaseNodeHeader className="border-b">
        <AppWindow className="size-4" />
        <BaseNodeHeaderTitle>Page</BaseNodeHeaderTitle>
      </BaseNodeHeader>

      {/* Target handle to receive edges */}
      <CustomHandle
        type="target"
        position={Position.Top}
         id={`${props.id}-target`}
      />

      {/* Source handle to send edges */}
      <CustomHandle
        type="source"
        position={Position.Bottom}
      id={`${props.id}-source`}
      />
    </BaseNode>
  )
})

PageNode.displayName = 'PageNode'
