import { useEffect, useState } from "react";
import { useChatStore } from "../store/useChatStore";
import { useAuthStore } from "../store/useAuthStore";
import SidebarSkeleton from "./skeletons/SidebarSkeleton";
import { Users } from "lucide-react";

const Sidebar = () => {
  const { getUsers, users, selectedUser, setSelectedUser, isUsersLoading } =
    useChatStore();

  const { onlineUsers } = useAuthStore();
  const [showOnlineOnly, setShowOnlineOnly] = useState(false);

  useEffect(() => {
    getUsers();
  }, [getUsers]);

  const filteredUsers = showOnlineOnly
    ? users.filter((user) => onlineUsers.includes(user._id))
    : users;

  if (isUsersLoading) return <SidebarSkeleton />;

  return (
    <aside className="h-full w-full lg:w-72 border-r border-gray-800/50 flex flex-col transition-all duration-200 bg-black/20 backdrop-blur-sm">
      {/* Header */}
      <div className="border-b border-gray-800/50 w-full p-5">
        <div className="flex items-center gap-2">
          <div className="p-2 rounded-xl bg-gradient-to-r from-gray-800 to-gray-900 shadow-lg">
            <Users className="w-5 h-5 text-white" />
          </div>
          <span className="font-semibold text-white">Contacts</span>
        </div>

        {/* Online filter toggle */}
        <div className="mt-4 flex items-center gap-2">
          <label className="cursor-pointer flex items-center gap-2">
            <input
              type="checkbox"
              checked={showOnlineOnly}
              onChange={(e) => setShowOnlineOnly(e.target.checked)}
              className="w-4 h-4 rounded bg-gray-800/50 border border-gray-700/50 text-gray-600 focus:ring-2 focus:ring-gray-600/50 focus:ring-offset-0"
            />
            <span className="text-sm text-gray-300 font-medium">
              Online only?
            </span>
          </label>
          <span className="text-xs text-gray-500 bg-gray-800/30 px-2 py-1 rounded-lg">
            ({onlineUsers.length - 1} online)
          </span>
        </div>
      </div>

      {/* User list */}
      <div className="overflow-y-auto overflow-x-hidden w-full py-3 space-y-1 custom-scrollbar flex-1">
        {filteredUsers.map((user) => (
          <button
            key={user._id}
            onClick={() => setSelectedUser(user)}
            className={`
              w-full p-3 mx-2 rounded-xl flex items-center gap-3
              transition-all duration-200 hover:bg-gray-800/40
              ${
                selectedUser?._id === user._id
                  ? "bg-gradient-to-r from-gray-800/60 to-gray-900/60 border border-gray-700/50 shadow-lg"
                  : "hover:bg-gray-800/20"
              }
            `}
          >
            <div className="relative flex-shrink-0">
              <img
                src={user.profilePic || "/avatar.png"}
                alt={user.name}
                className="w-12 h-12 object-cover rounded-full border-2 border-gray-700/50"
              />
              {onlineUsers.includes(user._id) && (
                <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-500 rounded-full ring-2 ring-black shadow-lg" />
              )}
            </div>

            {/* User info - always visible now */}
            <div className="text-left min-w-0 flex-1">
              <div className="font-semibold text-white truncate">
                {user.fullName}
              </div>
              <div className="text-sm text-gray-400 font-medium">
                {onlineUsers.includes(user._id) ? (
                  <span className="text-green-400">Online</span>
                ) : (
                  <span className="text-gray-500">Offline</span>
                )}
              </div>
            </div>
          </button>
        ))}

        {filteredUsers.length === 0 && (
          <div className="text-center text-gray-500 py-8">
            <div className="space-y-2">
              <Users className="w-8 h-8 mx-auto text-gray-600" />
              <p className="text-sm font-medium">No online users</p>
            </div>
          </div>
        )}
      </div>
    </aside>
  );
};

export default Sidebar;
