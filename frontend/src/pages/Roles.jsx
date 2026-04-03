import { useState, useEffect } from "react";
import { api } from "../services/api";
import { useAuth } from "../context/AuthContext";

export default function Roles() {
  const { hasPerm } = useAuth();
  const [roles, setRoles] = useState([]);
  const [permissions, setPermissions] = useState([]);
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    load();
  }, []);

  async function load() {
    try {
      const [r, p] = await Promise.all([api.listRoles(), api.listPermissions()]);
      setRoles(r);
      setPermissions(p);
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleCreate(e) {
    e.preventDefault();
    setError("");
    try {
      await api.createRole(name);
      setName("");
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  async function togglePermission(roleId, permId) {
    const role = roles.find((r) => r.id === roleId);
    const current = role?.permissions || [];
    const next = current.includes(permId)
      ? current.filter((p) => p !== permId)
      : [...current, permId];
    try {
      await api.assignPermissions(roleId, next);
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  function permName(id) {
    return permissions.find((p) => p.id === id)?.name || id;
  }

  return (
    <div>
      <h2>Roles</h2>
      {error && <p className="error">{error}</p>}

      {hasPerm("role:write") && (
        <form onSubmit={handleCreate} className="inline-form">
          <input placeholder="Role name" value={name} onChange={(e) => setName(e.target.value)} required />
          <button type="submit">Create role</button>
        </form>
      )}

      <table>
        <thead>
          <tr>
            <th>Role</th>
            <th>Permissions</th>
            {hasPerm("role:write") && <th>Toggle permission</th>}
          </tr>
        </thead>
        <tbody>
          {roles.map((r) => (
            <tr key={r.id}>
              <td>{r.name}</td>
              <td>{(r.permissions || []).map(permName).join(", ") || "—"}</td>
              {hasPerm("role:write") && (
                <td>
                  {permissions.map((p) => (
                    <button
                      key={p.id}
                      className={(r.permissions || []).includes(p.id) ? "btn-active" : "btn-outline"}
                      onClick={() => togglePermission(r.id, p.id)}
                    >
                      {p.name}
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
