# 💬 Real-Time Chat Application

A modern, full-stack real-time chat application built with **Go** backend and **React** frontend, featuring WebSocket communication, JWT authentication, and cloud-based image storage.

![Chat App](frontend/public/screenshot-for-readme.png)

## 🚀 Features

### Core Functionality
- 🔐 **Secure Authentication** - JWT-based authentication with HTTP-only cookies
- 💬 **Real-Time Messaging** - Instant message delivery using native WebSockets
- 👥 **Online User Tracking** - See who's online in real-time
- 🖼️ **Image Sharing** - Upload and share images in chats (Cloudinary integration)
- 📱 **Responsive Design** - Beautiful UI that works on desktop and mobile
- 🎨 **Theme Switching** - Multiple theme options with DaisyUI
- 👤 **Profile Management** - Update profile pictures and user information

### Technical Features
- ⚡ **Fast Performance** - Go backend for high-performance message handling
- 🔄 **Persistent Storage** - MongoDB for data persistence
- 🛡️ **Security** - Password hashing with bcrypt, secure cookie handling
- 🌐 **CORS Support** - Configured for frontend-backend communication
- 📦 **Modular Architecture** - Clean code structure following Go best practices

## 🏗️ Architecture

### Backend (Go)
```
go-backend/
├── cmd/api/              # Application entry point
│   └── main.go
├── config/               # Configuration management
│   └── config.go
├── internal/             # Private application code
│   ├── auth/             # Authentication handlers & middleware
│   │   ├── handler.go
│   │   └── middleware.go
│   ├── chat/             # Chat/messaging handlers
│   │   └── handler.go
│   ├── models/           # Data models
│   │   ├── user.go
│   │   └── message.go
│   └── server/           # Server setup & routing
│       └── server.go
└── pkg/                  # Public/shared packages
    ├── db/               # Database connection
    │   └── mongodb.go
    └── utils/            # Utility functions
        ├── jwt.go        # JWT token generation/validation
        ├── cloudinary.go # Image upload service
        └── socket.go     # WebSocket hub & handlers
```

### Frontend (React)
```
frontend/
├── src/
│   ├── components/       # Reusable UI components
│   │   ├── ChatContainer.jsx
│   │   ├── Sidebar.jsx
│   │   ├── Navbar.jsx
│   │   └── ...
│   ├── pages/            # Page components
│   │   ├── HomePage.jsx
│   │   ├── LoginPage.jsx
│   │   ├── SignUpPage.jsx
│   │   └── ...
│   ├── store/            # State management (Zustand)
│   │   ├── useAuthStore.js
│   │   ├── useChatStore.js
│   │   └── useThemeStore.js
│   └── lib/              # Utilities & configurations
│       ├── axios.js
│       └── utils.js
```

## 🛠️ Tech Stack

### Backend
- **Go 1.23.4** - High-performance backend language
- **Gin** - Fast HTTP web framework
- **MongoDB** - NoSQL database with MongoDB Go Driver
- **JWT** - Token-based authentication
- **WebSockets** - Gorilla WebSocket for real-time communication
- **Cloudinary** - Cloud-based image storage and management
- **Bcrypt** - Password hashing

### Frontend
- **React 18** - Modern UI library
- **Vite** - Fast build tool and dev server
- **Zustand** - Lightweight state management
- **React Router** - Client-side routing
- **Axios** - HTTP client
- **Native WebSocket API** - Real-time communication
- **TailwindCSS** - Utility-first CSS framework
- **DaisyUI** - Component library for Tailwind
- **Lucide React** - Beautiful icon library
- **React Hot Toast** - Toast notifications

## 📋 Prerequisites

- **Go 1.23+** installed
- **Node.js 18+** and npm/yarn installed
- **MongoDB Atlas** account (or local MongoDB instance)
- **Cloudinary** account for image uploads

## ⚙️ Setup & Installation

### 1. Clone the Repository
```bash
git clone https://github.com/YashSensei/chat-app.git
cd chat-app
```

### 2. Backend Setup

#### Navigate to Backend Directory
```bash
cd go-backend
```

#### Create `.env` File
Create a `.env` file in the `go-backend` directory:

```env
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/chat-db?retryWrites=true&w=majority
PORT=5000
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
NODE_ENV=development

CLOUDINARY_CLOUD_NAME=your-cloudinary-cloud-name
CLOUDINARY_API_KEY=your-cloudinary-api-key
CLOUDINARY_API_SECRET=your-cloudinary-api-secret
```

