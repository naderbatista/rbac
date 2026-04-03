import { useState, useEffect } from "react";
import { api } from "../services/api";
import { useAuth } from "../context/useAuth";

export default function Permissions() {
  const { hasPerm } = useAuth();
  const [permissions, setPermissions] = useState([]);
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  async function load() {
    try {
      setPermissions(await api.listPermissions());
    } catch (err) {
      setError(err.message);
    }
  }

  useEffect(() => {
    let cancel = false;
    api.listPermissions()
      .then((p) => { if (!cancel) setPermissions(p); })
      .catch((err) => { if (!cancel) setError(err.message); });
    return () => { cancel = true; };
  }, []);

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
      <h2>Permissões</h2>
      {error && <p className="error">{error}</p>}

      {hasPerm("role:write") && (
        <form onSubmit={handleCreate} className="inline-form">
          <input placeholder="ex: post:delete" value={name} onChange={(e) => setName(e.target.value)} required />
          <button type="submit">Criar permissão</button>
        </form>
      )}

      <table>
        <thead>
          <tr>
            <th>Nome</th>
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
