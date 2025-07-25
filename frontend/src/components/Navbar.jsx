import { Link, useLocation } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";
import { LogOut, MessageSquare, User } from "lucide-react";

const Navbar = () => {
  const { logout, authUser } = useAuthStore();
  const location = useLocation();
  const isAuthPage =
    location.pathname === "/login" || location.pathname === "/signup";

  return (
    <header
      className={`fixed w-full top-0 z-50 transition-all duration-300 ${
        isAuthPage
          ? "bg-black/20 backdrop-blur-xl border-b border-gray-800/30"
          : "bg-gradient-to-r from-black via-gray-900 to-black border-b border-gray-800/50"
      } shadow-2xl`}
    >
      <div className="container mx-auto px-6 h-16 flex items-center justify-between">
        {/* Logo */}
        <Link
          to="/"
          className="flex items-center gap-3 hover:opacity-80 transition-all duration-300 group"
        >
          <div className="w-10 h-10 rounded-xl bg-gradient-to-r from-gray-800 to-gray-900 flex items-center justify-center shadow-lg group-hover:shadow-gray-500/25 transition-all duration-300">
            <MessageSquare className="w-5 h-5 text-white" />
          </div>
          <h1 className="text-xl font-bold tracking-tight bg-gradient-to-r from-white via-gray-200 to-gray-300 bg-clip-text text-transparent">
            ChatApp
          </h1>
        </Link>

        {/* Navigation */}
        <div className="flex items-center gap-3">
          {!isAuthPage && (
            <Link
              to="/"
              className="px-4 py-2 rounded-xl bg-gray-900/50 border border-gray-700/50 text-gray-300 hover:bg-gray-800/50 hover:text-white hover:border-gray-600/50 transition-all duration-300 flex items-center gap-2 backdrop-blur-sm"
            >
              <MessageSquare className="w-4 h-4" />
              <span className="hidden sm:inline font-medium">Chat</span>
            </Link>
          )}

          {authUser && !isAuthPage && (
            <>
              <Link
                to="/profile"
                className="px-4 py-2 rounded-xl bg-gray-900/50 border border-gray-700/50 text-gray-300 hover:bg-gray-800/50 hover:text-white hover:border-gray-600/50 transition-all duration-300 flex items-center gap-2 backdrop-blur-sm"
              >
                <User className="w-4 h-4" />
                <span className="hidden sm:inline font-medium">Profile</span>
              </Link>
              <button
                className="px-4 py-2 rounded-xl bg-red-900/30 border border-red-800/50 text-red-300 hover:bg-red-800/40 hover:text-red-200 hover:border-red-700/50 transition-all duration-300 flex items-center gap-2 backdrop-blur-sm"
                onClick={logout}
              >
                <LogOut className="w-4 h-4" />
                <span className="hidden sm:inline font-medium">Logout</span>
              </button>
            </>
          )}

          {/* Auth page navigation */}
          {isAuthPage && (
            <div className="flex items-center gap-3">
              {location.pathname === "/login" ? (
                <Link
                  to="/signup"
                  className="px-4 py-2 rounded-xl bg-white/5 border border-gray-700/50 text-gray-300 hover:bg-white/10 hover:text-white hover:border-gray-600/50 transition-all duration-300 backdrop-blur-sm font-medium"
                >
                  Sign Up
                </Link>
              ) : (
                <Link
                  to="/login"
                  className="px-4 py-2 rounded-xl bg-white/5 border border-gray-700/50 text-gray-300 hover:bg-white/10 hover:text-white hover:border-gray-600/50 transition-all duration-300 backdrop-blur-sm font-medium"
                >
                  Sign In
                </Link>
              )}
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Navbar;
