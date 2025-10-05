"use client";

import {
  createContext,
  useContext,
  useState,
  ReactNode,
  useEffect,
} from "react";
import { useRouter } from "next/navigation";
import { User } from "@/types/types";
import { loginRequest } from "@/lib/login";
import { sha256Hex } from "@/lib/crypto";

type AuthContextType = {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
};
// Sometimes I wish that class based React components weren't so outdated bc this way of doing context shit is fucking disgusting. ~brtcrt
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  // I spent an embarasingly long time writing a custom hook that we probably don't need since
  // we really should be using the amazong thingy for authentication,
  // although this might still be useful in that case. ~brtcrt
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  // Initialize auth state
  useEffect(() => {
    try {
      const storedUser = localStorage.getItem("user");
      if (storedUser) {
        const parsed: User = JSON.parse(storedUser);
        setUser(parsed);
      }
    } catch (err) {
      console.error("Failed to parse stored user", err);
      localStorage.removeItem("user");
    } finally {
      setLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    try {
      const hashedPassword = await sha256Hex(password);
      const res = await loginRequest({
        email: email,
        password: hashedPassword,
      });
      const token = res.token;
      const userData = { email, token };
      setUser(userData);
      localStorage.setItem("user", JSON.stringify(userData));
      router.push("/");
      return;
    } catch (err) {
      throw err;
    }
  };

  // delete from local storage on logout
  const logout = () => {
    setUser(null);
    localStorage.removeItem("user");
    router.push("/login");
  };

  if (loading) {
    return null;
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout }}>
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
