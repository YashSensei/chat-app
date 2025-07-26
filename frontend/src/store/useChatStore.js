import { create } from "zustand";
import toast from "react-hot-toast";
import { axiosInstance } from "../lib/axios";
import { useAuthStore } from "./useAuthStore"; // Import useAuthStore to get the socket instance

export const useChatStore = create((set, get) => ({
  messages: [],
  users: [],
  selectedUser: null,
  isUsersLoading: false,
  isMessagesLoading: false,

  getUsers: async () => {
    set({ isUsersLoading: true });
    try {
      const res = await axiosInstance.get("/messages/users");
      set({ users: res.data });
    } catch (error) {
      toast.error(error.response.data.message);
    } finally {
      set({ isUsersLoading: false });
    }
  },

  getMessages: async (userId) => {
    set({ isMessagesLoading: true });
    try {
      const res = await axiosInstance.get(`/messages/${userId}`);
      set({ messages: res.data });
    } catch (error) {
      toast.error(error.response.data.message);
    } finally {
      set({ isMessagesLoading: false });
    }
  },

  sendMessage: async (messageData) => {
    const { selectedUser, messages } = get();
    if (!selectedUser) {
      toast.error("Please select a user to send a message.");
      return;
    }
    try {
      // Send message via HTTP POST (this is correct)
      const res = await axiosInstance.post(`/messages/send/${selectedUser._id}`, messageData);
      // Add the new message to the local state immediately
      set({ messages: [...messages, res.data] });
    } catch (error) {
      console.error("Error sending message:", error);
      toast.error(error.response.data.message || "Failed to send message.");
    }
  },

  // ADDED: addMessage function to update messages state from WebSocket
  addMessage: (receivedWsMessage) => { // Renamed parameter for clarity
    const { selectedUser } = get();
    const authUser = useAuthStore.getState().authUser;

    // Extract the actual message payload from the WebSocket message
    const rawNewMessage = receivedWsMessage.payload;

    // Normalize the message structure to match HTTP fetched messages
    const normalizedMessage = {
      _id: rawNewMessage.ID,
      senderId: rawNewMessage.SenderID,
      receiverId: rawNewMessage.ReceiverID,
      text: rawNewMessage.Text,
      image: rawNewMessage.Image,
      createdAt: rawNewMessage.CreatedAt,
      updatedAt: rawNewMessage.UpdatedAt,
      // Add any other fields if necessary, ensuring consistent naming
    };

    // Ensure selectedUser and authUser are available before comparison
    if (!selectedUser || !authUser) {
        console.log("addMessage: selectedUser or authUser is not set. Skipping display.");
        return;
    }

    // A message should be added to the current chat if:
    // 1. The message is from the currently selected user AND is for the authenticated user (incoming).
    // 2. The message is sent by the authenticated user AND is for the currently selected user (outgoing, for sender's own UI).
    // Use normalizedMessage properties for comparison
    const isMessageForCurrentChat =
        (normalizedMessage.senderId === selectedUser._id && normalizedMessage.receiverId === authUser._id) ||
        (normalizedMessage.senderId === authUser._id && normalizedMessage.receiverId === selectedUser._id);

    if (isMessageForCurrentChat) {
      set((state) => ({
        messages: [...state.messages, normalizedMessage] // Add the normalized message
      }));
    } else {
      console.log("Message not for current chat or user. Skipping display (not for current chat or user).");
    }
  },

  // MODIFIED: subscribeToMessages to use native WebSocket addEventListener
  subscribeToMessages: () => {
    const { selectedUser } = get();
    if (!selectedUser) return;

    // Get the native WebSocket instance from useAuthStore
    const socket = useAuthStore.getState().socket;

    if (!socket || socket.readyState !== WebSocket.OPEN) {
      console.warn("WebSocket not connected or not open. Cannot subscribe to messages.");
      return;
    }

    // Define the message handler function
    const messageHandler = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.event === "newMessage") {
          // Pass the entire WebSocket message object (with event and payload) to addMessage
          get().addMessage(data);
        }
      } catch (e) {
        console.error("Error parsing WebSocket message in useChatStore:", e, event.data);
      }
    };

    // Attach the event listener for "message" events
    socket.addEventListener("message", messageHandler);

    // Store the handler so we can remove it later
    set({ _messageHandler: messageHandler }); // Using a private-like key
  },

  // MODIFIED: unsubscribeFromMessages to use native WebSocket removeEventListener
  unsubscribeFromMessages: () => {
    const socket = useAuthStore.getState().socket;
    const messageHandler = get()._messageHandler; // Retrieve the stored handler

    if (socket && messageHandler) {
      socket.removeEventListener("message", messageHandler);
      set({ _messageHandler: null }); // Clear the stored handler
    }
  },

  setSelectedUser: (selectedUser) => set({ selectedUser }),
}));
