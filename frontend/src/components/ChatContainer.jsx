import { useChatStore } from "../store/useChatStore";
import { useEffect, useRef, useState } from "react";

import ChatHeader from "./ChatHeader";
import MessageInput from "./MessageInput";
import MessageSkeleton from "./skeletons/MessageSkeleton";
import { useAuthStore } from "../store/useAuthStore";
import { formatMessageTime } from "../lib/utils";
import { X } from "lucide-react";

const ChatContainer = () => {
  const {
    messages,
    getMessages,
    isMessagesLoading,
    selectedUser,
    subscribeToMessages,
    unsubscribeFromMessages,
  } = useChatStore();
  const { authUser } = useAuthStore();
  const messageEndRef = useRef(null);
  const [previewImage, setPreviewImage] = useState(null);
  const [isModalClosing, setIsModalClosing] = useState(false);

  useEffect(() => {
    getMessages(selectedUser._id);

    subscribeToMessages();

    return () => unsubscribeFromMessages();
  }, [
    selectedUser._id,
    getMessages,
    subscribeToMessages,
    unsubscribeFromMessages,
  ]);

  useEffect(() => {
    if (messageEndRef.current && messages) {
      messageEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  const openImagePreview = (imageSrc) => {
    setPreviewImage(imageSrc);
    setIsModalClosing(false);
  };

  const closeImagePreview = () => {
    setIsModalClosing(true);
    setTimeout(() => {
      setPreviewImage(null);
      setIsModalClosing(false);
    }, 200);
  };

  if (isMessagesLoading) {
    return (
      <div className="h-full flex flex-col bg-black/10 backdrop-blur-sm">
        <ChatHeader />
        <MessageSkeleton />
        <MessageInput />
      </div>
    );
  }

  return (
    <div className="h-full flex flex-col bg-black/10 backdrop-blur-sm relative">
      <ChatHeader />

      <div className="flex-1 overflow-y-auto overflow-x-hidden p-6 custom-scrollbar">
        <div className="space-y-4">
          {messages.map((message) => (
            <div
              key={message._id}
              className={`flex ${
                message.senderId === authUser._id
                  ? "justify-end"
                  : "justify-start"
              } mb-4`}
            >
              <div
                className={`flex items-end gap-3 max-w-[80%] ${
                  message.senderId === authUser._id
                    ? "flex-row-reverse"
                    : "flex-row"
                }`}
              >
                {/* Avatar */}
                <div className="flex-shrink-0">
                  <img
                    src={
                      message.senderId === authUser._id
                        ? authUser.profilePic || "/avatar.png"
                        : selectedUser.profilePic || "/avatar.png"
                    }
                    alt="profile pic"
                    className="w-8 h-8 rounded-full border-2 border-gray-700/50"
                  />
                </div>

                {/* Message content */}
                <div
                  className={`flex flex-col ${
                    message.senderId === authUser._id
                      ? "items-end"
                      : "items-start"
                  }`}
                >
                  {/* Time */}
                  <div
                    className={`text-xs text-gray-500 mb-1 ${
                      message.senderId === authUser._id ? "mr-2" : "ml-2"
                    }`}
                  >
                    {formatMessageTime(message.createdAt)}
                  </div>

                  {/* Message bubble */}
                  <div
                    className={`px-4 py-3 rounded-2xl max-w-full ${
                      message.senderId === authUser._id
                        ? "bg-gradient-to-r from-gray-700 to-gray-800 text-white rounded-br-md"
                        : "bg-gradient-to-r from-gray-800/60 to-gray-900/60 border border-gray-700/50 text-gray-100 rounded-bl-md"
                    } backdrop-blur-sm shadow-lg`}
                  >
                    {message.image && (
                      <div className="mb-2">
                        <img
                          src={message.image}
                          alt="Attachment"
                          className="max-w-[250px] rounded-xl border border-gray-600/30 cursor-pointer hover:opacity-80 transition-opacity duration-200"
                          onClick={() => openImagePreview(message.image)}
                        />
                      </div>
                    )}
                    {message.text && (
                      <p className="text-sm leading-relaxed break-words whitespace-pre-wrap">
                        {message.text}
                      </p>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}

          {/* Scroll anchor */}
          <div ref={messageEndRef} className="h-1" />
        </div>
      </div>

      <MessageInput />

      {/* Image Preview Modal */}
      {previewImage && (
        <div
          className={`fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50 transition-all duration-200 ${
            isModalClosing ? "opacity-0" : "opacity-100"
          }`}
          onClick={closeImagePreview}
        >
          <div
            className={`relative max-w-4xl max-h-[90vh] p-6 transition-all duration-200 ${
              isModalClosing ? "scale-95 opacity-0" : "scale-100 opacity-100"
            }`}
          >
            <button
              onClick={closeImagePreview}
              className="absolute top-2 right-2 w-10 h-10 rounded-full bg-gray-800/90 hover:bg-gray-700/90 flex items-center justify-center transition-colors duration-200 shadow-lg z-10"
            >
              <X className="w-5 h-5 text-white" />
            </button>
            <img
              src={previewImage}
              alt="Preview"
              className="max-w-full max-h-full rounded-xl shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default ChatContainer;