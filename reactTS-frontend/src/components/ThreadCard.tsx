import React from "react";
import { useNavigate } from "react-router-dom";
import { Thread } from "../models/models";
import { FaBookmark, FaEdit, FaTrash } from "react-icons/fa";
import { useUser } from "../context/UserContext";
import { useAuth } from "../context/AuthContext"; // Import useAuth to get the current user
// import { useThreads } from "../context/ThreadContext";

interface ThreadCardProps {
  thread: Thread;
  onDelete: () => void;
  onEdit: () => void;
}

const ThreadCard: React.FC<ThreadCardProps> = ({
  thread,
  onDelete,
  onEdit,
}) => {
  const navigate = useNavigate();
  const { fetchUserDetails } = useUser();
  const { user } = useAuth(); // Get the current user from AuthContext
  // const { fetchThreads } = useThreads();
  const [postedBy, setPostedBy] = React.useState<string>("Unknown User"); // Local state for user details

  React.useEffect(() => {
    const fetchUser = async () => {
      try {
        console.log("Current thread: ", thread.postedBy);
        // if (thread.postedBy == undefined) {
        //   await fetchThreads();
        // }
        const fetchedUser = await fetchUserDetails(thread.postedBy); // Fetch user details
        if (fetchedUser) {
          setPostedBy(fetchedUser.username); // Update local state
        }
      } catch (error) {
        console.error("Error fetching user details:", error);
      }
    };
    console.log("changing user");
    console.log("Curernt User ID: ", user?.id);
    fetchUser();
  }, [thread.postedBy, fetchUserDetails]);

  // Check if the thread belongs to the current user
  var isCurrentUserThread: Boolean = user?.id === thread.postedBy;

  return (
    <div
      className="bg-white p-4 rounded-lg shadow mb-4 hover:shadow-lg transition-shadow cursor-pointer"
      onClick={() => navigate(`/threads/${thread.id}`)}
    >
      <h2 className="text-xl font-bold mb-2">{thread.title}</h2>
      <p className="text-gray-700 mb-4 line-clamp-3">{thread.content}</p>
      <div className="flex flex-wrap gap-2">
        {thread.attachedImages &&
          thread.attachedImages.map((image, index) => (
            <img
              key={index}
              src={image}
              alt={`Thread Image ${index + 1}`}
              className="w-24 h-24 object-cover rounded-lg"
            />
          ))}
      </div>
      <div className="mt-4 text-sm text-gray-500">
        <span>Posted by: {postedBy}</span>
        <span className="mx-2">•</span>
        <span>{new Date(thread.createdAt).toLocaleDateString()}</span>
      </div>
      <div className="flex space-x-3 mt-4" onClick={(e) => e.stopPropagation()}>
        <FaBookmark
          className={`text-gray-500 hover:text-gray-700 cursor-pointer`}
        />

        {/* Show Edit and Trash icons only if the thread belongs to the current user */}
        {isCurrentUserThread && (
          <>
            <button
              onClick={onEdit}
              className="text-blue-500 hover:text-blue-600 cursor-pointer"
            >
              <FaEdit />
            </button>

            <button
              onClick={onDelete}
              className="text-red-500 hover:text-red-600 cursor-pointer"
            >
              <FaTrash />
            </button>
          </>
        )}
      </div>
    </div>
  );
};

export default ThreadCard;
