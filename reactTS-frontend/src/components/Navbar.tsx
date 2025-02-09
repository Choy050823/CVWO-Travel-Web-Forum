import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { FaHome, FaSearch, FaTripadvisor, FaUserCircle } from "react-icons/fa";
import { useAuth } from "../context/AuthContext";

const Navbar: React.FC = () => {
  const navigate = useNavigate();
  const { user, logout } = useAuth();

  const isLogin: boolean = user !== null;
  // const isAdmin: boolean = isLogin && user!.role === "admin";

  const handleLogOut = () => {
    logout();
    navigate("/");
  };

  return (
    <nav className="border-b border-gray-200 bg-white p-4 flex items-center justify-between relative shadow-sm">
      <div className="flex items-center gap-3 ml-4">
        {/* <img
          src="/logo.svg" // Replace with your actual logo path
          alt="WanderVerse Logo"
          className="h-8 w-8 rounded-full"
        /> */}
        <FaTripadvisor className={`h-8 w-8 rounded-full`} />
        <Link
          to="/"
          className="text-2xl font-bold bg-gradient-to-r from-blue-600 to-green-500 bg-clip-text text-transparent"
        >
          WanderVerse
        </Link>
      </div>

      {/* Centered Search Bar */}
      <div className="absolute left-1/2 transform -translate-x-1/2">
        <div className="flex items-center border border-gray-200 rounded-full px-4 py-2 w-96 bg-gray-50 focus-within:border-blue-500 focus-within:bg-white transition-colors">
          <FaSearch className={`text-gray-400`} />
          <input
            type="text"
            placeholder="Search threads..."
            className="bg-transparent outline-none w-full ml-2 placeholder-gray-400 text-gray-800"
          />
        </div>
      </div>

      {/* User Actions */}
      <div className="flex items-center space-x-6 w-[200px] justify-end">
        <Link
          to="/"
          className="text-gray-500 hover:text-blue-600 flex items-center gap-1 transition-colors"
        >
          <FaHome className={`text-xl`} />
          <span className="hidden sm:inline text-sm font-medium">Home</span>
        </Link>

        {isLogin ? (
          <div className="flex items-center space-x-4">
            <button
              onClick={handleLogOut}
              className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-full text-sm font-medium transition-colors"
            >
              Logout
            </button>
            <Link
              to={`/user/${user!.id}`}
              className="flex items-center gap-2 text-gray-600 hover:text-blue-600 group transition-colors"
            >
              <FaUserCircle className={`text-xl text-gray-500`} />
              <span className="hidden sm:inline max-w-[120px] truncate text-sm font-medium">
                {user?.username}
              </span>
            </Link>
          </div>
        ) : (
          <div className="flex items-center gap-3">
            <Link
              to="/login"
              className="text-gray-600 hover:text-blue-600 px-4 py-2 text-sm font-medium transition-colors"
            >
              Login
            </Link>
            <Link
              to="/signup"
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-full text-sm font-medium transition-colors whitespace-nowrap"
            >
              Sign Up
            </Link>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
