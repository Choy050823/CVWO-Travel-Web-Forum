import React from "react";
import { Comment } from "../models/models";
import { useUser } from "../context/UserContext";

interface CommentProps {
  comment: Comment;
}

const CommentComponent: React.FC<CommentProps> = ({ comment }) => {
  const { userDetails, fetchUserDetails } = useUser();
  const [postedBy, setPostedBy] = React.useState<string>("");

  React.useEffect(() => {
    const fetchUser = async () => {
      await fetchUserDetails(comment.postedBy);
      setPostedBy(userDetails?.username || "Unknown User");
    };
    fetchUser();
  }, [comment.postedBy, fetchUserDetails, userDetails]);

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
