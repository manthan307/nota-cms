"use client";

import { useContext, useEffect, useMemo, useState } from "react";
import { SchemaContext } from "@/context/schema";
import { ContentContext } from "@/context/contents";
import { fetch } from "@/lib/instance";
import { Spinner } from "@/components/ui/spinner";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Pencil, Save, Trash2 } from "lucide-react";
import { SidebarTrigger } from "@/components/ui/sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { ContentCreateDialog } from "@/components/contents/create";
import { Switch } from "@/components/ui/switch";
import { toast } from "sonner";
import { FileField } from "@/components/ui/FileField";

/* --- types used locally --- */
type NormalizedContent = {
  ID: string;
  SchemaID?: string;
  Data: Record<string, any>;
  Published: boolean;
  CreatedAt?: string;
  UpdatedAt?: string;
  _raw?: any;
};

export default function PageContentDashboard() {
  const { schemas, loading: schemaLoading } = useContext(SchemaContext);
  const {
    contents,
    loading: contentLoading,
    refreshContents,
    addContent,
    updateContent,
    deleteContent,
    setContents,
  } = useContext(ContentContext);

  const [selectedSchema, setSelectedSchema] = useState<any | null>(null);
  const [selectedContent, setSelectedContent] =
    useState<NormalizedContent | null>(null);
  const [editingContent, setEditingContent] =
    useState<NormalizedContent | null>(null);
  const [localLoading, setLocalLoading] = useState(false); // for actions like save/delete/publish

  // Normalize helper
  const normalizeContentItem = (raw: any): NormalizedContent => {
    const id =
      raw.ID ??
      raw.id ??
      raw.Id ??
      raw.ID?.toString() ??
      raw.id?.toString() ??
      "";
    const schemaID =
      raw.SchemaID ?? raw.schemaID ?? raw.schema_id ?? raw.schemaId;
    // data may come as .Data (object) or .data (object) or nested in response
    const data = raw.Data ?? raw.data ?? raw.payload ?? {};
    // published can be boolean or an object {Bool:bool, Valid:bool} or {published:true}
    let published = false;
    if (typeof raw.Published === "boolean") {
      published = raw.Published;
    } else if (raw.Published && typeof raw.Published === "object") {
      // pgtype style: {Bool: true, Valid: true}
      published =
        raw.Published.Bool ??
        raw.Published.bool ??
        raw.Published.value ??
        false;
    } else if (typeof raw.published === "boolean") {
      published = raw.published;
    } else if (raw.published && typeof raw.published === "object") {
      published = raw.published.Bool ?? raw.published.bool ?? false;
    }

    return {
      ID: id,
      SchemaID: schemaID,
      Data: typeof data === "object" && data !== null ? data : {},
      Published: Boolean(published),
      CreatedAt: raw.createdAt ?? raw.CreatedAt ?? raw.created_at,
      UpdatedAt: raw.updatedAt ?? raw.UpdatedAt ?? raw.updated_at,
      _raw: raw,
    };
  };

  // When schema is selected, refresh contents for it
  useEffect(() => {
    if (selectedSchema?.Name) {
      // default to fetching all (published + drafts) in the editor dashboard
      refreshContents(selectedSchema.Name);
      setSelectedContent(null);
      setEditingContent(null);
    } else {
      // clear contents if no schema selected
      setContents([]);
      setSelectedContent(null);
      setEditingContent(null);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedSchema]);

  // Normalize contents coming from context for local usage (memoized)
  const normalizedContents = useMemo(() => {
    if (!Array.isArray(contents)) return [];
    return contents.map((c) => normalizeContentItem(c));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [contents]);

  // select content (normalize on select)
  const handleSelectContent = (raw: any) => {
    const norm = normalizeContentItem(raw);
    setSelectedContent(norm);
    setEditingContent(null);
  };

  // Save edited content
  const handleSave = async () => {
    if (!editingContent) return;
    setLocalLoading(true);
    try {
      const body = {
        content_id: editingContent.ID,
        data: editingContent.Data,
        published: editingContent.Published,
      };

      // v1: router registers POST /api/v1/content/update
      const res = await fetch.post("/api/v1/content/update", body);

      // normalise server response if it returns updated content
      const returned = res?.data ?? res;
      const updatedNormalized = normalizeContentItem(returned);
      // update context
      updateContent(updatedNormalized);
      setSelectedContent(updatedNormalized);
      setEditingContent(null);
    } catch (err) {
      console.error("Failed to update content:", err);
      toast.error("Failed to update content.");
    } finally {
      setLocalLoading(false);
    }
  };

  // Delete content
  const handleDelete = async (id: string) => {
    if (!id) return;
    setLocalLoading(true);
    try {
      await fetch.delete(`/api/v1/content/delete/${id}`);
      deleteContent(id);
      if (selectedContent?.ID === id) {
        setSelectedContent(null);
        setEditingContent(null);
      }
    } catch (err) {
      console.error("Failed to delete content:", err);
    } finally {
      setLocalLoading(false);
    }
  };

  // Toggle publish: update server & state
  const handlePublishToggle = async (checked: boolean) => {
    if (!selectedContent) return;
    setLocalLoading(true);
    try {
      // backend expects { content_id, data, published } in POST /update
      const body = {
        content_id: selectedContent.ID,
        data: selectedContent.Data,
        published: checked,
      };
      const res = await fetch.post("/api/v1/content/update", body);
      const updated = normalizeContentItem(res.data ?? res);
      updateContent(updated);
      setSelectedContent(updated);
      // if editing, also update editingContent
      if (editingContent?.ID === updated.ID) setEditingContent(updated);
    } catch (err) {
      console.error("Failed to update publish state:", err);
    } finally {
      setLocalLoading(false);
    }
  };

  // UI helpers
  if (schemaLoading) {
    return (
      <div className="h-screen w-screen flex justify-center items-center">
        <Spinner className="size-8" />
      </div>
    );
  }

  if (!schemas || schemas.length === 0) {
    return (
      <div className="flex h-full w-full justify-center items-center">
        <p>No Schemas Found.</p>
      </div>
    );
  }

  return (
    <>
      <header className="flex h-14 shrink-0 items-center gap-2 px-3">
        <div className="flex flex-1 items-center gap-2">
          <SidebarTrigger />
          <Separator orientation="vertical" className="mr-2 h-4" />
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbPage>Contents</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>

        {/* create content dialog only when a schema is selected */}
        {selectedSchema && <ContentCreateDialog schema={selectedSchema} />}
      </header>

      <section className="flex h-full w-full">
        {/* Schemas Sidebar */}
        <div className="w-1/4 border-r p-4 overflow-y-auto">
          <h2 className="text-lg font-semibold mb-4">Schemas</h2>
          <div className="flex flex-col gap-2">
            {schemas.map((schema: any) => (
              <Button
                key={schema.ID ?? schema.id}
                variant={
                  selectedSchema?.ID === (schema.ID ?? schema.id)
                    ? "default"
                    : "outline"
                }
                className="justify-start"
                onClick={() => setSelectedSchema(schema)}
              >
                {schema.Name ?? schema.name}
              </Button>
            ))}
          </div>
        </div>

        {/* Contents List */}
        <div className="w-1/4 border-r p-4 overflow-y-auto">
          {!selectedSchema ? (
            <p className="text-muted-foreground">Select a schema.</p>
          ) : contentLoading ? (
            <div className="flex justify-center mt-10">
              <Spinner className="size-6" />
            </div>
          ) : normalizedContents.length === 0 ? (
            <p>No contents found.</p>
          ) : (
            <div className="flex flex-col gap-2">
              {normalizedContents.map((c) => (
                <Button
                  key={c.ID}
                  variant={selectedContent?.ID === c.ID ? "default" : "outline"}
                  className="justify-start"
                  onClick={() => handleSelectContent(c._raw ?? c)}
                >
                  {(() => {
                    const firstField = selectedSchema?.Definition?.[0]?.name;
                    return (
                      c.Data?.[firstField] ?? `Content ${c.ID.slice(0, 6)}`
                    );
                  })()}
                </Button>
              ))}
            </div>
          )}
        </div>

        {/* Content Details */}
        <div className="flex-1 p-4 overflow-y-auto">
          {!selectedContent ? (
            <div className="h-full flex items-center justify-center">
              <p className="text-muted-foreground">Select a content to view.</p>
            </div>
          ) : (
            <Card className="shadow-sm">
              <CardHeader className="flex justify-between items-center">
                <CardTitle>
                  {(() => {
                    const firstField = selectedSchema?.Definition?.[0]?.name;
                    return (
                      selectedContent.Data?.[firstField] ??
                      `Content ${selectedContent.ID.slice(0, 6)}`
                    );
                  })()}
                </CardTitle>

                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-2">
                    <label className="text-sm text-muted-foreground">
                      Published
                    </label>
                    <Switch
                      checked={selectedContent.Published}
                      onCheckedChange={async (v) => {
                        // update UI optimistically
                        setSelectedContent((prev) =>
                          prev ? { ...prev, Published: v } : prev
                        );
                        setEditingContent((prev) =>
                          prev ? { ...prev, Published: v } : prev
                        );
                        await handlePublishToggle(Boolean(v));
                      }}
                      disabled={localLoading}
                    />
                  </div>

                  <div className="flex gap-2">
                    {editingContent?.ID === selectedContent.ID ? (
                      <Button
                        size="sm"
                        onClick={handleSave}
                        disabled={localLoading}
                      >
                        <Save className="size-4 mr-1" /> Save
                      </Button>
                    ) : (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() =>
                          setEditingContent({
                            ...selectedContent,
                            Data: { ...selectedContent.Data },
                          })
                        }
                      >
                        <Pencil className="size-4 mr-1" /> Edit
                      </Button>
                    )}
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => handleDelete(selectedContent.ID)}
                      disabled={localLoading}
                    >
                      <Trash2 className="size-4" />
                    </Button>
                  </div>
                </div>
              </CardHeader>

              <CardContent className="grid gap-3">
                {/* Render fields from selected schema definition (if present) */}
                {selectedSchema?.Definition?.map((field: any, idx: number) => (
                  <div key={idx}>
                    <label className="text-sm font-medium">{field.name}</label>
                    {["image", "video", "file"].includes(field.type) ? (
                      <FileField
                        value={
                          editingContent?.ID === selectedContent.ID
                            ? editingContent.Data?.[field.name] ?? ""
                            : selectedContent.Data?.[field.name] ?? ""
                        }
                        editing={editingContent?.ID === selectedContent.ID}
                        onChange={(v) => {
                          if (editingContent?.ID === selectedContent.ID) {
                            setEditingContent((prev) =>
                              prev
                                ? {
                                    ...prev,
                                    Data: { ...prev.Data, [field.name]: v },
                                  }
                                : prev
                            );
                          }
                        }}
                      />
                    ) : (
                      <Input
                        className="mt-1"
                        type={
                          field.type === "number"
                            ? "number"
                            : field.type === "date"
                            ? "date"
                            : "text"
                        }
                        value={
                          editingContent?.ID === selectedContent.ID
                            ? editingContent.Data?.[field.name] ?? ""
                            : selectedContent.Data?.[field.name] ?? ""
                        }
                        disabled={
                          editingContent?.ID !== selectedContent.ID ||
                          localLoading
                        }
                        onChange={(e) => {
                          if (editingContent?.ID === selectedContent.ID) {
                            setEditingContent((prev) => {
                              if (!prev) return prev;
                              return {
                                ...prev,
                                Data: {
                                  ...prev.Data,
                                  [field.name]: e.target.value,
                                },
                              };
                            });
                          }
                        }}
                      />
                    )}
                  </div>
                ))}

                {/* If schema has no definition (defensive), show JSON editor for data */}
                {(!selectedSchema?.Definition ||
                  selectedSchema.Definition.length === 0) && (
                  <div>
                    <label className="text-sm font-medium">Data (JSON)</label>
                    <textarea
                      className="mt-1 w-full min-h-[200px] rounded border p-2"
                      value={JSON.stringify(
                        editingContent?.ID === selectedContent.ID
                          ? editingContent.Data
                          : selectedContent.Data,
                        null,
                        2
                      )}
                      onChange={(e) => {
                        try {
                          const parsed = JSON.parse(e.target.value);
                          setEditingContent((prev) =>
                            prev ? { ...prev, Data: parsed } : prev
                          );
                        } catch {
                          // ignore parse errors while editing
                        }
                      }}
                      disabled={
                        editingContent?.ID !== selectedContent.ID ||
                        localLoading
                      }
                    />
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>
      </section>
    </>
  );
}
