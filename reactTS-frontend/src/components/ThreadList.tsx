import React from "react";
import { useThreads } from "../context/ThreadContext";
import ThreadCard from "./ThreadCard";
import { Thread } from "../models/models";
import { useNavigate } from "react-router-dom";

interface ThreadListProps {
  onEdit: (thread: Thread) => void;
}

const ThreadList: React.FC<ThreadListProps> = ({}) => {
  const { threads, deleteThread } = useThreads();
  const navigate = useNavigate();

  if (!threads) {
    return <div>Loading threads...</div>;
  }

  return (
    <div>
      {threads.map((thread: Thread) => (
        <ThreadCard
          key={thread.id}
          thread={thread}
          onDelete={() => deleteThread(thread.id)}
          onEdit={() =>
            navigate(`/edit-thread/${thread.id}`, { state: { thread } })
          }
        />
      ))}
    </div>
  );
};

export default ThreadList;
