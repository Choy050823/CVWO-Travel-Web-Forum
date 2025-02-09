// import BASE_URL from "../config";
import { useState } from "react";
import { Category } from "../models/models";
import React from "react";

const BASE_URL = import.meta.env.VITE_API_URL;
const [topics, setTopics] = useState<Category[]>([]);

const fetchCategories = async (): Promise<Category[]> => {
  try {
    const response = await fetch(`${BASE_URL}/api/categories`);
    if (!response.ok) {
      throw new Error("Failed to fetch categories");
    }
    const res = await response.json();
    setTopics(res);
    return topics;
  } catch (error) {
    console.error("Error fetching categories:", error);
    return [];
  }
};

React.useEffect(() => {
  const fetchData = async () => {
    await fetchCategories();
  };
  fetchData();
}, [])




const TopicsList = () => {
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