**MongoDB Setup:**
1. Create a free account at [MongoDB Atlas](https://www.mongodb.com/atlas)
2. Create a new cluster
3. Create a database user with read/write permissions
4. Whitelist your IP address (or use 0.0.0.0/0 for development)
5. Get your connection string and replace in `.env`

**Cloudinary Setup:**
1. Create a free account at [Cloudinary](https://cloudinary.com/)
2. Get your Cloud Name, API Key, and API Secret from the dashboard
3. Add them to your `.env` file

#### Install Dependencies
```bash
go mod download
```

#### Run Backend Server
```bash
go run cmd/api/main.go
```

The backend server will start on `http://localhost:5000`

### 3. Frontend Setup

#### Navigate to Frontend Directory
```bash
cd ../frontend
```

#### Install Dependencies
```bash
npm install
```

#### Run Development Server
```bash
npm run dev
```

The frontend will start on `http://localhost:5173`

## 🚀 Running the Application

### Development Mode

1. **Start Backend** (Terminal 1):
```bash
cd go-backend
go run cmd/api/main.go
```

2. **Start Frontend** (Terminal 2):
```bash
cd frontend
npm run dev
```

3. Open your browser and navigate to `http://localhost:5173`

### Production Build

#### Build Frontend
```bash
cd frontend
npm run build
```

#### Run Production Server
The Go backend will automatically serve the built frontend from `frontend/dist` when `NODE_ENV=production`.

```bash
cd go-backend
NODE_ENV=production go run cmd/api/main.go
```

## 📡 API Endpoints

### Authentication
- `POST /api/auth/signup` - Register new user
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/check` - Check auth status (protected)
- `PUT /api/auth/update-profile` - Update profile (protected)

### Messages
- `GET /api/messages/users` - Get all users for sidebar (protected)
- `GET /api/messages/:id` - Get messages with specific user (protected)
- `POST /api/messages/send/:id` - Send message to user (protected)

### WebSocket
- `GET /ws` - WebSocket connection endpoint (protected)

## 🔒 Security Features

- **Password Hashing** - Bcrypt with default cost factor (10)
- **JWT Tokens** - HTTP-only cookies with 7-day expiration
- **CORS Protection** - Configured for specific origin
- **Authentication Middleware** - Protects sensitive routes
- **Input Validation** - Request body validation with Gin bindings
- **Secure Cookies** - HttpOnly and Secure flags in production

## 🌐 WebSocket Communication

### Connection Flow
1. User authenticates via HTTP (receives JWT cookie)
2. Frontend establishes WebSocket connection to `/ws`
3. Backend validates JWT from cookie before upgrading connection
4. Client is registered in WebSocket hub with their user ID
5. Real-time events are sent/received through this connection

### WebSocket Events

#### Sent by Server
```javascript
// Online users list
{
  "event": "getOnlineUsers",
  "payload": ["userId1", "userId2", ...]
}

// New message
{
  "event": "newMessage",
  "payload": {
    "ID": "messageId",
    "SenderID": "senderId",
    "ReceiverID": "receiverId",
    "Text": "message text",
    "Image": "image url",
    "CreatedAt": "timestamp",
    "UpdatedAt": "timestamp"
  }
}
```

## 🎨 Frontend State Management

### Zustand Stores

#### `useAuthStore`
- User authentication state
- WebSocket connection management
- Online users tracking
- Auth actions (signup, login, logout, updateProfile)

#### `useChatStore`
- Messages and users list
- Selected user state
- Message actions (getUsers, getMessages, sendMessage)
- WebSocket message subscription

#### `useThemeStore`
- Theme selection and persistence
- DaisyUI theme switching

## 🧪 Testing

### Manual Testing Checklist
- [ ] User registration with valid data
- [ ] User login with valid credentials
- [ ] Profile picture update
- [ ] Send text message
- [ ] Send image message
- [ ] Receive real-time messages
- [ ] Online user status updates
- [ ] Theme switching
- [ ] Logout functionality

## 🐛 Troubleshooting

### Backend Issues

**MongoDB Connection Failed**
- Verify MongoDB URI is correct
- Check if IP is whitelisted in MongoDB Atlas
- Ensure database user has proper permissions

**Cloudinary Upload Failed**
- Verify Cloudinary credentials
- Check API key permissions
- Ensure base64 image format is correct

**WebSocket Connection Failed**
- Check CORS settings
- Verify JWT token is valid
- Ensure WebSocket hub is initialized

### Frontend Issues

**Axios 401 Errors**
- Check if backend is running
- Verify `withCredentials: true` in axios config
- Clear browser cookies and retry

**WebSocket Not Connecting**
- Verify WebSocket URL matches backend port
- Check browser console for errors
- Ensure user is authenticated first

## 📝 Environment Variables

### Backend (`.env` in `go-backend/`)
| Variable | Description | Example |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection string | `mongodb+srv://user:pass@cluster.mongodb.net/chat-db` |
| `PORT` | Server port | `5000` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key` |
| `NODE_ENV` | Environment mode | `development` or `production` |
| `CLOUDINARY_CLOUD_NAME` | Cloudinary cloud name | `your-cloud-name` |
| `CLOUDINARY_API_KEY` | Cloudinary API key | `123456789012345` |
| `CLOUDINARY_API_SECRET` | Cloudinary API secret | `your-api-secret` |

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Yash Agrawal**
- GitHub: [@YashSensei](https://github.com/YashSensei)

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)
- [Cloudinary](https://cloudinary.com/)
- [React](https://react.dev/)
- [TailwindCSS](https://tailwindcss.com/)
- [DaisyUI](https://daisyui.com/)

---

⭐ If you found this project helpful, please consider giving it a star!