"use client";

import React, { createContext, useState, useEffect, ReactNode } from "react";
import { decodeToken } from "@/lib/utils/token";
import { getCookie, deleteCookie } from "@/lib/utils/cookies";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";

interface AuthContextType {
  isAuthenticated: boolean;
  privilege: number;
  logout: () => void;
  processLogin: () => void;
}

export const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  privilege: 0,
  logout: () => {},
  processLogin: () => {},
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [privilege, setPrivilege] = useState<number>(0);

  useEffect(() => {
    const token = Cookies.get("token");
    if (token) {
      const payload = JSON.parse(
        Buffer.from(token.split(".")[1], "base64").toString(),
      );
      console.log(payload);
      setPrivilege(payload.privilege);
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
      setPrivilege(0);
    }
  }, []);

  const logout = () => {
    deleteCookie("token");
    setPrivilege(0);
    setIsAuthenticated(false);
  };

  const processLogin = () => {
    const token = Cookies.get("token");
    if (token) {
      const payload = JSON.parse(
        Buffer.from(token.split(".")[1], "base64").toString(),
      );
      setPrivilege(payload.privilege);
      setIsAuthenticated(true);
    }
  };

  return (
    <AuthContext.Provider
      value={{ isAuthenticated, privilege, logout, processLogin }}
    >
      {children}
    </AuthContext.Provider>
  );
};
