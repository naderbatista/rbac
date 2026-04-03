import { useState, useEffect } from "react";
import { api } from "../services/api";
import { useAuth } from "../context/useAuth";

export default function Policies() {
  const { hasPerm } = useAuth();
  const [policies, setPolicies] = useState([]);
  const [name, setName] = useState("");
  const [type, setType] = useState("horario");
  const [value, setValue] = useState("");
  const [error, setError] = useState("");

  async function load() {
    try {
      setPolicies(await api.listPolicies());
    } catch (err) {
      setError(err.message);
    }
  }

  useEffect(() => {
    let cancel = false;
    api.listPolicies()
      .then((p) => { if (!cancel) setPolicies(p); })
      .catch((err) => { if (!cancel) setError(err.message); });
    return () => { cancel = true; };
  }, []);

  async function handleCreate(e) {
    e.preventDefault();
    setError("");
    try {
      await api.createPolicy(name, type, value);
      setName("");
      setValue("");
      load();
    } catch (err) {
      setError(err.message);
    }
  }

  const typeLabels = { horario: "Horário", ip: "IP" };

  return (
    <div>
      <h2>Políticas ABAC</h2>
      {error && <p className="error">{error}</p>}

      {hasPerm("role:write") && (
        <form onSubmit={handleCreate} className="inline-form">
          <input placeholder="Nome da política" value={name} onChange={(e) => setName(e.target.value)} required />
          <select value={type} onChange={(e) => setType(e.target.value)}>
            <option value="horario">Horário (ex: 08:00-18:00)</option>
            <option value="ip">IP (ex: 127.0.0.1,::1)</option>
          </select>
          <input placeholder={type === "horario" ? "08:00-18:00" : "127.0.0.1,::1"} value={value} onChange={(e) => setValue(e.target.value)} required />
          <button type="submit">Criar política</button>
        </form>
      )}

      <table>
        <thead>
          <tr>
            <th>Nome</th>
            <th>Tipo</th>
            <th>Valor</th>
            <th>ID</th>
          </tr>
        </thead>
        <tbody>
          {policies.map((p) => (
            <tr key={p.id}>
              <td>{p.name}</td>
              <td>{typeLabels[p.type] || p.type}</td>
              <td><code>{p.value}</code></td>
              <td className="muted">{p.id}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
