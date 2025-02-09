import React, { useState, useEffect } from "react";
import { useNavigate, useParams, useLocation } from "react-router-dom";
import { useThreads } from "../context/ThreadContext";
import { useAuth } from "../context/AuthContext";
// import BASE_URL from "../config";

const BASE_URL = import.meta.env.VITE_API_URL;

const CreateEditThreadPage: React.FC = () => {
  const { id } = useParams();
  const { createThread, updateThread } = useThreads();
  const { user } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  // Get the thread data from the navigation state (if editing)
  const thread = location.state?.thread;

  // State for form fields
  const [title, setTitle] = useState(thread?.title || "");
  const [content, setContent] = useState(thread?.content || "");
  const [attachedImages, setAttachedImages] = useState<File[]>([]);
  const [imagePreviews, setImagePreviews] = useState<string[]>([]);

  // If editing, pre-fill the form with the thread data
  useEffect(() => {
    if (thread) {
      setTitle(thread.title);
      setContent(thread.content);
      setImagePreviews(thread.attachedImages || []);
    }
  }, [thread]);

  // Handle image selection
  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files);
      setAttachedImages(files);

      // Create image previews
      const previews = files.map((file) => URL.createObjectURL(file));
      setImagePreviews(previews);
    }
  };

  // Handle form submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      // Upload images and get their URLs
      const imageUrls = await uploadImages(attachedImages);

      // Create or update the thread
      const threadData = {
        id: thread?.id || Date.now(),
        title,
        content,
        attachedImages: imageUrls,
        postedBy: thread?.userId || user?.id || -1,
        categoryId: thread?.categoryId || 1, // Assuming a default category ID of 1
        createdAt: thread?.createdAt || new Date(),
        likes: thread?.likes || 0,
        lastActive: "Just now",
      };

      if (id) {
        await updateThread(threadData);
      } else {
        await createThread(threadData);
      }

      navigate("/"); // Redirect to home page
    } catch (error) {
      console.error("Error during thread submission:", error);
    }
  };

  // Function to upload images to the backend
  const uploadImages = async (images: File[]) => {
    const formData = new FormData();
    images.forEach((image) => formData.append("images", image));

    const token = localStorage.getItem("token");
    if (!token) {
      throw new Error("No authentication token found");
    }

    const response = await fetch(`${BASE_URL}/api/upload-images`, {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    });

    if (!response.ok) {
      const errorText = await response.text(); // Get detailed error from backend
      console.error("Server response:", errorText);
      throw new Error("Failed to upload images");
    }

    const data = await response.json();
    return data.imageUrls; // Array of uploaded image URLs
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-6">
        {id ? "Edit Thread" : "Create New Thread"}
      </h1>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Title Input */}
        <div>
          <label className="block text-sm font-medium mb-1">Title</label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-black"
            required
          />
        </div>

        {/* Content Input */}
        <div>
          <label className="block text-sm font-medium mb-1">Content</label>
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            className="w-full px-4 py-2 border rounded-lg h-48 focus:ring-2 focus:ring-black"
            required
          />
        </div>

        {/* Image Upload */}
        <div>
          <label className="block text-sm font-medium mb-1">
            Attach Images
          </label>
          <input
            type="file"
            multiple
            onChange={handleImageUpload}
            className="w-full px-4 py-2 border rounded-lg"
          />
        </div>

        {/* Image Previews */}
        <div className="flex flex-wrap gap-2">
          {imagePreviews.map((preview, index) => (
            <img
              key={index}
              src={preview}
              alt={`Preview ${index + 1}`}
              className="w-24 h-24 object-cover rounded-lg"
            />
          ))}
        </div>

        {/* Form Actions */}
        <div className="flex justify-end space-x-4">
          <button
            type="button"
            onClick={() => navigate(-1)}
            className="px-6 py-2 border rounded-lg hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="px-6 py-2 bg-black text-white rounded-lg hover:bg-gray-800"
          >
            {id ? "Update Thread" : "Create Thread"}
          </button>
        </div>
      </form>
    </div>
  );
};

export default CreateEditThreadPage;
