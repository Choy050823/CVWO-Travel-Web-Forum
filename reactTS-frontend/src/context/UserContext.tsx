import React, { createContext, useState, useContext } from "react";
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
    try {
      // Check cache first
      if (userCache.has(userId)) {
        return userCache.get(userId)!;
      }

      const token = localStorage.getItem("token");

      // Use the public endpoint if the user is not authenticated
      const url = token
        ? `${BASE_URL}/api/users/${userId}` // Authenticated endpoint
        : `${BASE_URL}/api/public/users/${userId}`; // Public endpoint

      const response = await fetch(url, {
        headers: token
          ? { Authorization: `Bearer ${token}` } // Include token if authenticated
          : {}, // No headers for public endpoint
      });

      if (!response.ok) {
        throw new Error("Failed to fetch user details");
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
