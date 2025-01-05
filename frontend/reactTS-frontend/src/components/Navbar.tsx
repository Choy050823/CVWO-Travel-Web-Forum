import React from "react";
import { useNavigate } from "react-router-dom";
import { User } from "../models/models";

const Navbar: React.FC<{ user: User | null }> = ({ user }) => {
  const navigate = useNavigate();
  const isLogin = user !== null;
  const isAdmin: boolean = isLogin && user!.role === "admin";

  const handleLogOut = () => {
    // Handle Log Out here
  };

  return (
    <div>
      <ul>
        <li>
          <button onClick={() => navigate("/")}>Home</button>
        </li>
        <li>
          <button onClick={() => navigate("/search")}>
            Browse specific threads!
          </button>
        </li>
        <li>
          {isAdmin && <button onClick={() => navigate("/admin")}>Admin</button>}
        </li>
        {isLogin ? (
          <div>
            <li>
              <button onClick={handleLogOut}>Log Out</button>
            </li>
            <li>
              <button onClick={() => navigate("/notifications")}>
                Notifications
              </button>
            </li>
            <li>
              <button onClick={() => navigate(`/user/:${user.id}`)}>
                User Profile
                <img
                  src={user.profilePicture === null ? "" : user.profilePicture!}
                  alt=""
                />
              </button>
            </li>
          </div>
        ) : (
          <div>
            <li>
              <button onClick={() => navigate("/login")}>Login</button>
            </li>
            <li>
              <button onClick={() => navigate("/signup")}>Sign Up</button>
            </li>
          </div>
        )}
      </ul>
    </div>
  );
};

export default Navbar;
