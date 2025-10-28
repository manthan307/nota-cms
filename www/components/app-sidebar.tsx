"use client";

import * as React from "react";
import {
  AudioWaveform,
  FileCode2,
  Home,
  Image,
  Inbox,
  Search,
  Sparkles,
} from "lucide-react";

import { NavMain } from "@/components/nav-main";
import { Header } from "@/components/navheader";
import {
  Sidebar,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { NavUser } from "./nav-user";

const data = {
  name: "Nota CMS",
  logo: AudioWaveform,
  navMain: [
    {
      title: "Schema",
      url: "#",
      icon: FileCode2,
    },
    {
      title: "Content",
      url: "#",
      icon: Inbox,
    },
    {
      title: "Media",
      url: "#",
      icon: Image,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar
      variant="inset"
      collapsible="icon"
      className="border-r-0 flex flex-col h-full"
      {...props}
    >
      <div className="grow">
        <SidebarHeader>
          <Header name={data.name} Logo={data.logo} />
          <NavMain items={data.navMain} />
        </SidebarHeader>
      </div>
      <SidebarFooter>
        <NavUser
          user={{
            name: "Manthan Patel",
            email: "pmanthan549@gmail.com",
          }}
        />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
