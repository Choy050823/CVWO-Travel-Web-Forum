// CommentContext.tsx
import React, { createContext, useState, useContext } from "react";
import { Comment } from "../models/models";
import BASE_URL from "../config";

interface CommentContextType {
  comments: Comment[];
  loading: boolean;
  error: string | null;
  fetchComments: (threadId: number) => Promise<void>;
  createComment: (comment: Omit<Comment, "id">) => Promise<void>;
  updateComment: (comment: Comment) => Promise<void>;
  deleteComment: (commentId: number) => Promise<void>;
  upvoteComment: (commentId: number) => Promise<void>;
  downvoteComment: (commentId: number) => Promise<void>;
}

const CommentContext = createContext<CommentContextType | undefined>(undefined);

export const CommentProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [comments, setComments] = useState<Comment[]>([]); // Initialize as empty array
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchComments = async (threadId: number) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `${BASE_URL}/api/threads/${threadId}/comments`
      );
      if (!response.ok) throw new Error("Failed to fetch comments");
      const data = await response.json();
      setComments(data || []); // Ensure data is an array
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
      setComments([]); // Reset comments on error
    } finally {
      setLoading(false);
    }
  };

  const createComment = async (comment: Omit<Comment, "id">) => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No authentication token found");

    try {
      const response = await fetch(`${BASE_URL}/api/comments`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(comment),
      });

      if (!response.ok) throw new Error("Failed to create comment");

      const newComment = await response.json();
      setComments((prev) => [newComment, ...(prev || [])]); // Ensure prev is an array
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  const updateComment = async (comment: Comment) => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No authentication token found");

    try {
      const response = await fetch(`${BASE_URL}/api/comments/${comment.id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(comment),
      });

      if (!response.ok) throw new Error("Failed to update comment");

      const updatedComment = await response.json();
      setComments((prev) =>
        prev.map((c) => (c.id === updatedComment.id ? updatedComment : c))
      );
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  const deleteComment = async (commentId: number) => {
    const token = localStorage.getItem("token");
    if (!token) throw new Error("No authentication token found");

    try {
      const response = await fetch(`${BASE_URL}/api/comments/${commentId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) throw new Error("Failed to delete comment");

      setComments((prev) => prev.filter((c) => c.id !== commentId));
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("An unknown error occurred");
      }
    }
  };

  // const voteComment = async (commentId: number, type: "up" | "down") => {
  //   const token = localStorage.getItem("token");
  //   if (!token) throw new Error("No authentication token found");

  //   try {
  //     const response = await fetch(
  //       `${BASE_URL}/api/comments/${commentId}/vote`,
  //       {
  //         method: "POST",
  //         headers: {
  //           "Content-Type": "application/json",
  //           Authorization: `Bearer ${token}`,
  //         },
  //         body: JSON.stringify({ type }),
  //       }
  //     );

  //     if (!response.ok) throw new Error("Failed to vote on comment");

  //     const updatedComment = await response.json();
  //     setComments((prev) =>
  //       prev.map((c) => (c.id === updatedComment.id ? updatedComment : c))
  //     );
  //   } catch (err) {
  //     if (err instanceof Error) {
  //       setError(err.message);
  //     } else {
  //       setError("An unknown error occurred");
  //     }
  //   }
  // };

  const upvoteComment = async (commentId: number) => {
    try {
      const response = await fetch(`/api/comments/${commentId}/upvote`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (!response.ok) throw new Error("Failed to upvote comment");
      // Update the local state
      setComments((prev) =>
        prev.map((c) =>
          c.id === commentId ? { ...c, upvotes: c.upvotes + 1 } : c
        )
      );
    } catch (err) {
      console.error("Error upvoting comment:", err);
    }
  };

  const downvoteComment = async (commentId: number) => {
    try {
      const response = await fetch(`/api/comments/${commentId}/downvote`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
      });
      if (!response.ok) throw new Error("Failed to downvote comment");
      // Update the local state
      setComments((prev) =>
        prev.map((c) =>
          c.id === commentId ? { ...c, downvotes: c.downvotes + 1 } : c
        )
      );
    } catch (err) {
      console.error("Error downvoting comment:", err);
    }
  };

  return (
    <CommentContext.Provider
      value={{
        comments,
        loading,
        error,
        fetchComments,
        createComment,
        updateComment,
        deleteComment,
        upvoteComment,
        downvoteComment,
      }}
    >
      {children}
    </CommentContext.Provider>
  );
};

export const useComments = () => {
  const context = useContext(CommentContext);
  if (!context) {
    throw new Error("useComments must be used within a CommentProvider");
  }
  return context;
};
