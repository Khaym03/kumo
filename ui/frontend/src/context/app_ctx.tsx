
import  { createContext, useContext, useState,type ReactNode } from 'react';

interface AppContextType {
  isDarkMode: boolean;
  setIsDarkMode: (value: boolean) => void;
}

// 2. Create the context with an initial value of `undefined`.
// This helps TypeScript ensure that the context is always used within a provider.
const AppContext = createContext<AppContextType | undefined>(undefined);

// 3. Create a provider component that will wrap the application.
export const AppProvider = ({ children }: { children: ReactNode }) => {
  const [isDarkMode, setIsDarkMode] = useState<boolean>(false);

  // The value object holds all the states and functions to be shared.
  const value = {
    isDarkMode,
    setIsDarkMode,
  };

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
};

// 4. Create a custom hook to easily consume the context.
export const useAppContext = () => {
  const context = useContext(AppContext);
  if (context === undefined) {
    throw new Error('useAppContext must be used within an AppProvider');
  }
  return context;
};