"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
  useCallback,
} from "react";
import { User, initializeApi, getMe } from "~/lib/api";

interface AuthContextType {
  accessToken: string | null;
  user: User | null;
  isLoggedIn: boolean;
  login: (token: string, user: User) => void;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const logout = useCallback(() => {
    setAccessToken(null);
    setUser(null);
    localStorage.removeItem("user");
  }, []);

  const login = useCallback((token: string, user: User) => {
    setAccessToken(token);
    setUser(user);
    localStorage.setItem("user", JSON.stringify(user));
  }, []);

  useEffect(() => {
    initializeApi(
      () => accessToken,
      (token) => setAccessToken(token),
      logout
    );
  }, [accessToken, logout]);

  useEffect(() => {
    const authenticateUser = async () => {
      try {
        const refreshResponse = await fetch("/api/auth/refresh-proxy", {
          method: "POST",
        });

        if (!refreshResponse.ok) {
          throw new Error("No active session");
        }

        const { access_token } = await refreshResponse.json();

        initializeApi(
          () => access_token,
          () => {},
          logout
        );

        const currentUser = await getMe();

        login(access_token, currentUser);
      } catch (error) {
        logout();
      } finally {
        setIsLoading(false);
      }
    };

    authenticateUser();
  }, [login, logout]);

  return (
    <AuthContext.Provider
      value={{
        accessToken,
        user,
        isLoggedIn: !!user && !!accessToken,
        login,
        logout,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
