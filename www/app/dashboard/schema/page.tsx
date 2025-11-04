"use client";
import { Separator } from "@/components/ui/separator";
import { fetch } from "@/lib/instance";
import { useEffect } from "react";

export default function PageSchemaDashboard() {
  useEffect(() => {
    (async () => {
      try {
        const res = await fetch.post("/api/v1/schemas/list");

        console.log(res);
      } catch (err) {
        console.error(err);
      }
    })();
  }, []);
  return (
    <section className="flex h-full w-full">
      <div className="h-full w-full"></div>
      <Separator className="my-4" orientation="vertical" />
      <div className="h-full w-full"></div>
    </section>
  );
}
