import { useState } from "react";
import { useAuthStore } from "../store/useAuthStore";
import { Link } from "react-router-dom";
import { Eye, EyeOff, Loader2, Lock, Mail, MessageSquare } from "lucide-react";

const LoginPage = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const { login, isLoggingIn } = useAuthStore();

  const handleSubmit = async (e) => {
    e.preventDefault();
    login(formData);
  };

  return (
    <div className="min-h-screen pt-16 flex items-center justify-center relative overflow-hidden bg-gradient-to-br from-black via-gray-900 to-black">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gray-800 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-blob"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gray-700 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-blob animation-delay-2000"></div>
        <div className="absolute top-40 left-40 w-80 h-80 bg-gray-600 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-blob animation-delay-4000"></div>
      </div>

      {/* Main container */}
      <div className="relative z-10 w-full max-w-md px-6">
        {/* Main card */}
        <div className="backdrop-blur-xl bg-black/40 border border-gray-800/50 rounded-3xl shadow-2xl p-8 space-y-8">
          {/* Logo and header */}
          <div className="text-center space-y-4">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-r from-gray-900 to-gray-800 shadow-lg">
              <MessageSquare className="w-8 h-8 text-white" />
            </div>
            <div className="space-y-2">
              <h1 className="text-3xl font-bold bg-gradient-to-r from-white via-gray-300 to-gray-400 bg-clip-text text-transparent">
                Welcome Back
              </h1>
              <p className="text-gray-400 text-sm font-medium">
                Sign in to continue your conversation
              </p>
            </div>
          </div>

          {/* Form */}
        <form onSubmit={handleSubmit} className="space-y-6">
            {/* Email field */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-300 ml-1">
                Email Address
            </label>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors duration-200 group-focus-within:text-gray-300">
                  <Mail className="w-5 h-5 text-gray-500 group-focus-within:text-gray-300" />
              </div>
              <input
                type="email"
                  className="w-full pl-12 pr-4 py-4 bg-black/20 border border-gray-700/50 rounded-2xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-600/50 focus:border-gray-600/50 transition-all duration-300 hover:bg-black/30"
                  placeholder="Enter your email"
                value={formData.email}
                  onChange={(e) =>
                    setFormData({ ...formData, email: e.target.value })
                  }
                  required
              />
            </div>
          </div>

            {/* Password field */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-300 ml-1">
                Password
            </label>
              <div className="relative group">
                <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors duration-200">
                  <Lock className="w-5 h-5 text-gray-500 group-focus-within:text-gray-300" />
              </div>
              <input
                type={showPassword ? "text" : "password"}
                  className="w-full pl-12 pr-12 py-4 bg-black/20 border border-gray-700/50 rounded-2xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-600/50 focus:border-gray-600/50 transition-all duration-300 hover:bg-black/30"
                  placeholder="Enter your password"
                value={formData.password}
                  onChange={(e) =>
                    setFormData({ ...formData, password: e.target.value })
                  }
                  required
              />
              <button
                type="button"
                  className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-500 hover:text-gray-300 transition-colors duration-200"
                onClick={() => setShowPassword((prev) => !prev)}
                tabIndex={-1}
              >
                  {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                  ) : (
                    <Eye className="w-5 h-5" />
                  )}
              </button>
            </div>
          </div>

            {/* Submit button */}
          <button
            type="submit"
              className="w-full py-4 px-6 bg-gradient-to-r from-gray-900 to-black hover:from-gray-800 hover:to-gray-900 text-white font-semibold rounded-2xl shadow-lg focus:outline-none focus:ring-2 focus:ring-gray-600/50 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
            disabled={isLoggingIn}
          >
              {isLoggingIn ? (
                <div className="flex items-center justify-center space-x-2">
                  <Loader2 className="w-5 h-5 animate-spin" />
                  <span>Signing in...</span>
                </div>
              ) : (
                "Sign In"
              )}
          </button>
        </form>

          {/* Footer link */}
          <div className="text-center pt-2">
            <p className="text-gray-400 text-sm">
              Don't have an account?{" "}
              <Link
                to="/signup"
                className="text-gray-300 hover:text-white font-medium transition-colors duration-200 hover:underline"
              >
                Create one
          </Link>
            </p>
          </div>
        </div>

        {/* Subtle bottom text */}
        <div className="text-center mt-8">
          <p className="text-gray-600 text-xs">
            Secure login powered by modern encryption
          </p>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;