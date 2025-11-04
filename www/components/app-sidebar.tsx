"use client";

import * as React from "react";

import { NavMain } from "@/components/nav-main";
import { Header } from "@/components/navheader";
import {
  Sidebar,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from "@/components/ui/sidebar";
import { NavUser } from "./nav-user";
import { LucideIcon } from "lucide-react";

export function AppSidebar({
  user,
  data,
  ...props
}: {
  user: { email: string; role: string };
  data: {
    name: string;
    logo: React.ElementType;
    navMain: {
      title: string;
      url: string;
      icon: LucideIcon;
    }[];
  };
} & React.ComponentProps<typeof Sidebar>) {
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
        <NavUser user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
}
