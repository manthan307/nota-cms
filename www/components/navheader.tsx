"use client";

import * as React from "react";
import { SidebarMenu, SidebarMenuItem } from "@/components/ui/sidebar";

export function Header({
  name,
  Logo,
}: {
  name: string;
  Logo: React.ElementType;
}) {
  return (
    <SidebarMenu>
      <SidebarMenuItem className="flex items-center gap-2 px-2 py-3">
        <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-5 items-center justify-center rounded-md">
          <Logo className="size-3" />
        </div>
        <span className="truncate font-medium">{name}</span>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
