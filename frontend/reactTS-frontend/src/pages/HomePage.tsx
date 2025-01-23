import React, { useEffect, useState } from "react";
import { Thread } from "../models/models";
import Navbar from "../components/Navbar";
import BASE_URL from "../config";

const HomePage: React.FC = () => {
  const [threads, setThreads] = useState<Thread[]>([]);

  useEffect(() => {
    const fetchThreads = async () => {
      try {
        const response = await fetch(`${BASE_URL}/api/threads`);
        if (!response.ok) {
          throw new Error("Failed to fetch threads");
        }
        const data = await response.json();
        setThreads(data);
      } catch (error) {
        console.error("Error fetching threads:", error);
      }
    };

    fetchThreads();
  }, []);

  return (
    <div>
      <Navbar user={null} />
      <h1>Threads</h1>
      <ul>
        {threads.map((thread) => (
          <li key={thread.id}>
            <h2>{thread.title}</h2>
            <p>Last Active: {new Date(thread.createdAt).toLocaleString()}</p>
            <p>{thread.content}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default HomePage;
