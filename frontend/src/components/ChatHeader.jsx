import { X, ChevronLeft } from "lucide-react";
import { useAuthStore } from "../store/useAuthStore";
import { useChatStore } from "../store/useChatStore";

const ChatHeader = () => {
  const { selectedUser, setSelectedUser } = useChatStore();
  const { onlineUsers } = useAuthStore();

  const handleBack = () => {
    setSelectedUser(null);
  };

  return (
    <div className="p-4 border-b border-gray-800/50 bg-black/20 backdrop-blur-sm">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          {/* iOS style back button for mobile */}
          <button
            onClick={handleBack}
            className="lg:hidden p-1 text-blue-400 hover:text-blue-300 transition-colors duration-200 -ml-1"
          >
            <ChevronLeft className="w-6 h-6" strokeWidth={2.5} />
          </button>

          {/* Avatar */}
          <div className="relative">
            <img
              src={selectedUser.profilePic || "/avatar.png"}
              alt={selectedUser.fullName}
              className="w-12 h-12 rounded-full border-2 border-gray-700/50 shadow-lg"
            />
            {onlineUsers.includes(selectedUser._id) && (
              <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full ring-2 ring-black shadow-lg" />
            )}
          </div>

          {/* User info */}
          <div>
            <h3 className="font-semibold text-white text-lg">
              {selectedUser.fullName}
            </h3>
            <p className="text-sm font-medium">
              {onlineUsers.includes(selectedUser._id) ? (
                <span className="text-green-400">Online</span>
              ) : (
                <span className="text-gray-500">Offline</span>
              )}
            </p>
          </div>
        </div>

        {/* Close button - hidden on mobile, back button handles this */}
        <button
          onClick={() => setSelectedUser(null)}
          className="hidden lg:flex p-2 rounded-xl bg-gray-800/50 border border-gray-700/50 text-gray-400 hover:bg-gray-700/50 hover:text-white transition-all duration-200"
        >
          <X className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
};

export default ChatHeader;