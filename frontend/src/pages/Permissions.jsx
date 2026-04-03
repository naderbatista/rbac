import { useState, useEffect } from "react";
import { api } from "../services/api";
import { useAuth } from "../context/AuthContext";

export default function Permissions() {
  const { hasPerm } = useAuth();
  const [permissions, setPermissions] = useState([]);
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    load();
  }, []);

  async function load() {
    try {
      setPermissions(await api.listPermissions());
    } catch (err) {
      setError(err.message);
    }
  }

  async function handleCreate(e) {
    e.preventDefault();
    setError("");
    try {
      await api.createPermission(name);
      setName("");
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div>
      <h2>Permissions</h2>
      {error && <p className="error">{error}</p>}

      {hasPerm("role:write") && (
        <form onSubmit={handleCreate} className="inline-form">
          <input placeholder="e.g. post:delete" value={name} onChange={(e) => setName(e.target.value)} required />
          <button type="submit">Create permission</button>
        </form>
      )}

      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>ID</th>
          </tr>
        </thead>
        <tbody>
          {permissions.map((p) => (
            <tr key={p.id}>
              <td><code>{p.name}</code></td>
              <td className="muted">{p.id}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
