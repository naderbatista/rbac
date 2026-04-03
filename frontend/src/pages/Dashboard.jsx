import { useAuth } from "../context/useAuth";

const typeLabels = { horario: "Horário", ip: "IP" };

export default function Dashboard() {
  const { user, permissions, policies } = useAuth();

  return (
    <div>
      <h2>Painel</h2>
      <p>Bem-vindo, <strong>{user?.username}</strong></p>

      <h3>Suas permissões</h3>
      {permissions.length === 0 ? (
        <p className="muted">Nenhuma permissão atribuída.</p>
      ) : (
        <ul>
          {permissions.map((p) => (
            <li key={p}><code>{p}</code></li>
          ))}
        </ul>
      )}

      <h3>Políticas ABAC ativas</h3>
      {policies.length === 0 ? (
        <p className="muted">Nenhuma política ABAC aplicada — acesso irrestrito por atributos.</p>
      ) : (
        <ul>
          {policies.map((p) => (
            <li key={p.id}>
              <strong>{p.name}</strong> ({typeLabels[p.type] || p.type}): <code>{p.value}</code>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
