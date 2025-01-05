import React from "react";
// import { useNavigate } from 'react-router-dom';
import { Thread } from "../models/models";
import Navbar from "../components/Navbar";

const HomePage: React.FC = () => {
  const threads: Thread[] = [
    {
      id: 1,
      title: "Thread 1",
      content: "Content of thread 1",
      postedBy: 0,
      categories: [],
      attachedImages: [],
      createdAt: new Date(),
      likes: 0,
      lastActive: "",
    },
    {
      id: 1,
      title: "Thread 1",
      content: "Content of thread 1",
      postedBy: 0,
      categories: [],
      attachedImages: [],
      createdAt: new Date(),
      likes: 0,
      lastActive: "",
    },
    {
      id: 1,
      title: "Thread 1",
      content: "Content of thread 1",
      postedBy: 0,
      categories: [],
      attachedImages: [],
      createdAt: new Date(),
      likes: 0,
      lastActive: "",
    },
  ];

  return (
    <div>
      <Navbar user={null} />
      <h1>Threads</h1>
      <ul>
        {threads.map((thread) => (
          <li key={thread.id}>
            <h2>{thread.title}</h2>
            {/* Find when is it last active to see whether it is outdated */}
            <p>Last Active: {Date.now() - thread.createdAt.getTime()}</p>
            <p>{thread.content}</p>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default HomePage;
