"use client";

import { useState, useContext } from "react";
import { fetch } from "@/lib/instance";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Spinner } from "@/components/ui/spinner";
import { Plus } from "lucide-react";
import { ContentContext } from "@/context/contents";

export function ContentCreateDialog({ schema }: { schema: any }) {
  const { addContent, refreshContents } = useContext(ContentContext);

  const [formData, setFormData] = useState<Record<string, any>>({});
  const [published, setPublished] = useState(false);
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);

  if (!schema || !Array.isArray(schema.Definition)) {
    return (
      <Button variant="outline" disabled>
        <Plus className="mr-1" /> Create Content
      </Button>
    );
  }

  const handleChange = (field: string, value: any) =>
    setFormData((prev) => ({ ...prev, [field]: value }));

  const uploadMedia = async (file: File) => {
    const form = new FormData();
    form.append("file", file);

    const res = await fetch.post("/api/v1/media/upload", form, {
      headers: { "Content-Type": "multipart/form-data" },
    });

    return res.data.url;
  };

  const handleFile = async (field: string, e: any) => {
    const file = e.target.files?.[0];
    if (!file) return;

    try {
      setLoading(true);
      const url = await uploadMedia(file);
      handleChange(field, url);
    } catch {
      setError("Failed to upload media.");
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: any) => {
    e.preventDefault();

    for (const field of schema.Definition) {
      if (field.isRequired && !formData[field.name]) {
        setError(`Field "${field.name}" is required.`);
        return;
      }
    }

    setLoading(true);
    setError(undefined);

    try {
      const body = {
        schema_id: schema.ID || schema.id,
        data: formData,
        published,
      };

      const res = await fetch.post("/api/v1/content/create", body);
      const created = res.data?.data || res.data;

      addContent(created);
      refreshContents(schema.Name);

      setFormData({});
      setPublished(false);
      setOpen(false);
    } catch {
      setError("Something went wrong while creating content.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">
          <Plus className="mr-1" /> Create Content
        </Button>
      </DialogTrigger>

      <DialogContent className="max-w-[700px] w-full">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>Create New {schema.Name} Content</DialogTitle>
            <DialogDescription>
              Fill in the fields according to the schema definition.
            </DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 mt-4">
            {schema.Definition.map((field: any, idx: number) => {
              const value = formData[field.name] || "";

              return (
                <div key={idx}>
                  <label className="block text-sm font-medium mb-1">
                    {field.name}
                    {field.isRequired && (
                      <span className="text-red-500 ml-1">*</span>
                    )}
                  </label>

                  {field.type === "text" && (
                    <Input
                      value={value}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                    />
                  )}

                  {field.type === "number" && (
                    <Input
                      type="number"
                      value={value}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                    />
                  )}

                  {field.type === "date" && (
                    <Input
                      type="date"
                      value={value}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                    />
                  )}

                  {field.type === "textarea" && (
                    <textarea
                      className="border rounded p-2 w-full"
                      rows={4}
                      value={value}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                    />
                  )}

                  {field.type === "boolean" && (
                    <Switch
                      checked={value === "true" || value === true}
                      onCheckedChange={(v) => handleChange(field.name, `${v}`)}
                    />
                  )}

                  {(field.type === "image" ||
                    field.type === "video" ||
                    field.type === "file") && (
                    <div>
                      <Input
                        type="file"
                        accept={
                          field.type === "image"
                            ? "image/*"
                            : field.type === "video"
                            ? "video/*"
                            : undefined
                        }
                        onChange={(e) => handleFile(field.name, e)}
                      />

                      {value && typeof value === "string" && (
                        <>
                          {field.type === "image" && (
                            <img
                              src={value}
                              className="mt-2 h-24 w-auto rounded border"
                            />
                          )}

                          {field.type === "video" && (
                            <video
                              src={value}
                              controls
                              className="mt-2 h-32 w-auto rounded border"
                            />
                          )}

                          {field.type === "file" && (
                            <a
                              href={value}
                              target="_blank"
                              className="text-blue-600 underline mt-2 block"
                            >
                              View file
                            </a>
                          )}
                        </>
                      )}
                    </div>
                  )}
                </div>
              );
            })}

            <div className="flex items-center gap-2 mt-3">
              <Switch
                checked={published}
                onCheckedChange={setPublished}
                id="published"
              />
              <label htmlFor="published" className="text-sm">
                Published
              </label>
            </div>
          </div>

          {error && <p className="text-red-500 mt-2">{error}</p>}

          <DialogFooter className="mt-6 flex justify-end gap-2">
            <DialogClose asChild>
              <Button variant="outline" disabled={loading}>
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={loading}>
              {loading && <Spinner className="size-4 mr-2" />}
              Create
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
