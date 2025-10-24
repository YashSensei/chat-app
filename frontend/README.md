# Chat App - Frontend

This is the React frontend for the Chat App (Vite + React).

What this frontend does:
- Uses Axios to call backend REST APIs for authentication, user lists and messages.
- Uses native WebSocket API to receive real-time events (online users and incoming messages).

Deployed link: (add your frontend deployed URL here)

HTTP requests used by the frontend (examples):
- POST /api/auth/signup — body: { fullName, email, password }
- POST /api/auth/login — body: { email, password }
- POST /api/auth/logout — no body
- GET /api/auth/check — no body (returns current user)
- PUT /api/auth/update-profile — body: { profilePic: base64String }
- GET /api/messages/users — no body (returns users list)
- GET /api/messages/:id — no body (returns message history)
- POST /api/messages/send/:id — body: { text?, image? }

WebSocket
- Connect to: ws://localhost:5000/ws (development) or wss://<your-domain>/ws (production)
- The backend authenticates the connection using the JWT cookie.
- Events: { event: "getOnlineUsers", payload: [...] } and { event: "newMessage", payload: { ... } }

How to run (development):
```
cd frontend
npm install
npm run dev
```

Environment sample: `frontend/.env.sample` (API base URL and WS URL)

If you want additional examples of request/response shapes or to remove/add endpoints, tell me which ones and I will update this file.
- 📱 **Fully Responsive** - Optimized for both mobile and desktop
- 🎭 **Multiple Themes** - Theme switcher with 10+ DaisyUI themes
- 🌓 **Dark Mode Support** - Beautiful dark theme options
- ⚡ **Smooth Animations** - Fluid transitions and loading states
- 🖼️ **Image Preview** - Full-screen image preview modal

### Functionality
- 💬 **Real-Time Messaging** - Instant message delivery via WebSockets
- 👥 **Online Status** - Live online/offline user indicators
- 📸 **Image Sharing** - Upload and share images in chats
- 🔍 **Online Filter** - Filter to show only online users
- 📝 **Message History** - Persistent chat history
- 🔔 **Toast Notifications** - User-friendly notifications for actions
- 🔐 **Secure Authentication** - JWT-based auth with HTTP-only cookies

## 🛠️ Tech Stack

### Core
- **React 18.3.1** - Modern UI library with hooks
- **Vite 5.4.10** - Next-generation frontend tooling
- **React Router DOM 6.28.0** - Client-side routing

### State Management
- **Zustand 5.0.1** - Lightweight state management
  - `useAuthStore` - Authentication & WebSocket connection
  - `useChatStore` - Messages & users
  - `useThemeStore` - Theme preferences

### UI & Styling
- **TailwindCSS 3.4.15** - Utility-first CSS framework
- **DaisyUI 4.12.14** - Beautiful component library
- **Lucide React 0.459.0** - Modern icon library

### Communication
- **Axios 1.7.7** - HTTP client with interceptors
- **Native WebSocket API** - Real-time bidirectional communication

### User Experience
- **React Hot Toast 2.4.1** - Beautiful toast notifications

### Development
- **ESLint** - Code linting and formatting
- **PostCSS** - CSS transformations
- **Autoprefixer** - Automatic vendor prefixing

## 📁 Project Structure

```
frontend/
├── public/                 # Static assets
│   ├── avatar.png         # Default avatar image
│   ├── screenshot-for-readme.png
│   └── vite.svg
├── src/
│   ├── components/        # React components
│   │   ├── AuthImagePattern.jsx    # Auth page background pattern
│   │   ├── ChatContainer.jsx       # Main chat interface
│   │   ├── ChatHeader.jsx          # Chat header with user info
│   │   ├── MessageInput.jsx        # Message input with image upload
│   │   ├── Navbar.jsx              # Top navigation bar
│   │   ├── NoChatSelected.jsx      # Welcome screen
│   │   ├── Sidebar.jsx             # User list sidebar
│   │   └── skeletons/              # Loading skeletons
│   │       ├── MessageSkeleton.jsx
│   │       └── SidebarSkeleton.jsx
│   ├── constants/         # App constants
│   │   └── index.js       # Theme options & constants
│   ├── lib/               # Utilities & configurations
│   │   ├── axios.js       # Axios instance configuration
│   │   └── utils.js       # Helper functions
│   ├── pages/             # Page components
│   │   ├── HomePage.jsx   # Main chat page
│   │   ├── LoginPage.jsx  # User login
│   │   ├── SignUpPage.jsx # User registration
│   │   ├── ProfilePage.jsx # User profile management
│   │   └── SettingsPage.jsx # App settings
│   ├── store/             # Zustand stores
│   │   ├── useAuthStore.js   # Authentication state
│   │   ├── useChatStore.js   # Chat state
│   │   └── useThemeStore.js  # Theme state
│   ├── App.jsx            # Root component
│   ├── index.css          # Global styles
│   └── main.jsx           # App entry point
├── eslint.config.js       # ESLint configuration
├── index.html             # HTML template
├── package.json           # Dependencies & scripts
├── postcss.config.js      # PostCSS configuration
├── tailwind.config.js     # Tailwind configuration
└── vite.config.js         # Vite configuration
```

