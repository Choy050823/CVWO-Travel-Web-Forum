import React, { useState, useEffect } from "react";
import { Category } from "../models/models";

const BASE_URL = import.meta.env.VITE_API_URL;

const TopicsList: React.FC = () => {
  const [topics, setTopics] = useState<Category[]>([]);

  const fetchCategories = async (): Promise<Category[]> => {
    try {
      const response = await fetch(`${BASE_URL}/api/categories`);
      if (!response.ok) {
        throw new Error("Failed to fetch categories");
      }
      const res = await response.json();
      return res;
    } catch (error) {
      console.error("Error fetching categories:", error);
      return [];
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      const res = await fetchCategories();
      setTopics(res);
      console.log(topics);
    };
    fetchData();
  }, []);

  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h3 className="font-semibold text-lg mb-4">Topics For You</h3>
      <div className="flex flex-wrap gap-2">
        {topics.map((topic, index) => (
          <span
            key={index}
            className="px-3 py-1 bg-gray-200 text-gray-700 text-sm rounded-full cursor-pointer"
          >
            {topic.name}
          </span>
        ))}
      </div>
    </div>
  );
};

export default TopicsList;
