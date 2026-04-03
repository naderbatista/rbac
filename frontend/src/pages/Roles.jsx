import { useState, useEffect } from "react";
import { api } from "../services/api";
import { useAuth } from "../context/useAuth";

export default function Roles() {
  const { hasPerm } = useAuth();
  const [roles, setRoles] = useState([]);
  const [permissions, setPermissions] = useState([]);
  const [policies, setPolicies] = useState([]);
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  async function load() {
    try {
      const [r, p, pol] = await Promise.all([api.listRoles(), api.listPermissions(), api.listPolicies()]);
      setRoles(r);
      setPermissions(p);
      setPolicies(pol);
    } catch (err) {
      setError(err.message);
    }
  }

  useEffect(() => {
    let cancel = false;
    Promise.all([api.listRoles(), api.listPermissions(), api.listPolicies()])
      .then(([r, p, pol]) => { if (!cancel) { setRoles(r); setPermissions(p); setPolicies(pol); } })
      .catch((err) => { if (!cancel) setError(err.message); });
    return () => { cancel = true; };
  }, []);

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

  async function togglePolicy(roleId, policyId) {
    const role = roles.find((r) => r.id === roleId);
    const current = role?.policies || [];
    const next = current.includes(policyId)
      ? current.filter((p) => p !== policyId)
      : [...current, policyId];
    try {
      await api.assignPolicies(roleId, next);
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  function policyName(id) {
    return policies.find((p) => p.id === id)?.name || id;
  }

  return (
    <div>
      <h2>Perfis</h2>
      {error && <p className="error">{error}</p>}

      {hasPerm("role:write") && (
        <form onSubmit={handleCreate} className="inline-form">
          <input placeholder="Nome do perfil" value={name} onChange={(e) => setName(e.target.value)} required />
          <button type="submit">Criar perfil</button>
        </form>
      )}

      <table>
        <thead>
          <tr>
            <th>Perfil</th>
            <th>Permissões</th>
            <th>Políticas ABAC</th>
            {hasPerm("role:write") && <th>Alternar permissão</th>}
            {hasPerm("role:write") && <th>Alternar política</th>}
          </tr>
        </thead>
        <tbody>
          {roles.map((r) => (
            <tr key={r.id}>
              <td>{r.name}</td>
              <td>{(r.permissions || []).map(permName).join(", ") || "—"}</td>
              <td>{(r.policies || []).map(policyName).join(", ") || "—"}</td>
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
              {hasPerm("role:write") && (
                <td>
                  {policies.map((p) => (
                    <button
                      key={p.id}
                      className={(r.policies || []).includes(p.id) ? "btn-active" : "btn-outline"}
                      onClick={() => togglePolicy(r.id, p.id)}
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
