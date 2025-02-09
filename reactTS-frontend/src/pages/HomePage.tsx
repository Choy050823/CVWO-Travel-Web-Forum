// import { useThreads } from "../context/ThreadContext";
import ThreadList from "../components/ThreadList";
import { IoMdAddCircle } from "react-icons/io";
import TrendsList from "../components/TrendsLists";
import TopicsList from "../components/TopicsList";
import { useNavigate } from "react-router";
import { useAuth } from "../context/AuthContext";
// import React from "react";

const HomePage = () => {
  // const { threads, fetchThreads } = useThreads(); // Fetch threads from ThreadContext
  const { user } = useAuth();
  const navigate = useNavigate();

  // React.useEffect(() => {
  //   fetchThreads(); // Fetch threads when the component mounts
  // }, []);

  return (
    <div className="flex justify-center gap-8 p-6">
      {/* Left Side - Thread List */}
      <section className="w-full max-w-3xl">
        <ThreadList
          onEdit={(thread) => navigate(`/edit-thread/${thread!.id}`)}
        />
      </section>

      {/* Right Sidebar */}
      <div className="w-[30%] min-w-[300px] space-y-6">
        {/* Create New Thread Button */}
        {user && (
          <div className="bg-white p-4 rounded-lg shadow">
            <button
              onClick={() => navigate("/create-thread")}
              className="w-full bg-black text-white px-6 py-3 rounded-lg font-semibold hover:bg-gray-800 transition-colors flex items-center justify-center"
            >
              Create New Thread
              <IoMdAddCircle className={`ml-2 text-xl`} />
            </button>
          </div>
        )}

        {/* Today's Top Trends */}
        <TrendsList />

        {/* Topics for You */}
        <TopicsList />
      </div>
    </div>
  );
};

export default HomePage;
