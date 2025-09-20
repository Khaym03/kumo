import { memo, useEffect } from 'react'
import { BaseNode, BaseNodeContent } from '@/components/base-node'
import { LaptopMinimal } from 'lucide-react'
import { Position, useReactFlow, type NodeProps } from '@xyflow/react'

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger
} from '@/components/ui/sheet'
import { Input } from '../../components/ui/input'
import { RadioGroup, RadioGroupItem } from '../../components/ui/radio-group'
import { Switch } from '../../components/ui/switch'
import CustomHandle from '../handle/custom_handle'

// Form and Zod imports
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from '@/components/ui/form'

import { type BrowserConfig } from '@/infrastruture/node_factory'

type BrowserFormValues = z.infer<typeof browserFormSchema>

export const BrowserNode = memo((props: NodeProps) => {
  return (
    <BrowserNodeSheet node={props}>
      <BaseNode>
        <BaseNodeContent>
          <div className="flex gap-4 items-center">
            <LaptopMinimal className="size-4" />
            <div>Browser</div>
          </div>
        </BaseNodeContent>
        <CustomHandle
          type="target"
          position={Position.Top}
          id={`${props.id}-target`}
        />
        <CustomHandle
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

const browserFormSchema = z.object({
  type: z.enum(['local', 'remote'], {
    message: 'Please select a browser type.'
  }),
  withProxy: z.boolean(),
  address: z.string()
})

export function BrowserNodeSheet({ children, node }: BrowserNodeProps) {
  const { setNodes } = useReactFlow()

  const initialData = node.data as unknown as BrowserConfig


  const form = useForm<BrowserFormValues>({
    resolver: zodResolver(browserFormSchema),
    defaultValues: {
      type: initialData.type,
      withProxy: initialData.withProxy,
      address: initialData.address
    }
  })

 
  useEffect(() => {
    const subscription = form.watch(values => {
      setNodes(nds =>
        nds.map(n => (n.id === node.id ? { ...n, data: values } : n))
      )
    })

    return () => subscription.unsubscribe()
  }, [form, node.id, setNodes])

  return (
    <Sheet>
      <SheetTrigger asChild>{children}</SheetTrigger>
      <SheetContent className="max-w-md">
        <SheetHeader>
          <SheetTitle className="text-lg font-semibold">
            Node ID: #{node.id}
          </SheetTitle>
          <SheetDescription className="text-sm text-foreground">
            Update the browser node configuration below.
          </SheetDescription>
        </SheetHeader>

        <Form {...form}>
          <form className="space-y-6 px-4 py-2">
            {/* Type Selector (Radio Group) */}
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem className="space-y-3">
                  <FormLabel>Type</FormLabel>
                  <FormControl>
                    <RadioGroup
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      className="flex space-x-6"
                    >
                      <FormItem className="flex items-center space-x-2">
                        <FormControl>
                          <RadioGroupItem value="local" />
                        </FormControl>
                        <FormLabel>Local</FormLabel>
                      </FormItem>
                      <FormItem className="flex items-center space-x-2">
                        <FormControl>
                          <RadioGroupItem value="remote" />
                        </FormControl>
                        <FormLabel>Remote</FormLabel>
                      </FormItem>
                    </RadioGroup>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Proxy Switch */}
            <FormField
              control={form.control}
              name="withProxy"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">Use Proxy</FormLabel>
                    <FormDescription>
                      Enable this to use a proxy for the browser.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />

            {/* Remote Address Input */}
            {form.watch('type') === 'remote' && (
              <FormField
                control={form.control}
                name="address"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Remote Address</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter remote address" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            )}
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  )
}
