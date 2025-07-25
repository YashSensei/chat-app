import { useChatStore } from "../store/useChatStore";

import Sidebar from "../components/Sidebar";
import NoChatSelected from "../components/NoChatSelected";
import ChatContainer from "../components/ChatContainer";

const HomePage = () => {
  const { selectedUser } = useChatStore();

  return (
    <div className="h-screen bg-gradient-to-br from-black via-gray-900 to-black relative overflow-hidden">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gray-800 rounded-full mix-blend-multiply filter blur-xl opacity-5 animate-blob"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gray-700 rounded-full mix-blend-multiply filter blur-xl opacity-5 animate-blob animation-delay-2000"></div>
        <div className="absolute top-40 left-40 w-80 h-80 bg-gray-600 rounded-full mix-blend-multiply filter blur-xl opacity-5 animate-blob animation-delay-4000"></div>
      </div>

      {/* Mobile: Full screen layout, Desktop: Centered container */}
      <div className="relative z-10 h-full pt-16 lg:flex lg:items-center lg:justify-center lg:pt-20 lg:px-4 overflow-hidden">
        <div className="h-full w-full backdrop-blur-xl bg-black/30 border border-gray-800/50 shadow-2xl overflow-hidden lg:rounded-3xl lg:max-w-6xl lg:h-[calc(100vh-8rem)]">
          <div className="flex h-full">
            {/* Sidebar - always visible on desktop, conditional on mobile */}
            <div
              className={`${
                selectedUser ? "hidden lg:block" : "w-full lg:w-auto"
              } lg:w-72`}
            >
              <Sidebar />
            </div>

            {/* Chat area - shows chat when selected, welcome screen when not */}
            <div
              className={`flex-1 ${
                selectedUser ? "flex" : "hidden lg:flex"
              } flex-col min-h-0`}
            >
              {!selectedUser ? <NoChatSelected /> : <ChatContainer />}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default HomePage;
