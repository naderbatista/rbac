import { createContext, useContext, useState, useEffect } from "react";
import { api } from "../services/api";

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [permissions, setPermissions] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) {
      setLoading(false);
      return;
    }
    api.me()
      .then((data) => {
        setUser(data.user);
        setPermissions(data.permissions || []);
      })
      .catch(() => localStorage.removeItem("token"))
      .finally(() => setLoading(false));
  }, []);

  function login(token, userData, perms) {
    localStorage.setItem("token", token);
    setUser(userData);
    setPermissions(perms || []);
  }

  function logout() {
    localStorage.removeItem("token");
    setUser(null);
    setPermissions([]);
  }

  function hasPerm(perm) {
    return permissions.includes(perm);
  }

  if (loading) return null;

  return (
    <AuthContext.Provider value={{ user, permissions, login, logout, hasPerm }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be inside AuthProvider");
  return ctx;
}
