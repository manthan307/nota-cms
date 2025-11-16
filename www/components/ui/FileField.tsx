"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { fetch } from "@/lib/instance";

export function FileField({
  value,
  onChange,
  editing,
}: {
  value: string;
  onChange: (url: string) => void;
  editing: boolean;
}) {
  const [uploading, setUploading] = useState(false);

  async function replaceFile(oldUrl: string, newFile: File) {
    try {
      setUploading(true);

      if (oldUrl) {
        await fetch.delete("/api/v1/media/delete", {
          data: { file_url: oldUrl },
        });
      }

      const form = new FormData();
      form.append("file", newFile);

      const res = await fetch.post("/api/v1/media/upload", form);
      const newUrl = res.data?.url;

      if (newUrl) onChange(newUrl);
    } finally {
      setUploading(false);
    }
  }

  function chooseFile() {
    const picker = document.createElement("input");
    picker.type = "file";
    picker.onchange = (e: any) => {
      const file = e.target.files?.[0];
      if (file) replaceFile(value, file);
    };
    picker.click();
  }

  function renderPreview() {
    if (!value) return null;

    const isImage = /\.(jpg|jpeg|png|gif|webp)$/i.test(value);
    const isVideo = /\.(mp4|mov|webm|mkv)$/i.test(value);

    if (isImage) {
      return (
        <img
          src={value}
          className="rounded max-h-48 object-cover border"
          alt="uploaded preview"
        />
      );
    }

    if (isVideo) {
      return <video src={value} controls className="rounded max-h-48 border" />;
    }

    return (
      <div className="border rounded p-3 flex items-center justify-between">
        <p className="truncate max-w-[80%]">{value}</p>
        <a
          href={value}
          target="_blank"
          className="text-sm underline text-blue-600"
        >
          open
        </a>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-3">
      {renderPreview()}

      {editing && (
        <Button
          type="button"
          variant="outline"
          disabled={uploading}
          className="w-fit"
          onClick={chooseFile}
        >
          {uploading ? "Uploading..." : value ? "Replace File" : "Upload File"}
        </Button>
      )}
    </div>
  );
}
