import React from "react";
import { Comment, User } from "../models/models";
import { useUser } from "../context/UserContext";

interface CommentProps {
  comment: Comment;
}

const CommentComponent: React.FC<CommentProps> = ({ comment }) => {
  const { fetchUserDetails } = useUser();
  const [postedBy, setPostedBy] = React.useState<string>("");
  const [user, setUser] = React.useState<User | null>(null);

  React.useEffect(() => {
    const fetchUser = async () => {
      const res = await fetchUserDetails(comment.userId);
      setUser(res);
      setPostedBy(user?.username || "Unknown User");
    };
    fetchUser();
  }, [comment.userId, fetchUserDetails]);

  return (
    <div className="bg-gray-50 p-4 rounded-lg mb-4">
      <p className="text-gray-700 mb-4">{comment.content}</p>
      <div className="flex space-x-2">
        {comment.attachedImages.map((image, index) => (
          <img
            key={index}
            src={image}
            alt={`Comment Image ${index + 1}`}
            className="w-24 h-24 object-cover rounded-lg"
          />
        ))}
      </div>
      <div className="mt-4 text-sm text-gray-500">
        Posted by: {postedBy} | {comment.createdAt.toLocaleDateString()}
      </div>
    </div>
  );
};

export default CommentComponent;
