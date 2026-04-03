import { NavLink } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export default function Nav() {
  const { user, logout } = useAuth();

  return (
    <nav className="topnav">
      <div className="nav-links">
        <NavLink to="/">Dashboard</NavLink>
        <NavLink to="/users">Users</NavLink>
        <NavLink to="/roles">Roles</NavLink>
        <NavLink to="/permissions">Permissions</NavLink>
      </div>
      <div className="nav-right">
        <span>{user?.username}</span>
        <button onClick={logout}>Logout</button>
      </div>
    </nav>
  );
}
