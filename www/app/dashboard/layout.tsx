"use client";
import { AppSidebar } from "@/components/app-sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import { Spinner } from "@/components/ui/spinner";
import { useAuth } from "@/context/auth";
import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";
import {
  AudioWaveform,
  FileCode2,
  Home,
  Inbox,
  Plus,
  User,
} from "lucide-react";
import { Button } from "@/components/ui/button";

const data = {
  name: "Nota CMS",
  logo: AudioWaveform,
  navMain: [
    {
      title: "Home",
      url: "/dashboard",
      icon: Home,
    },
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
  const pathName = usePathname();
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
      <SidebarInset>
        <header className="flex h-14 shrink-0 items-center gap-2">
          <div className="flex flex-1 items-center gap-2 px-3">
            <SidebarTrigger />
            <Separator
              orientation="vertical"
              className="mr-2 data-[orientation=vertical]:h-4"
            />
            <Breadcrumb>
              <BreadcrumbList>
                <BreadcrumbItem>
                  <BreadcrumbPage className="line-clamp-1">
                    {data.navMain.find((v) => v.url === pathName)?.title || ""}
                  </BreadcrumbPage>
                </BreadcrumbItem>
              </BreadcrumbList>
            </Breadcrumb>
          </div>
          <Button className="mx-8 self-center">
            <Plus /> Create
          </Button>
        </header>
        {children}
      </SidebarInset>
    </SidebarProvider>
  );
}