## 🔧 Component Details

### Core Components

#### `ChatContainer.jsx`
- Main chat interface with message history
- Auto-scrolling to latest message
- Image preview modal
- Message bubbles with sender/receiver styling
- Loading skeletons during fetch

#### `Sidebar.jsx`
- User list with profile pictures
- Online/offline status indicators
- Filter toggle for online users
- Selected user highlighting
- Responsive design (full screen on mobile)

#### `MessageInput.jsx`
- Text message input with emoji support
- Image upload with preview
- Base64 encoding for images
- Send button with loading state

#### `Navbar.jsx`
- App branding and logo
- User profile dropdown
- Settings and logout buttons
- Responsive mobile menu

### Page Components

#### `HomePage.jsx`
- Container for Sidebar and ChatContainer
- Conditional rendering based on selected user
- Animated background elements
- Responsive layout switching

#### `LoginPage.jsx` & `SignUpPage.jsx`
- Form validation
- Error handling with toast notifications
- Background pattern component
- Loading states during submission

#### `ProfilePage.jsx`
- Profile picture upload
- User information display
- Avatar selection from predefined set

#### `SettingsPage.jsx`
- Theme selection grid
- Preview of theme colors
- Persistent theme storage

## 📦 State Management

### useAuthStore
```javascript
{
  authUser: null,           // Current user object
  isSigningUp: false,       // Loading state
  isLoggingIn: false,       // Loading state
  isUpdatingProfile: false, // Loading state
  isCheckingAuth: true,     // Initial auth check
  onlineUsers: [],          // Array of online user IDs
  socket: null,             // WebSocket instance
  
  // Actions
  checkAuth(),              // Verify authentication
  signup(data),             // Register new user
  login(data),              // Login user
  logout(),                 // Logout user
  updateProfile(data),      // Update profile
  connectSocket(),          // Establish WebSocket
  disconnectSocket()        // Close WebSocket
}
```

### useChatStore
```javascript
{
  messages: [],             // Message array
  users: [],                // All users list
  selectedUser: null,       // Currently selected chat user
  isUsersLoading: false,    // Loading state
  isMessagesLoading: false, // Loading state
  
  // Actions
  getUsers(),               // Fetch users list
  getMessages(userId),      // Fetch chat messages
  sendMessage(data),        // Send new message
  addMessage(message),      // Add WebSocket message
  subscribeToMessages(),    // Subscribe to WebSocket
  unsubscribeFromMessages(), // Unsubscribe from WebSocket
  setSelectedUser(user)     // Set active chat
}
```

### useThemeStore
```javascript
{
  theme: 'coffee',          // Current theme name
  
  // Actions
  setTheme(theme)           // Change and persist theme
}
```

## 🌐 API Integration

### Axios Configuration (`lib/axios.js`)
```javascript
export const axiosInstance = axios.create({
  baseURL: 'http://localhost:5000/api', // Development
  withCredentials: true,                 // Include cookies
});
```

### API Endpoints
- `POST /auth/signup` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `GET /auth/check` - Auth verification
- `PUT /auth/update-profile` - Profile update
- `GET /messages/users` - Get users list
- `GET /messages/:id` - Get messages with user
- `POST /messages/send/:id` - Send message

## 🔌 WebSocket Integration

### Connection Management
```javascript
// Connection established after authentication
const socket = new WebSocket('ws://localhost:5000/ws');

// Event handlers
socket.onopen = () => console.log('Connected');
socket.onmessage = (event) => handleMessage(event);
socket.onerror = (error) => console.error(error);
socket.onclose = () => console.log('Disconnected');
```

### WebSocket Events
```javascript
// Incoming: Online users update
{
  "event": "getOnlineUsers",
  "payload": ["userId1", "userId2"]
}

// Incoming: New message
{
  "event": "newMessage",
  "payload": {
    "ID": "messageId",
    "SenderID": "senderId",
    "ReceiverID": "receiverId",
    "Text": "Hello!",
    "Image": "",
    "CreatedAt": "2025-10-24T12:00:00Z",
    "UpdatedAt": "2025-10-24T12:00:00Z"
  }
}
```

