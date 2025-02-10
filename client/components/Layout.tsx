import { ReactNode } from "react";
import Navbar from "@/components/Navbar";

const Layout = ({ children }: { children: ReactNode }) => {
  return (
    <div className="min-h-screen flex flex-col bg-gray-900 text-gray-100">
      <Navbar />
      <main className="flex-grow container mx-auto px-4 py-8">{children}</main>
      <footer className="py-4 text-center text-sm text-gray-500">
        &copy; {new Date().getFullYear()} Secret Phrase
      </footer>
    </div>
  );
};

export default Layout;
