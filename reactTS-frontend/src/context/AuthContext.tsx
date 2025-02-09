import React, { createContext, useState, useContext, useEffect } from "react";
import { User } from "../models/models";
// import BASE_URL from "../config";
const BASE_URL = import.meta.env.VITE_API_URL;

interface AuthContextType {
  user: User | null;
  login: (token: string) => Promise<void>;
  logout: () => void;
  fetchUserData: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);

  // Fetch user data using the token
  const fetchUserData = async () => {
    const token = localStorage.getItem("token");
    if (token) {
      try {
        const response = await fetch(`${BASE_URL}/api/me`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        if (response.ok) {
          const userData = await response.json();
          setUser(userData);
        }
      } catch (error) {
        console.error("Error fetching user data:", error);
      }
    }
  };

  // Login functionality
  const login = async (token: string) => {
    localStorage.setItem("token", token);
    await fetchUserData();
  };

  // Logout functionality
  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
  };

  // Fetch user data on mount
  useEffect(() => {
    fetchUserData();
  }, []);

  return (
    <AuthContext.Provider value={{ user, login, logout, fetchUserData }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
