import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { Thread } from "../models/models";
import BASE_URL from "../config";

const ThreadDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [thread, setThread] = useState<Thread | null>(null);

  useEffect(() => {
    const fetchThreadDetails = async () => {
      try {
        const response = await fetch(`${BASE_URL}/api/threads/${id}`);
        if (!response.ok) {
          throw new Error("Failed to fetch thread details");
        }
        const data = await response.json();
        setThread(data);
      } catch (error) {
        console.error("Error fetching thread details:", error);
      }
    };

    fetchThreadDetails();
  }, [id]);

  if (!thread) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>{thread.title}</h1>
      <p>{thread.content}</p>
      <p>Posted by: {thread.postedBy}</p>
      <p>Created at: {new Date(thread.createdAt).toLocaleString()}</p>
    </div>
  );
};

export default ThreadDetailsPage;
