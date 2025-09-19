import { memo, use, useState } from 'react'
import {
  BaseNode,
  BaseNodeHeader,
  BaseNodeHeaderTitle
} from '@/components/base-node'
import { LaptopMinimal } from 'lucide-react'
import {
  Handle,
  Position,
  useReactFlow,
  type Node,
  type NodeProps
} from '@xyflow/react'
import { nanoid } from 'nanoid'

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger
} from '@/components/ui/sheet'
import { Button } from '../ui/button'
import { Label } from '@/components/ui/label'
import { Input } from '../ui/input'
import { type BrowserConfig } from '@/lib/node_factory'
import { RadioGroup, RadioGroupItem } from '../ui/radio-group'
import { Switch } from '../ui/switch'

export const BrowserNode = memo((props: NodeProps) => {
  return (
    <BrowserNodeSheet node={props}>
      <BaseNode className="">
        <BaseNodeHeader className="border-b">
          <LaptopMinimal className="size-4" />
          <BaseNodeHeaderTitle>Browser</BaseNodeHeaderTitle>
        </BaseNodeHeader>

        {/* Target handle to receive edges */}
        <Handle
          type="target"
          position={Position.Top}
          id={`${props.id}-target`}
        />

        {/* Source handle to send edges */}
        <Handle
          type="source"
          position={Position.Bottom}
          id={`${props.id}-source`}
        />
      </BaseNode>
    </BrowserNodeSheet>
  )
})

BrowserNode.displayName = 'BrowserNode'

interface BrowserNodeProps {
  children: React.ReactNode
  node: NodeProps
}

export function BrowserNodeSheet({ children, node }: BrowserNodeProps) {
  const browserConf = node.data as BrowserConfig
  const [conf, setConf] = useState<BrowserConfig>(browserConf)
  const { setNodes } = useReactFlow()

  const handleTypeChange = (value: 'local' | 'remote') => {
    const newConf = { ...conf, type: value }
    setConf(newConf)
    setNodes(nds =>
      nds.map(n => (n.id === node.id ? { ...n, data: newConf } : n))
    )
  }

  const handleWithProxyChange = (value: boolean) => {
    const newConf = { ...conf, withProxy: value }
    setConf(newConf)
    setNodes(nds =>
      nds.map(n => (n.id === node.id ? { ...n, data: newConf } : n))
    )
  }

  const handleAddressChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newConf = { ...conf, address: e.target.value }
    setConf(newConf)
    setNodes(nds =>
      nds.map(n => (n.id === node.id ? { ...n, data: newConf } : n))
    )
  }

  return (
    <Sheet>
      <SheetTrigger asChild>{children}</SheetTrigger>
      <SheetContent className="max-w-md">
        <SheetHeader>
          <SheetTitle className="text-lg font-semibold">
            Node ID: #{node.id}
          </SheetTitle>
          <SheetDescription className="text-sm text-secondary">
            Update the browser node configuration below.
          </SheetDescription>
        </SheetHeader>

        <form className="space-y-6 px-4 py-2">
          {/* Type Selector */}
          <div>
            <Label
              className="mb-1 block text-sm font-medium"
              htmlFor={`type-${node.id}`}
            >
              Type
            </Label>
            <RadioGroup
              name="type"
              value={conf.type}
              onValueChange={handleTypeChange}
              className="flex space-x-6"
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="local" id={`local-${node.id}`} />
                <Label htmlFor={`local-${node.id}`}>Local</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="remote" id={`remote-${node.id}`} />
                <Label htmlFor={`remote-${node.id}`}>Remote</Label>
              </div>
            </RadioGroup>
          </div>

          {/* Proxy Switch */}
          <div className="flex items-center space-x-4">
            <Switch
              checked={conf.withProxy}
              onCheckedChange={handleWithProxyChange}
              id={`proxy-${node.id}`}
            />
            <Label htmlFor={`proxy-${node.id}`} className="text-sm font-medium">
              Use Proxy
            </Label>
          </div>

          {conf.type === 'remote' && (
            <div>
              <Label
                className="mb-1 block text-sm font-medium"
                htmlFor={`address-${node.id}`}
              >
                Proxy Address
              </Label>
              <Input
                id={`address-${node.id}`}
                placeholder="Enter proxy address"
                value={conf.address}
                onChange={handleAddressChange}
                className="w-full"
              />
            </div>
          )}
        </form>
      </SheetContent>
    </Sheet>
  )
}
