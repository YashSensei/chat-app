import { MessageSquare } from "lucide-react";

const NoChatSelected = () => {
  return (
    <div className="w-full flex flex-1 flex-col items-center justify-center p-16 bg-black/10 backdrop-blur-sm">
      <div className="max-w-md text-center space-y-8">
        {/* Icon Display */}
        <div className="flex justify-center gap-4 mb-4">
          <div className="relative">
            <div className="w-20 h-20 rounded-3xl bg-gradient-to-r from-gray-800 to-gray-900 flex items-center justify-center shadow-2xl animate-pulse">
              <MessageSquare className="w-10 h-10 text-white" />
            </div>
            {/* Decorative elements */}
            <div className="absolute -top-2 -right-2 w-6 h-6 bg-gradient-to-r from-gray-700 to-gray-800 rounded-full opacity-60"></div>
            <div className="absolute -bottom-2 -left-2 w-4 h-4 bg-gradient-to-r from-gray-600 to-gray-700 rounded-full opacity-40"></div>
          </div>
        </div>

        {/* Welcome Text */}
        <div className="space-y-4">
          <h2 className="text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-gray-300 bg-clip-text text-transparent">
            Welcome to ChatApp!
          </h2>
          <p className="text-gray-400 text-lg leading-relaxed">
            Select a conversation from the sidebar to start chatting
          </p>
        </div>

        {/* Feature highlights */}
        <div className="space-y-3 pt-6">
          <div className="flex items-center justify-center gap-3 text-gray-300">
            <div className="w-2 h-2 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
            <span className="text-sm font-medium">Real-time messaging</span>
          </div>
          <div className="flex items-center justify-center gap-3 text-gray-300">
            <div className="w-2 h-2 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
            <span className="text-sm font-medium">Instant notifications</span>
          </div>
          <div className="flex items-center justify-center gap-3 text-gray-300">
            <div className="w-2 h-2 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
            <span className="text-sm font-medium">Secure conversations</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NoChatSelected;
