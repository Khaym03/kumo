import type { BrowserConfig } from '@/infrastruture/node_factory'
import {
  useNodesState,
  type OnNodesChange,
  type Node,
  useEdgesState,
  type Edge,
  addEdge
} from '@xyflow/react'
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode
} from 'react'
import { initialEdges, initialNodes } from '@/infrastruture/constants'
import { useBuildConfig } from '@/infrastruture/useBuildConfig'
import type { LogEntry } from '@/logs/types'
import { EventsOff, EventsOn } from '@wailsjs/runtime/runtime'

interface AppContextType {
  isDarkMode: boolean
  setIsDarkMode: (value: boolean) => void

  nodes: Node[]
  setNodes: React.Dispatch<React.SetStateAction<Node[]>>
  onNodesChange: OnNodesChange<Node>

  edges: Edge[]
  setEdges: React.Dispatch<React.SetStateAction<Edge[]>>
  onEdgesChange: import('@xyflow/react').OnEdgesChange

  onConnect: (params: any) => void

  config: {
    root?: string | undefined
    browsers: any[]
  }
  isBuilding: boolean
  setIsBuilding: React.Dispatch<React.SetStateAction<boolean>>


  // logs
  logs: LogEntry[]
  setLogs: React.Dispatch<React.SetStateAction<LogEntry[]>>
}

// 2. Create the context with an initial value of `undefined`.
// This helps TypeScript ensure that the context is always used within a provider.
const AppContext = createContext<AppContextType | undefined>(undefined)

// 3. Create a provider component that will wrap the application.
export const AppProvider = ({ children }: { children: ReactNode }) => {
  const [isDarkMode, setIsDarkMode] = useState<boolean>(false)

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)

  const [logs, setLogs] = useState<LogEntry[]>([]);

  useEffect(() => {
    EventsOn('log_message', (data: LogEntry) => {
      setLogs((prevLogs) => [...prevLogs, data]);
    });
    return () => {
      EventsOff('log_message');
    };
  }, []);

  const onConnect = useCallback(
    (params: any) => setEdges(els => addEdge(params, els)),
    [setEdges]
  )
  const { config } = useBuildConfig(nodes, edges)

  const [isBuilding, setIsBuilding] = useState(false)

  // The value object holds all the states and functions to be shared.
  const value = {
    isDarkMode,
    setIsDarkMode,
    nodes,
    setNodes,
    onNodesChange,
    edges,
    setEdges,
    onEdgesChange,
    onConnect,
    config,
    isBuilding,
    setIsBuilding,
    logs,
    setLogs
  }

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>
}

// 4. Create a custom hook to easily consume the context.
export const useAppContext = () => {
  const context = useContext(AppContext)
  if (context === undefined) {
    throw new Error('useAppContext must be used within an AppProvider')
  }
  return context
}
