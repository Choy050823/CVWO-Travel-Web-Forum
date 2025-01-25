import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { User } from "../models/models";

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const storedUser = localStorage.getItem("user");
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
  }, []);

  const isLogin: boolean = user !== null;
  const isAdmin: boolean = isLogin && user!.role === "admin";

  const handleLogOut = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    localStorage.removeItem("name");
    setUser(null);
    navigate("/");
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
        {isAdmin && (
          <li>
            <button onClick={() => navigate("/admin")}>Admin</button>
          </li>
        )}
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
              <button onClick={() => navigate(`/user/:${user!.id}`)}>
                User Profile
                <img
                  src={
                    user!.profilePicture === null ? "" : user!.profilePicture!
                  }
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
