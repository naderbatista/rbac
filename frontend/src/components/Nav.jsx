import { NavLink } from "react-router-dom";
import { useAuth } from "../context/useAuth";

export default function Nav() {
  const { user, logout } = useAuth();

  return (
    <nav className="topnav">
      <div className="nav-links">
        <NavLink to="/">Painel</NavLink>
        <NavLink to="/users">Usuários</NavLink>
        <NavLink to="/roles">Perfis</NavLink>
        <NavLink to="/permissions">Permissões</NavLink>
        <NavLink to="/policies">Políticas</NavLink>
      </div>
      <div className="nav-right">
        <span>{user?.username}</span>
        <button onClick={logout}>Sair</button>
      </div>
    </nav>
  );
}
