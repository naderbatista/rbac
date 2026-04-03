import { useAuth } from "../context/AuthContext";

export default function Dashboard() {
  const { user, permissions } = useAuth();

  return (
    <div>
      <h2>Dashboard</h2>
      <p>Welcome, <strong>{user?.username}</strong></p>
      <h3>Your permissions</h3>
      {permissions.length === 0 ? (
        <p className="muted">No permissions assigned.</p>
      ) : (
        <ul>
          {permissions.map((p) => (
            <li key={p}><code>{p}</code></li>
          ))}
        </ul>
      )}
    </div>
  );
}
