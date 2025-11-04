"use client";
import { createContext, useContext, useEffect, useState } from "react";
import { fetch } from "@/lib/instance";

const AuthContext = createContext({
  user: null,
  loading: true,
});

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      if (user !== null) return;
      try {
        const res = await fetch.post("/api/v1/auth/verify");
        if (res.data.auth) setUser(res.data.user);
      } catch {
        setUser(null);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
