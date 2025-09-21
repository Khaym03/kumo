// src/components/NavBar.tsx
import { useAppContext } from "@/context/app_ctx";
import RunButton from "@/infrastruture/run_kumo_btn";
import React from "react";
import { NavLink } from "react-router";

const NavBar: React.FC = () => {

  return (
    <>
      <NavLink
        to="/infrastructure"
        className={({ isActive }) =>
          isActive ? "font-medium text-primary" : "text-foreground"
        }
      >
        Infrastructure
      </NavLink>
      <NavLink
        to="/logs"
        className={({ isActive }) =>
            isActive ? "font-medium text-primary" : "text-foreground"
        }
      >
        Logs
      </NavLink>

       
    </>
  );
};

export default NavBar;
