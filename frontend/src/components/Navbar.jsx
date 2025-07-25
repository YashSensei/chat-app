import { Link } from "react-router-dom";
import { useAuthStore } from "../store/useAuthStore";
import { LogOut, MessageSquare, User } from "lucide-react";

const Navbar = () => {
  const { logout, authUser } = useAuthStore();

  return (
    <header className="bg-surface border-b border-border fixed w-full top-0 z-40 shadow-lg">
      <div className="container mx-auto px-4 h-14 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 hover:opacity-80 transition-all">
          <div className="size-8 rounded-lg bg-primary flex items-center justify-center">
            <MessageSquare className="w-5 h-5 text-white" />
          </div>
          <h1 className="text-xl font-display font-bold tracking-tight text-white">chatt-app</h1>
        </Link>
        <div className="flex items-center gap-2">
          <Link to="/" className="px-3 py-1 rounded-md bg-background text-muted hover:bg-primary hover:text-white transition-all flex items-center gap-1">
            <MessageSquare className="w-4 h-4" />
            <span className="hidden sm:inline">Chat</span>
          </Link>
          {authUser && (
            <>
              <Link to={"/profile"} className="px-3 py-1 rounded-md bg-background text-muted hover:bg-accent hover:text-white transition-all flex items-center gap-1">
                <User className="size-5" />
                <span className="hidden sm:inline">Profile</span>
              </Link>
              <button className="px-3 py-1 rounded-md bg-background text-muted hover:bg-secondary hover:text-white transition-all flex items-center gap-1" onClick={logout}>
                <LogOut className="size-5" />
                <span className="hidden sm:inline">Logout</span>
              </button>
            </>
          )}
        </div>
      </div>
    </header>
  );
};
export default Navbar;
