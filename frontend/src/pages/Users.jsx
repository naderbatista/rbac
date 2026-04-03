import { useState, useEffect } from "react";
import { api } from "../services/api";
import { useAuth } from "../context/useAuth";

export default function Users() {
  const { hasPerm } = useAuth();
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function load() {
    try {
      const [u, r] = await Promise.all([api.listUsers(), api.listRoles()]);
      setUsers(u);
      setRoles(r);
    } catch (err) {
      setError(err.message);
    }
  }

  useEffect(() => {
    let cancel = false;
    Promise.all([api.listUsers(), api.listRoles()])
      .then(([u, r]) => { if (!cancel) { setUsers(u); setRoles(r); } })
      .catch((err) => { if (!cancel) setError(err.message); });
    return () => { cancel = true; };
  }, []);

  async function handleCreate(e) {
    e.preventDefault();
    setError("");
    try {
      await api.createUser(username, password);
      setUsername("");
      setPassword("");
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleAssignRole(userId, roleId) {
    const user = users.find((u) => u.id === userId);
    const currentRoles = user?.roles || [];
    const next = currentRoles.includes(roleId)
      ? currentRoles.filter((r) => r !== roleId)
      : [...currentRoles, roleId];
    try {
      await api.assignRoles(userId, next);
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  function roleName(id) {
    return roles.find((r) => r.id === id)?.name || id;
  }

  return (
    <div>
      <h2>Usuários</h2>
      {error && <p className="error">{error}</p>}

      {hasPerm("user:write") && (
        <form onSubmit={handleCreate} className="inline-form">
          <input placeholder="Usuário" value={username} onChange={(e) => setUsername(e.target.value)} required />
          <input type="password" placeholder="Senha" value={password} onChange={(e) => setPassword(e.target.value)} required />
          <button type="submit">Criar usuário</button>
        </form>
      )}

      <table>
        <thead>
          <tr>
            <th>Usuário</th>
            <th>Perfis</th>
            {hasPerm("role:write") && <th>Alternar perfil</th>}
          </tr>
        </thead>
        <tbody>
          {users.map((u) => (
            <tr key={u.id}>
              <td>{u.username}</td>
              <td>{(u.roles || []).map(roleName).join(", ") || "—"}</td>
              {hasPerm("role:write") && (
                <td>
                  {roles.map((r) => (
                    <button
                      key={r.id}
                      className={(u.roles || []).includes(r.id) ? "btn-active" : "btn-outline"}
                      onClick={() => handleAssignRole(u.id, r.id)}
                    >
                      {r.name}
                    </button>
                  ))}
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
