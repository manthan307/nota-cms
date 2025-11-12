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

  // handle case where schema is undefined or malformed
  if (!schema || !Array.isArray(schema.Definition)) {
    return (
      <Button variant="outline" disabled>
        <Plus className="mr-1" /> Create Content
      </Button>
    );
  }

  const handleChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Required fields validation
    for (const field of schema.Definition) {
      if (field.isRequired && !formData[field.name]?.trim()) {
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

      if (res.status === 200 || res.status === 201) {
        // support either shape: { data: { ... } } or raw
        const created = res.data?.data || res.data;
        addContent(created);
        refreshContents(schema.Name);
        setFormData({});
        setPublished(false);
        setOpen(false);
      } else {
        setError(res.data?.error || "Failed to create content.");
      }
    } catch (err) {
      console.error("Create content error:", err);
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
            {schema.Definition.map((field: any, idx: number) => (
              <div key={idx}>
                <label className="block text-sm font-medium mb-1">
                  {field.name}{" "}
                  {field.isRequired && <span className="text-red-500">*</span>}
                </label>
                <Input
                  type={
                    field.type === "number"
                      ? "number"
                      : field.type === "date"
                      ? "date"
                      : "text"
                  }
                  value={formData[field.name] || ""}
                  onChange={(e) => handleChange(field.name, e.target.value)}
                  placeholder={`Enter ${field.name}`}
                />
              </div>
            ))}

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
