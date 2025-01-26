import React, { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { useComments } from "../context/CommentContext";
import { useThreads } from "../context/ThreadContext";
import { useAuth } from "../context/AuthContext";
import { FaComment, FaRegClock, FaUserCircle } from "react-icons/fa";
import { AiFillLike } from "react-icons/ai";
import { FiArrowUp, FiArrowDown } from "react-icons/fi";

const ThreadDetailsPage: React.FC = () => {
  // All hooks declared at the top (before any conditionals)
  const { id: threadId } = useParams<{ id: string }>();
  const { comments, loading, error, fetchComments, createComment } =
    useComments();
  const { threads, likeThread } = useThreads();
  const { user } = useAuth();
  const [newComment, setNewComment] = useState("");
  const { upvoteComment, downvoteComment } = useComments();
  const [votedComments, setVotedComments] = useState<number[]>([]);
  const [votedThreads, setVotedThreads] = useState<number[]>([]);

  // Fetch comments when the threadId changes
  useEffect(() => {
    if (threadId) {
      fetchComments(parseInt(threadId));
    }
  }, [threadId]);

  // Find the current thread
  const thread = threads.find((t) => t.id === parseInt(threadId || ""));

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newComment.trim() || !thread || !user) return;

    await createComment({
      content: newComment,
      threadId: thread.id,
      userId: user.id,
      attachedImages: [],
      upvotes: 0,
      downvotes: 0,
      createdAt: new Date(),
    });

    setNewComment("");
  };

  const handleUpvote = async (commentId: number) => {
    if (votedComments.includes(commentId)) return;
    await upvoteComment(commentId);
    setVotedComments((prev) => [...prev, commentId]);
  };

  const handleDownvote = async (commentId: number) => {
    if (votedComments.includes(commentId)) return;
    await downvoteComment(commentId);
    setVotedComments((prev) => [...prev, commentId]);
  };

  const handleLike = async (threadId: number) => {
    if (votedThreads.includes(threadId)) return;
    await likeThread(threadId);
    setVotedThreads((prev) => [...prev, threadId]);
  };

  // Conditional returns moved below all hook declarations
  if (!thread) {
    return <div className="text-center p-8">Thread not found</div>;
  }

  if (loading) {
    return <div className="text-center p-8">Loading comments...</div>;
  }

  if (error) {
    return <div className="text-center p-8 text-red-500">{error}</div>;
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      {/* Thread Header */}
      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100 mb-6">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">
          {thread.title}
        </h1>

        <div className="flex items-center gap-3 text-sm text-gray-500 mb-4">
          <span className="flex items-center gap-1">
            <FaUserCircle className="text-gray-400" />
            {thread.postedBy}
          </span>
          <span>•</span>
          <span className="flex items-center gap-1">
            <FaRegClock className="text-gray-400" />
            {new Date(thread.createdAt).toLocaleDateString()}
          </span>
        </div>

        <p className="text-gray-700 leading-relaxed mb-6">{thread.content}</p>

        {(thread.attachedImages?.length ?? 0) > 0 && (
          <div className="grid grid-cols-2 gap-4 mb-6">
            {thread.attachedImages?.map((img, index) => (
              <img
                key={index}
                src={img}
                alt={`Attachment ${index + 1}`}
                className="rounded-lg object-cover h-48 w-full"
              />
            ))}
          </div>
        )}

        <div className="flex items-center gap-6 text-gray-500 border-t pt-4">
          <button
            onClick={() => handleLike(thread.id)}
            disabled={votedThreads.includes(thread.id)}
            className="flex items-center gap-2 hover:text-blue-600 disabled:opacity-50"
          >
            <AiFillLike className="text-lg" />
            {thread.likes}
          </button>
          <button className="flex items-center gap-2 hover:text-blue-600">
            <FaComment className="text-lg" />
            {comments?.length || 0} Comments
          </button>
        </div>
      </div>

      {/* Comments Section */}
      <div className="bg-white rounded-xl p-6 shadow-sm border border-gray-100">
        <h2 className="text-xl font-semibold text-gray-800 mb-6">
          Comments ({comments?.length || 0})
        </h2>

        {/* Comment Form */}
        {user ? (
          <form onSubmit={handleCommentSubmit} className="mb-8">
            <textarea
              value={newComment}
              onChange={(e) => setNewComment(e.target.value)}
              placeholder="Write your comment..."
              className="w-full p-4 border border-gray-200 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
              rows={4}
              required
            />
            <div className="mt-4 flex justify-end">
              <button
                type="submit"
                className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-full font-medium transition-colors"
              >
                Post Comment
              </button>
            </div>
          </form>
        ) : (
          <div className="text-center py-6 text-gray-500">
            Please{" "}
            <Link to="/login" className="text-blue-600 hover:underline">
              login
            </Link>{" "}
            to comment
          </div>
        )}

        {/* Display Comments */}
        <div className="space-y-6">
          {comments?.map((comment) => (
            <div
              key={comment.id}
              className="border-b border-gray-100 pb-6 last:border-0"
            >
              <div className="flex items-center gap-3 text-sm text-gray-500 mb-3">
                <FaUserCircle className="text-gray-400" />
                <span>{comment.author}</span>
                <span>•</span>
                <span>{new Date(comment.createdAt).toLocaleDateString()}</span>
              </div>
              <p className="text-gray-700 mb-4">{comment.content}</p>
              <div className="flex items-center gap-4 text-gray-500">
                <button
                  onClick={() => handleUpvote(comment.id)}
                  disabled={votedComments.includes(comment.id)}
                  className="flex items-center gap-1 hover:text-blue-600 disabled:opacity-50"
                >
                  <FiArrowUp className="text-lg" />
                  {comment.upvotes}
                </button>

                <button
                  onClick={() => handleDownvote(comment.id)}
                  disabled={votedComments.includes(comment.id)}
                  className="flex items-center gap-1 hover:text-blue-600 disabled:opacity-50"
                >
                  <FiArrowDown className="text-lg" />
                  {comment.downvotes}
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ThreadDetailsPage;
