// src/App.tsx
import React from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router";
import TitleBar from "./components/tittle_bar";
import Infrastructure from "./infrastruture/Infrastructure";
import Logs from "@/logs/logs"; // Assuming you have a Logs component
import { AppProvider } from "./context/app_ctx";

const App: React.FC = () => {
  return (
    <AppProvider>
      <BrowserRouter>
        <div className="flex flex-col h-screen">
          {/* The custom title bar goes at the top */}
          <TitleBar />

          {/* The main content area */}
          <main className="flex-1 overflow-auto">
            <Routes>
              {/* Redirect root "/" to "/infrastructure" */}
              <Route path="/" element={<Navigate to="/infrastructure" replace />} />

              {/* Infrastructure route */}
              <Route path="/infrastructure" element={<Infrastructure />} />

              {/* Logs route */}
              <Route path="/logs" element={<Logs />} />
            </Routes>
          </main>
        </div>
      </BrowserRouter>
    </AppProvider>
  );
};

export default App;
