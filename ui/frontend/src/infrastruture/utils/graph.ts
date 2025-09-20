import type { Edge, Node } from '@xyflow/react';
import type { BrowserConfig } from '@/infrastruture/node_factory';

export function buildConfigFromGraph(
  nodes: Node[],
  edges: Edge[]
): { root?: string; browsers: (BrowserConfig & { pages: Node[] })[] } {
  // Build adjacency list: nodeId -> array of target nodeIds (children)
  const adjacencyList: Record<string, string[]> = {};
  edges.forEach(edge => {
    if (!adjacencyList[edge.source]) adjacencyList[edge.source] = [];
    adjacencyList[edge.source].push(edge.target);
  });

  // Find root node(s) as those without incoming edges
  const allTargets = new Set(edges.map(e => e.target));
  const rootIds = nodes.filter(n => !allTargets.has(n.id)).map(n => n.id);
  const rootId = rootIds[0]; // Take first root or adapt logic to multiple roots

  // Recursive function to build BrowserConfig + nested pages (full Nodes)
  function buildBrowserConfig(nodeId: string): BrowserConfig & { pages: Node[] } {
    const node = nodes.find(n => n.id === nodeId);
    if (!node) throw new Error(`Node ${nodeId} not found`);

    // cast node.data to BrowserConfig (adjust if different)
    const data = node.data as unknown as BrowserConfig;

    // Get children node IDs from adjacency list
    const childrenIds = adjacencyList[nodeId] || [];

    // Map children node IDs to full Node objects
    const childrenNodes = childrenIds
      .map(childId => nodes.find(n => n.id === childId))
      .filter(Boolean) as Node[];

    return {
      type: data.type,
      withProxy: data.withProxy,
      address: data.address,
      pages: childrenNodes,
    };
  }

  // Build browsers list from root's children
  const browsers = (adjacencyList[rootId] || [])
    .map(id => buildBrowserConfig(id))
    .filter(Boolean);

  return {
    root: rootId,
    browsers,
  };
}
