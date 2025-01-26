import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { FaHome, FaSearch, FaUserCircle } from "react-icons/fa";
import { useAuth } from "../context/AuthContext";

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth(); // Get user from AuthProvider

  const isLogin: boolean = user !== null;
  const isAdmin: boolean = isLogin && user!.role === "admin";

  const handleLogOut = () => {
    logout(); // Use logout from AuthProvider
    navigate("/");
  };

  return (
    <nav>
      {/* Search Bar */}
      <div className="flex items-center border rounded-full px-4 py-2 w-96 ml-8">
        <FaSearch className="text-gray-400" />
        <input
          type="text"
          placeholder="Search threads..."
          className="bg-transparent outline-none w-full ml-2"
        />
      </div>

      {/* User Actions */}
      <div className="flex items-center mr-8 space-x-6">
        <Link to="/" className="text-gray-600 hover:text-black">
          <FaHome className="text-3xl text-gray-600" /> Home
        </Link>

        {/* Admin Actions */}
        {isAdmin && (
          <Link to="/admin" className="text-gray-600 hover:text-black">
            Admin
          </Link>
        )}

        {/* Login user actions */}
        {isLogin ? (
          <div>
            <Link
              to="/"
              className="text-gray-600 hover:text-black"
              onClick={handleLogOut}
            >
              Logout
            </Link>
            <Link
              to={`/user/${user!.id}`}
              className="text-gray-600 hover:text-black"
            >
              <FaUserCircle className="text-3xl text-gray-600" />
              {user?.username} {/* Display username */}
            </Link>
          </div>
        ) : (
          <div>
            <Link to="/login" className="text-gray-600 hover:text-black">
              Login
            </Link>
            <Link to="/signup" className="text-gray-600 hover:text-black">
              Sign Up
            </Link>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
