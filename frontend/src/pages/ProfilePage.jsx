import { useState } from "react";
import { useAuthStore } from "../store/useAuthStore";
import { Camera, Mail, User, Calendar, Shield } from "lucide-react";

const ProfilePage = () => {
  const { authUser, isUpdatingProfile, updateProfile } = useAuthStore();
  const [selectedImg, setSelectedImg] = useState(null);

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();

    reader.readAsDataURL(file);

    reader.onload = async () => {
      const base64Image = reader.result;
      setSelectedImg(base64Image);
      await updateProfile({ profilePic: base64Image });
    };
  };

  return (
    <div className="min-h-screen pt-16 bg-gradient-to-br from-black via-gray-900 to-black relative overflow-hidden">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-gray-800 rounded-full mix-blend-multiply filter blur-xl opacity-5 animate-blob"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-gray-700 rounded-full mix-blend-multiply filter blur-xl opacity-5 animate-blob animation-delay-2000"></div>
        <div className="absolute top-40 left-40 w-80 h-80 bg-gray-600 rounded-full mix-blend-multiply filter blur-xl opacity-5 animate-blob animation-delay-4000"></div>
      </div>

      <div className="relative z-10 max-w-4xl mx-auto p-6 py-8">
        <div className="backdrop-blur-xl bg-black/30 border border-gray-800/50 rounded-3xl shadow-2xl overflow-hidden">
          {/* Header Section */}
          <div className="p-8 border-b border-gray-800/50 bg-black/20">
            <div className="text-center space-y-2">
              <h1 className="text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-gray-300 bg-clip-text text-transparent">
                Profile
              </h1>
              <p className="text-gray-400 text-lg">Your profile information</p>
            </div>
          </div>

          {/* Avatar Section */}
          <div className="p-8 flex flex-col items-center space-y-6">
            <div className="relative group">
              <div className="absolute inset-0 bg-gradient-to-r from-gray-600 to-gray-700 rounded-full opacity-20 group-hover:opacity-30 transition-opacity duration-300 blur-lg"></div>
              <img
                src={selectedImg || authUser.profilePic || "/avatar.png"}
                alt="Profile"
                className="relative w-32 h-32 rounded-full object-cover border-4 border-gray-700/50 shadow-2xl group-hover:scale-105 transition-transform duration-300"
              />
              <label
                htmlFor="avatar-upload"
                className={`
                  absolute bottom-2 right-2 
                  bg-gradient-to-r from-gray-700 to-gray-800 hover:from-gray-600 hover:to-gray-700
                  p-3 rounded-full cursor-pointer 
                  border border-gray-600/50 shadow-lg
                  transition-all duration-300 hover:scale-110
                  ${
                    isUpdatingProfile ? "animate-pulse pointer-events-none" : ""
                  }
                `}
              >
                <Camera className="w-5 h-5 text-white" />
                <input
                  type="file"
                  id="avatar-upload"
                  className="hidden"
                  accept="image/*"
                  onChange={handleImageUpload}
                  disabled={isUpdatingProfile}
                />
              </label>
            </div>
            <p className="text-sm text-gray-400 text-center max-w-md">
              {isUpdatingProfile
                ? "Uploading your new profile picture..."
                : "Click the camera icon to update your photo"}
            </p>
          </div>

          {/* Profile Information */}
          <div className="p-8 space-y-6">
            <div className="grid md:grid-cols-2 gap-6">
              {/* Full Name */}
              <div className="space-y-3">
                <div className="flex items-center gap-3 text-gray-300">
                  <div className="p-2 rounded-xl bg-gradient-to-r from-gray-800 to-gray-900">
                    <User className="w-4 h-4" />
                  </div>
                  <span className="text-sm font-medium">Full Name</span>
                </div>
                <div className="px-4 py-3 bg-black/20 border border-gray-700/50 rounded-2xl text-white backdrop-blur-sm">
                  {authUser?.fullName}
                </div>
              </div>

              {/* Email Address */}
              <div className="space-y-3">
                <div className="flex items-center gap-3 text-gray-300">
                  <div className="p-2 rounded-xl bg-gradient-to-r from-gray-800 to-gray-900">
                    <Mail className="w-4 h-4" />
                  </div>
                  <span className="text-sm font-medium">Email Address</span>
                </div>
                <div className="px-4 py-3 bg-black/20 border border-gray-700/50 rounded-2xl text-white backdrop-blur-sm">
                  {authUser?.email}
                </div>
              </div>
            </div>

            {/* Account Information Card */}
            <div className="mt-8 backdrop-blur-xl bg-black/20 border border-gray-800/50 rounded-2xl p-6">
              <div className="flex items-center gap-3 mb-6">
                <div className="p-2 rounded-xl bg-gradient-to-r from-gray-800 to-gray-900">
                  <Shield className="w-5 h-5 text-white" />
                </div>
                <h2 className="text-xl font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                  Account Information
                </h2>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between py-3 border-b border-gray-700/30">
                  <div className="flex items-center gap-3 text-gray-300">
                    <Calendar className="w-4 h-4" />
                    <span className="font-medium">Member Since</span>
                  </div>
                  <span className="text-white font-semibold">
                    {new Date(authUser.createdAt).toLocaleDateString("en-US", {
                      year: "numeric",
                      month: "long",
                      day: "numeric",
                    })}
                  </span>
                </div>

                <div className="flex items-center justify-between py-3">
                  <div className="flex items-center gap-3 text-gray-300">
                    <div className="w-2 h-2 rounded-full bg-green-500"></div>
                    <span className="font-medium">Account Status</span>
                  </div>
                  <span className="px-3 py-1 bg-green-500/20 border border-green-500/30 text-green-400 rounded-xl text-sm font-semibold">
                    Active
                  </span>
                </div>
              </div>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-2 md:grid-cols-3 gap-4 mt-6">
              <div className="p-4 backdrop-blur-xl bg-gradient-to-br from-gray-800/40 to-gray-900/40 border border-gray-700/30 rounded-2xl text-center">
                <div className="text-2xl font-bold text-white">0</div>
                <div className="text-sm text-gray-400">Messages Sent</div>
              </div>
              <div className="p-4 backdrop-blur-xl bg-gradient-to-br from-gray-800/40 to-gray-900/40 border border-gray-700/30 rounded-2xl text-center">
                <div className="text-2xl font-bold text-white">0</div>
                <div className="text-sm text-gray-400">Images Shared</div>
              </div>
              <div className="col-span-2 md:col-span-1 p-4 backdrop-blur-xl bg-gradient-to-br from-gray-800/40 to-gray-900/40 border border-gray-700/30 rounded-2xl text-center">
                <div className="text-2xl font-bold text-white">0</div>
                <div className="text-sm text-gray-400">Conversations</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProfilePage;
