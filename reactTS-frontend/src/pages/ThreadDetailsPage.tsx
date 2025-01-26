import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useThreads } from "../context/ThreadContext";
import { useUser } from "../context/UserContext";
import { Thread } from "../models/models";

const ThreadDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { threads } = useThreads();
  const { fetchUserDetails } = useUser();
  const [thread, setThread] = useState<Thread | null>(null);
  const [postedBy, setPostedBy] = useState<string>("Loading...");

  useEffect(() => {
    const fetchThreadDetails = () => {
      const foundThread = threads.find((t) => t.id === parseInt(id!));
      if (foundThread) {
        setThread(foundThread);

        // Fetch the username of the user who posted the thread
        fetchUserDetails(foundThread.postedBy)
          .then((user) => {
            if (user) {
              setPostedBy(user.username);
            } else {
              setPostedBy("Unknown User");
            }
          })
          .catch(() => {
            setPostedBy("Unknown User");
          });
      }
    };

    fetchThreadDetails();
  }, [id, threads, fetchUserDetails]);

  if (!thread) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>{thread.title}</h1>
      <p>{thread.content}</p>
      <p>Posted by: {postedBy}</p>
      <p>Created at: {new Date(thread.createdAt).toLocaleString()}</p>
    </div>
  );
};

export default ThreadDetailsPage;
