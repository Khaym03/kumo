// src/constants.ts

import { Position, type Node, type Edge } from '@xyflow/react';
import { BrowserNode } from './node/browser_node';
import { RootNode } from './node/root';
import { PageNode } from './node/page_node';

export const initialNodes: Node[] = [
  {
    id: 'root',
    data: { label: 'Root' },
    position: { x: 0, y: 0 },
    sourcePosition: Position.Bottom,
    type: 'root',
  },
];

export const initialEdges: Edge[] = [];

export const NodeTypes = {
  browser: BrowserNode,
  root: RootNode,
  page: PageNode,
};

export const EdgeTypes = {};

export type NodeComponentType = keyof typeof NodeTypes;

export const SupportedNodeComponentTypes = Object.keys(
  NodeTypes,
) as NodeComponentType[];