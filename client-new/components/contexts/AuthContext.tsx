"use client";

import React, { createContext, useState, useEffect, ReactNode } from "react";
import { decodeToken } from "@/lib/utils/token";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { DeleteToken, GetToken } from "@/lib/auth";

interface AuthContextType {
  isAuthenticated: boolean;
  privilege: number;
  logout: () => void;
}

export const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  privilege: 0,
  logout: () => {},
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [privilege, setPrivilege] = useState<number>(0);

  useEffect(() => {
    const token = GetToken();
    if (token) {
      const payload = JSON.parse(
        Buffer.from(token.split(".")[1], "base64").toString(),
      );
      console.log("Token saved: ", payload);
      setPrivilege(payload.privilege);
      setIsAuthenticated(true);
    } else {
      console.log("No token saved");
      setIsAuthenticated(false);
      setPrivilege(0);
    }
  }, []);

  const logout = () => {
    DeleteToken();
    setPrivilege(0);
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, privilege, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
