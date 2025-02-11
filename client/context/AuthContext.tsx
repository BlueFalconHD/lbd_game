import { createContext, useState, useEffect, ReactNode } from "react";
import Cookies from "js-cookie";

interface AuthContextProps {
  isAuthenticated: boolean;
  isAdmin: boolean;
  setAuth: (auth: boolean, admin: boolean) => void;
}

export const AuthContext = createContext<AuthContextProps>({
  isAuthenticated: false,
  isAdmin: false,
  setAuth: () => {},
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [isAdmin, setIsAdmin] = useState<boolean>(false);

  useEffect(() => {
    const token = Cookies.get("token");
    if (token) {
      // Decode token to check if user is admin
      // If token is valid, set isAuthenticated to true
      const payload = JSON.parse(
        Buffer.from(token.split(".")[1], "base64").toString(),
      );
      setIsAdmin(payload.is_admin);
      setIsAuthenticated(true);
    } else {
      setIsAuthenticated(false);
      setIsAdmin(false);
    }
  }, []);

  const setAuth = (auth: boolean, admin: boolean) => {
    setIsAuthenticated(auth);
    setIsAdmin(admin);
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, isAdmin, setAuth }}>
      {children}
    </AuthContext.Provider>
  );
};
