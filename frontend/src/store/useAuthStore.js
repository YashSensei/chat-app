import { create } from "zustand";
import { axiosInstance } from "../lib/axios.js";
import toast from "react-hot-toast";

// WS_URL for WebSocket connection
const WS_URL = import.meta.env.MODE === "development" ? "ws://localhost:5000/ws" : "wss://your-production-domain.com/ws"; // Use wss:// for production HTTPS

export const useAuthStore = create((set, get) => ({
  authUser: null,
  isSigningUp: false,
  isLoggingIn: false,
  isUpdatingProfile: false,
  isCheckingAuth: true,
  onlineUsers: [],
  socket: null, // This will now hold a native WebSocket object

  checkAuth: async () => {
    try {
      const res = await axiosInstance.get("/auth/check");
      set({ authUser: res.data });
      get().connectSocket(); // Connect WebSocket after successful auth check
    } catch (error) {
      console.log("Error in checkAuth:", error);
      set({ authUser: null });
    } finally {
      set({ isCheckingAuth: false });
    }
  },

  signup: async (data) => {
    set({ isSigningUp: true });
    try {
      const res = await axiosInstance.post("/auth/signup", data);
      set({ authUser: res.data });
      toast.success("Account created successfully");
      get().connectSocket(); // Connect WebSocket after successful signup
    } catch (error) {
      toast.error(error.response.data.message);
    } finally {
      set({ isSigningUp: false });
    }
  },

  login: async (data) => {
    set({ isLoggingIn: true });
    try {
      const res = await axiosInstance.post("/auth/login", data);
      set({ authUser: res.data });
      toast.success("Logged in successfully");
      get().connectSocket(); // Connect WebSocket after successful login
    } catch (error) {
      toast.error(error.response.data.message);
    } finally {
      set({ isLoggingIn: false });
    }
  },

  logout: async () => {
    try {
      await axiosInstance.post("/auth/logout");
      set({ authUser: null });
      toast.success("Logged out successfully");
      get().disconnectSocket(); // Disconnect WebSocket on logout
    } catch (error) {
      toast.error(error.response.data.message);
    }
  },

  updateProfile: async (data) => {
    set({ isUpdatingProfile: true });
    try {
      const res = await axiosInstance.put("/auth/update-profile", data);
      set({ authUser: res.data });
      toast.success("Profile updated successfully");
    } catch (error) {
      console.log("error in update profile:", error);
      toast.error(error.response.data.message);
    } finally {
      set({ isUpdatingProfile: false });
    }
  },

  connectSocket: () => {
    const { authUser } = get();
    // Only connect if authUser exists and socket is not already open
    if (!authUser || (get().socket && get().socket.readyState === WebSocket.OPEN)) {
      return;
    }

    // Create a new native WebSocket connection
    const socket = new WebSocket(WS_URL);

    // Event listener for when the WebSocket connection is established
    socket.onopen = () => {
      console.log("WebSocket connected to Go backend!");
    };

    // Event listener for incoming messages from the WebSocket
    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data); // Parse the JSON message from Go backend
        console.log("Received WebSocket message:", data);

        // MODIFIED: Read payload for getOnlineUsers.
        // Removed the `else` block as `useChatStore` now handles "newMessage" directly.
        if (data.event === "getOnlineUsers") {
          set({ onlineUsers: data.payload });
        }
        // No `else if (data.event === "newMessage")` here.
        // `useChatStore`'s `subscribeToMessages` will handle "newMessage" events directly
        // by listening to the same `socket` instance.
      } catch (e) {
        console.error("Error parsing WebSocket message:", e, event.data);
      }
    };

    // Event listener for WebSocket errors
    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    // Event listener for when the WebSocket connection is closed
    socket.onclose = (event) => {
      console.log("WebSocket disconnected:", event.code, event.reason);
      set({ socket: null, onlineUsers: [] }); // Clear socket and online users on close
    };

    set({ socket: socket }); // Store the native WebSocket object in state
  },

  disconnectSocket: () => {
    if (get().socket && get().socket.readyState === WebSocket.OPEN) {
      get().socket.close(); // Close the native WebSocket connection
    }
    set({ socket: null, onlineUsers: [] }); // Clear state
  },
}));