## 🎨 Styling & Themes

### TailwindCSS Custom Classes
```css
.custom-scrollbar        /* Custom scrollbar styling */
.animate-blob            /* Blob animation for background */
.animation-delay-2000    /* Animation delay utilities */
.animation-delay-4000
```

### Available Themes (DaisyUI)
- light, dark, cupcake, bumblebee, emerald
- corporate, synthwave, retro, cyberpunk, valentine
- halloween, garden, forest, aqua, lofi
- pastel, fantasy, wireframe, black, luxury
- dracula, cmyk, autumn, business, acid
- lemonade, night, coffee, winter, dim
- nord, sunset

## 🚀 Getting Started

### Prerequisites
- Node.js 18+ and npm/yarn
- Backend server running on `http://localhost:5000`

### Installation

1. **Install Dependencies**
```bash
npm install
```

2. **Run Development Server**
```bash
npm run dev
```

3. **Build for Production**
```bash
npm run build
```

4. **Preview Production Build**
```bash
npm run preview
```

### Available Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start Vite dev server with HMR |
| `npm run build` | Build production bundle |
| `npm run lint` | Run ESLint checks |
| `npm run preview` | Preview production build locally |

## 🔧 Configuration Files

### `vite.config.js`
```javascript
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:5000'
    }
  }
})
```

### `tailwind.config.js`
```javascript
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: { extend: {} },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark", "cupcake", ...],
  },
}
```

## 📱 Responsive Design

### Breakpoints
- **Mobile**: < 768px (Full screen chat, sidebar overlay)
- **Tablet**: 768px - 1024px (Sidebar visible, optimized layout)
- **Desktop**: > 1024px (Full sidebar + chat side-by-side)

### Mobile Optimizations
- Touch-friendly buttons and inputs
- Optimized image loading
- Conditional rendering for small screens
- Swipe gestures for navigation

## 🐛 Troubleshooting

### Common Issues

**WebSocket Connection Failed**
- Ensure backend is running on port 5000
- Check CORS settings in backend
- Verify authentication is successful

**Images Not Uploading**
- Check Cloudinary configuration in backend
- Verify file size is under limit
- Ensure base64 encoding is correct

**Theme Not Persisting**
- Check localStorage is enabled
- Verify browser privacy settings
- Clear cache and retry

**Messages Not Showing**
- Check WebSocket connection status
- Verify selected user is set
- Check browser console for errors

## 🔒 Security Best Practices

- ✅ HTTP-only cookies for JWT storage
- ✅ XSS protection through React's built-in escaping
- ✅ CORS configured for specific origin
- ✅ Input validation before sending to backend
- ✅ Secure WebSocket connection in production
- ✅ No sensitive data in localStorage

## 📈 Performance Optimizations

- ⚡ React lazy loading for routes
- ⚡ Image optimization with Cloudinary
- ⚡ Virtual scrolling for long message lists
- ⚡ Debounced input handlers
- ⚡ Memoized components where needed
- ⚡ Optimized re-renders with Zustand

## 🧪 Testing Checklist

### Manual Testing
- [ ] User can sign up with valid data
- [ ] User can log in with credentials
- [ ] Profile picture updates successfully
- [ ] Messages send and receive in real-time
- [ ] Online status updates correctly
- [ ] Theme changes persist
- [ ] Images upload and display
- [ ] Mobile responsive layout works
- [ ] Logout clears session

## 📚 Learn More

### React + Vite
This template uses [@vitejs/plugin-react](https://github.com/vitejs/vite-plugin-react/blob/main/packages/plugin-react/README.md) with Babel for Fast Refresh.

### Resources
- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)
- [TailwindCSS Documentation](https://tailwindcss.com/)
- [DaisyUI Components](https://daisyui.com/components/)
- [Zustand Documentation](https://docs.pmnd.rs/zustand/)

## 👨‍💻 Development

### Code Style
- ESLint for code quality
- Prettier for formatting (recommended)
- Consistent component structure
- Descriptive variable names

### Best Practices
- Keep components small and focused
- Use custom hooks for shared logic
- Implement proper error boundaries
- Add loading states for async operations
- Write descriptive comments for complex logic

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is part of the main chat-app repository. See the main [LICENSE](../LICENSE) file for details.

---

Built with ❤️ using React + Vite
