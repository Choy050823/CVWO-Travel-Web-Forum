import React, { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import BASE_URL from "../config";

const CreateEditThreadPage: React.FC = () => {
  const { id } = useParams<{ id?: string }>();
  const [title, setTitle] = useState<string>("");
  const [content, setContent] = useState<string>("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const threadData = { title, content };

    try {
      const url = id
        ? `${BASE_URL}/api/threads/${id}`
        : `${BASE_URL}/api/threads`;
      const method = id ? "PUT" : "POST";

      const response = await fetch(url, {
        method,
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify(threadData),
      });

      if (!response.ok) {
        throw new Error(
          id ? "Failed to update thread" : "Failed to create thread"
        );
      }

      navigate("/"); // Redirect to home page
    } catch (error) {
      console.error("Error during thread submission:", error);
    }
  };

  return (
    <div>
      <h1>{id ? "Edit Thread" : "Create Thread"}</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="title">Title:</label>
          <input
            type="text"
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
        <div>
          <label htmlFor="content">Content:</label>
          <textarea
            id="content"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
          />
        </div>
        <button type="submit">{id ? "Update" : "Create"}</button>
      </form>
    </div>
  );
};

export default CreateEditThreadPage;
