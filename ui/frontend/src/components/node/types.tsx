import { BrowserNode } from './browser_node'
import { RootNode } from './root'
import { PageNode } from './page_node'

const nodeTypes = {
  browser: BrowserNode,
  root: RootNode,
  page: PageNode
}

const edgeTypes = {}

type NodeComponentTypes = keyof typeof nodeTypes

const NodeComponents = [...Object.keys(nodeTypes)] as  NodeComponentTypes[]

export { nodeTypes, edgeTypes, NodeComponents,type NodeComponentTypes }
