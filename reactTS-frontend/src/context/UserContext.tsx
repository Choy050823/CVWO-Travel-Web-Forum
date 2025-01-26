import React, { createContext, useContext } from "react";
import { User } from "../models/models";
import BASE_URL from "../config";

interface UserContextType {
  fetchUserDetails: (userId: number) => Promise<User | null>; // Return user details directly
}

const UserContext = createContext<UserContextType | undefined>(undefined);

export const UserProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const userCache = new Map<number, User>(); // Cache for user details

  // Fetch user details by ID
  const fetchUserDetails = async (userId: number): Promise<User | null> => {
    if (!userId) {
      console.warn(
        "fetchUserDetails called with an undefined or invalid userId"
      );
      return null;
    }

    try {
      // Check cache first
      if (userCache.has(userId)) {
        console.log(`Returning cached user for userId: ${userId}`);
        return userCache.get(userId)!;
      }

      console.log(`Fetching user details from API for userId: ${userId}`);
      const token = localStorage.getItem("token");
      const url = token
        ? `${BASE_URL}/api/users/${userId}`
        : `${BASE_URL}/api/public/users/${userId}`;

      const response = await fetch(url, {
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch user details for userId: ${userId}`);
      }

      const data = await response.json();
      userCache.set(userId, data); // Cache the response
      return data;
    } catch (error) {
      console.error("Error fetching user details:", error);
      return null; // Return null on error
    }
  };

  return (
    <UserContext.Provider value={{ fetchUserDetails }}>
      {children}
    </UserContext.Provider>
  );
};

export const useUser = () => {
  const context = useContext(UserContext);
  if (!context) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
};
