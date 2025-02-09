import React, { useState, useEffect } from "react";
import { Category } from "../models/models";

const BASE_URL = import.meta.env.VITE_API_URL.replace(/\/$/, "");

const TopicsList: React.FC = () => {
  const [topics, setTopics] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCategories = async (): Promise<Category[]> => {
    try {
      const response = await fetch(`${BASE_URL}/api/categories`);
      if (!response.ok) {
        throw new Error(`Failed to fetch categories: ${response.statusText}`);
      }
      const data = await response.json();
      // Ensure the response is an array, even if empty
      return Array.isArray(data) ? data : [];
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "An error occurred";
      setError(errorMessage);
      return [];
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        const res = await fetchCategories();
        setTopics(res || []);
      } catch (error) {
        setError("Failed to fetch topics");
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  if (isLoading) {
    return (
      <div className="bg-white p-4 rounded-lg shadow">
        <h3 className="font-semibold text-lg mb-4">Topics For You</h3>
        <div className="animate-pulse flex space-x-2">
          <div className="h-8 w-20 bg-gray-200 rounded-full"></div>
          <div className="h-8 w-24 bg-gray-200 rounded-full"></div>
          <div className="h-8 w-16 bg-gray-200 rounded-full"></div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white p-4 rounded-lg shadow">
        <h3 className="font-semibold text-lg mb-4">Topics For You</h3>
        <p className="text-red-500">
          Failed to load topics. Please try again later.
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white p-4 rounded-lg shadow">
      <h3 className="font-semibold text-lg mb-4">Topics For You</h3>
      <div className="flex flex-wrap gap-2">
        {topics.length === 0 ? (
          <p className="text-gray-500">No topics available</p>
        ) : (
          topics.map((topic, index) => (
            <span
              key={index}
              className="px-3 py-1 bg-gray-200 text-gray-700 text-sm rounded-full cursor-pointer hover:bg-gray-300 transition-colors"
            >
              {topic.name}
            </span>
          ))
        )}
      </div>
    </div>
  );
};

export default TopicsList;
