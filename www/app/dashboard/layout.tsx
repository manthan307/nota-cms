"use client";
import { AppSidebar } from "@/components/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { Spinner } from "@/components/ui/spinner";
import { useAuth } from "@/context/auth";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { AudioWaveform, FileCode2, Home, Inbox, User } from "lucide-react";

const data = {
  name: "Nota CMS",
  logo: AudioWaveform,
  navMain: [
    // {
    //   title: "Home",
    //   url: "/dashboard",
    //   icon: Home,
    // },
    {
      title: "Schema",
      url: "/dashboard/schema",
      icon: FileCode2,
    },
    {
      title: "Content",
      url: "/dashboard/content",
      icon: Inbox,
    },
    {
      title: "Users",
      url: "/dashboard/users",
      icon: User,
    },
  ],
};

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const router = useRouter();
  const { user, loading } = useAuth();

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/");
    }
  }, [loading, user, router]);

  if (loading) {
    return (
      <div className="h-screen w-screen flex justify-center items-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  // Don't render the layout until the redirect happens
  if (!user) return null;

  return (
    <SidebarProvider>
      <AppSidebar user={user} data={data} />
      <SidebarInset>{children}</SidebarInset>
    </SidebarProvider>
  );
}
