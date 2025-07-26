import { useState } from "react";
import { useAuthStore } from "../store/useAuthStore";
import {
  Eye,
  EyeOff,
  Loader2,
  Lock,
  Mail,
  MessageSquare,
  User,
} from "lucide-react";
import { Link } from "react-router-dom";
import toast from "react-hot-toast";

const SignUpPage = () => {
  const [showPassword, setShowPassword] = useState(false);
  const [formData, setFormData] = useState({
    fullName: "",
    email: "",
    password: "",
  });

  const { signup, isSigningUp } = useAuthStore();

  const validateForm = () => {
    if (!formData.fullName.trim()) return toast.error("Full name is required");
    if (!formData.email.trim()) return toast.error("Email is required");
    if (!/\S+@\S+\.\S+/.test(formData.email))
      return toast.error("Invalid email format");
    if (!formData.password) return toast.error("Password is required");
    if (formData.password.length < 6)
      return toast.error("Password must be at least 6 characters");

    return true;
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    const success = validateForm();

    if (success === true) signup(formData);
  };

  return (
    <div className="min-h-screen pt-16 grid lg:grid-cols-2 bg-gradient-to-br from-black via-gray-900 to-black relative overflow-hidden">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gray-800 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-blob"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gray-700 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-blob animation-delay-2000"></div>
        <div className="absolute top-40 left-40 w-80 h-80 bg-gray-600 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-blob animation-delay-4000"></div>
      </div>

      {/* Left side - Form */}
      <div className="flex flex-col justify-center items-center p-6 sm:p-12 relative z-10">
        <div className="w-full max-w-md">
          {/* Main card */}
          <div className="backdrop-blur-xl bg-black/40 border border-gray-800/50 rounded-3xl shadow-2xl p-8 space-y-8">
            {/* Logo and header */}
            <div className="text-center space-y-4">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-r from-gray-900 to-gray-800 shadow-lg">
                <MessageSquare className="w-8 h-8 text-white" />
              </div>
              <div className="space-y-2">
                <h1 className="text-3xl font-bold bg-gradient-to-r from-white via-gray-300 to-gray-400 bg-clip-text text-transparent">
                  Create Account
                </h1>
                <p className="text-gray-400 text-sm font-medium">
                  Get started with your free account
                </p>
              </div>
            </div>

            {/* Form */}
            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Full Name field */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-gray-300 ml-1">
                  Full Name
                </label>
                <div className="relative group">
                  <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors duration-200 group-focus-within:text-gray-300">
                    <User className="w-5 h-5 text-gray-500 group-focus-within:text-gray-300" />
                  </div>
                  <input
                    type="text"
                    className="w-full pl-12 pr-4 py-4 bg-black/20 border border-gray-700/50 rounded-2xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-gray-600/50 focus:border-gray-600/50 transition-all duration-300 hover:bg-black/30"
                    placeholder="Enter your full name"
                    value={formData.fullName}
                    onChange={(e) =>
                      setFormData({ ...formData, fullName: e.target.value })
                    }
                    required
                  />
                </div>
              </div>

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
                    placeholder="Create a password"
                    value={formData.password}
                    onChange={(e) =>
                      setFormData({ ...formData, password: e.target.value })
                    }
                    required
                  />
                  <button
                    type="button"
                    className="absolute inset-y-0 right-0 pr-4 flex items-center text-gray-500 hover:text-gray-300 transition-colors duration-200"
                    onClick={() => setShowPassword(!showPassword)}
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
                disabled={isSigningUp}
              >
                {isSigningUp ? (
                  <div className="flex items-center justify-center space-x-2">
                    <Loader2 className="w-5 h-5 animate-spin" />
                    <span>Creating account...</span>
                  </div>
                ) : (
                  "Create Account"
                )}
              </button>
            </form>

            {/* Footer link */}
            <div className="text-center pt-2">
              <p className="text-gray-400 text-sm">
                Already have an account?{" "}
                <Link
                  to="/login"
                  className="text-gray-300 hover:text-white font-medium transition-colors duration-200 hover:underline"
                >
                  Sign in
                </Link>
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Right side - Community Section */}
      <div className="hidden lg:flex items-center justify-center p-12 relative z-10">
        <div className="max-w-md text-center space-y-8">
          {/* Animated grid pattern */}
          <div className="grid grid-cols-3 gap-4 mb-8">
            {[...Array(9)].map((_, i) => (
              <div
                key={i}
                className={`aspect-square rounded-2xl bg-gradient-to-br from-gray-800/40 to-gray-900/40 border border-gray-700/30 backdrop-blur-sm ${
                  i % 2 === 0 ? "animate-pulse" : ""
                }`}
                style={{
                  animationDelay: `${i * 0.2}s`,
                  animationDuration: "3s",
                }}
              />
            ))}
          </div>

          {/* Community content */}
          <div className="space-y-6">
            <h2 className="text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-gray-300 bg-clip-text text-transparent">
              Join our community
            </h2>
            <p className="text-gray-400 text-lg leading-relaxed">
              Connect with friends, share moments, and stay in touch with your
              loved ones.
            </p>

            {/* Feature highlights */}
            <div className="space-y-3 pt-4">
              <div className="flex items-center gap-3 text-gray-300">
                <div className="w-2 h-2 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
                <span className="text-sm">Real-time messaging</span>
              </div>
              <div className="flex items-center gap-3 text-gray-300">
                <div className="w-2 h-2 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
                <span className="text-sm">Secure conversations</span>
              </div>
              <div className="flex items-center gap-3 text-gray-300">
                <div className="w-2 h-2 rounded-full bg-gradient-to-r from-gray-600 to-gray-500"></div>
                <span className="text-sm">Connect with anyone</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SignUpPage;