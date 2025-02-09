import React, { createContext, useState, useContext, useEffect } from "react";
import { Thread } from "../models/models";
// import BASE_URL from "../config";

const BASE_URL = import.meta.env.VITE_API_URL;

interface ThreadContextType {
  threads: Thread[];
  fetchThreads: () => Promise<void>;
  createThread: (thread: Omit<Thread, "id">) => Promise<void>;
  updateThread: (thread: Thread) => Promise<void>;
  deleteThread: (threadId: number) => Promise<void>;
  likeThread: (threadId: number) => Promise<void>;
}

const ThreadContext = createContext<ThreadContextType | undefined>(undefined);

export const ThreadProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [threads, setThreads] = useState<Thread[]>([]);

  // Fetch all threads
  const fetchThreads = async () => {
    try {
      const response = await fetch(`${BASE_URL}/api/threads`);

      if (!response.ok) {
        throw new Error("Failed to fetch threads");
      }

      const data = await response.json();
      console.log(data);

      // Transform the data to match the frontend Thread model
      const transformedThreads = data.map((thread: any) => ({
        id: thread.id,
        title: thread.title,
        content: thread.content,
        attachedImages: thread.attachedImages || [],
        postedBy: thread.userId, // Ensure this matches the backend
        categoryId: thread.categoryId, // Use categoryId instead of categories
        createdAt: new Date(thread.createdAt),
        likes: thread.likes || 0,
        lastActive: thread.lastActive || "Just now",
      }));

      setThreads(transformedThreads);
    } catch (error) {
      console.error("Error fetching threads:", error);
    }
  };

  // Create a new thread
  const createThread = async (thread: Omit<Thread, "id">) => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        throw new Error("No authentication token found");
      }

      const response = await fetch(`${BASE_URL}/api/threads`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(thread),
      });

      if (!response.ok) {
        throw new Error("Failed to create thread");
      }

      const newThread = await response.json();
      setThreads((prev) => [...prev, newThread]);
    } catch (error) {
      console.error("Error creating thread:", error);
    }
  };

  // Update a thread
  const updateThread = async (thread: Thread) => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        throw new Error("No authentication token found");
      }

      const response = await fetch(`${BASE_URL}/api/threads/${thread.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(thread),
      });

      if (!response.ok) {
        throw new Error("Failed to update thread");
      }

      const updatedThread = await response.json();
      setThreads((prev) =>
        prev.map((t) => (t.id === updatedThread.id ? updatedThread : t))
      );
    } catch (error) {
      console.error("Error updating thread:", error);
    }
  };

  // Delete a thread
  const deleteThread = async (threadId: number) => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        throw new Error("No authentication token found");
      }

      const response = await fetch(`${BASE_URL}/api/threads/${threadId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to delete thread");
      }

      setThreads((prev) => prev.filter((t) => t.id !== threadId));
    } catch (error) {
      console.error("Error deleting thread:", error);
    }
  };

  const likeThread = async (threadId: number) => {
    try {
      const response = await fetch(`${BASE_URL}/api/threads/${threadId}/like`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (!response.ok) throw new Error("Failed to like thread");
      // Update the local state
      setThreads((prev) =>
        prev.map((t) => (t.id === threadId ? { ...t, likes: t.likes + 1 } : t))
      );
    } catch (err) {
      console.error("Error liking thread:", err);
    }
  };

  // Fetch threads on mount
  useEffect(() => {
    fetchThreads();
  }, []);

  return (
    <ThreadContext.Provider
      value={{
        threads,
        fetchThreads,
        createThread,
        updateThread,
        deleteThread,
        likeThread,
      }}
    >
      {children}
    </ThreadContext.Provider>
  );
};

export const useThreads = () => {
  const context = useContext(ThreadContext);
  if (!context) {
    throw new Error("useThreads must be used within a ThreadProvider");
  }
  return context;
};
