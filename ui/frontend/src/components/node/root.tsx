import { memo } from 'react'
import { BaseNode, BaseNodeContent } from '@/components/base-node'
import { Network } from 'lucide-react'
import { Handle, Position } from '@xyflow/react'

const numberOfHandles = 4

export const RootNode = memo(() => {
  // Calculate left positions for handles as percentages
  const leftPositions = Array.from({ length: numberOfHandles }, (_, i) => `${(i + 1) * 100 / (numberOfHandles + 1)}%`)

  return (
    <BaseNode>
      <BaseNodeContent>
        <div className="flex gap-4 items-center">
          <Network size={16} /> Root
        </div>
      </BaseNodeContent>

      {leftPositions.map((left, index) => (
        <Handle
          key={`bottom-${index}`}
          type="source"
          position={Position.Bottom}
          id={`bottom-${index}`}
          style={{ left, width: 10, height: 10 }} // customize handle style here
        />
      ))}
    </BaseNode>
  )
})

RootNode.displayName = 'RootNode'
