import { useState, useEffect } from "react";
import { api } from "../services/api";
import { AuthContext } from "./authContext";

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [permissions, setPermissions] = useState([]);
  const [policies, setPolicies] = useState([]);
  const [loading, setLoading] = useState(() => !!localStorage.getItem("token"));

  useEffect(() => {
    if (!loading) return;
    api.me()
      .then((data) => {
        setUser(data.user);
        setPermissions(data.permissions || []);
        setPolicies(data.policies || []);
      })
      .catch(() => localStorage.removeItem("token"))
      .finally(() => setLoading(false));
  }, [loading]);

  function login(token, userData, perms, pols) {
    localStorage.setItem("token", token);
    setUser(userData);
    setPermissions(perms || []);
    setPolicies(pols || []);
  }

  function logout() {
    localStorage.removeItem("token");
    setUser(null);
    setPermissions([]);
    setPolicies([]);
  }

  function hasPerm(perm) {
    return permissions.includes(perm);
  }

  if (loading) return null;

  return (
    <AuthContext.Provider value={{ user, permissions, policies, login, logout, hasPerm }}>
      {children}
    </AuthContext.Provider>
  );
}
