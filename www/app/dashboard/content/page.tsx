"use client";

import { useContext, useEffect, useState } from "react";
import { SchemaContext } from "@/context/schema";
import { Spinner } from "@/components/ui/spinner";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Pencil, Save, Trash2 } from "lucide-react";
import { ContentContext } from "@/context/contents";
import { SidebarTrigger } from "@/components/ui/sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbPage,
} from "@/components/ui/breadcrumb";
import { ContentCreateDialog } from "@/components/contents/create";
import { fetch } from "@/lib/instance";
import { Switch } from "@/components/ui/switch";

export default function PageContentDashboard() {
  const { schemas, loading: schemaLoading } = useContext(SchemaContext);
  const {
    contents,
    loading: contentLoading,
    refreshContents,
    updateContent,
    deleteContent,
  } = useContext(ContentContext);

  const [selectedSchema, setSelectedSchema] = useState<any | null>(null);
  const [selectedContent, setSelectedContent] = useState<any | null>(null);
  const [editingContent, setEditingContent] = useState<any | null>(null);

  // Fetch contents when schema changes
  useEffect(() => {
    if (selectedSchema?.Name) {
      refreshContents(selectedSchema.Name);
      setSelectedContent(null);
    }
  }, [selectedSchema]);

  const handleSave = async () => {
    if (!editingContent) return;
    try {
      await fetch.post(
        `/api/v1/content/update/${editingContent.ID}`,
        editingContent
      );
      updateContent(editingContent);
      setSelectedContent(editingContent);
      setEditingContent(null);
    } catch (err) {
      console.error("Failed to update content:", err);
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await fetch.delete(`/api/v1/content/delete/${id}`);
      deleteContent(id);
      if (selectedContent?.ID === id) setSelectedContent(null);
    } catch (err) {
      console.error("Failed to delete content:", err);
    }
  };

  const handlePublishToggle = async (checked: boolean) => {
    if (!selectedContent) return;
    const updated = { ...selectedContent, Published: checked };
    try {
      await fetch.post(`/api/v1/content/update/${selectedContent.ID}`, updated);
      updateContent(updated);
      setSelectedContent(updated);
    } catch (err) {
      console.error("Failed to update publish state:", err);
    }
  };

  if (schemaLoading)
    return (
      <div className="h-screen w-screen flex justify-center items-center">
        <Spinner className="size-8" />
      </div>
    );

  if (!schemas?.length)
    return (
      <div className="flex h-full w-full justify-center items-center">
        <p>No Schemas Found.</p>
      </div>
    );

  return (
    <>
      <header className="flex h-14 shrink-0 items-center gap-2">
        <div className="flex flex-1 items-center gap-2 px-3">
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
        {selectedSchema && <ContentCreateDialog schema={selectedSchema} />}
      </header>

      <section className="flex h-full w-full">
        {/* Schemas Sidebar */}
        <div className="w-1/4 border-r p-4 overflow-y-auto">
          <h2 className="text-lg font-semibold mb-4">Schemas</h2>
          <div className="flex flex-col gap-2">
            {schemas.map((schema) => (
              <Button
                key={schema.ID}
                variant={
                  selectedSchema?.ID === schema.ID ? "default" : "outline"
                }
                className="justify-start"
                onClick={() => setSelectedSchema(schema)}
              >
                {schema.Name}
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
          ) : contents.length === 0 ? (
            <p>No contents found.</p>
          ) : (
            <div className="flex flex-col gap-2">
              {contents.map((c) => (
                <Button
                  key={c.ID ?? c.id}
                  variant={
                    selectedContent?.ID === (c.ID ?? c.id)
                      ? "default"
                      : "outline"
                  }
                  className="justify-start"
                  onClick={() => setSelectedContent(c)}
                >
                  {c.data?.Title ||
                    c.Data?.Title ||
                    `Content ${c.ID.slice(0, 6)}`}
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
                  {selectedContent.Data?.Title ||
                    selectedContent.data?.Title ||
                    "Untitled"}
                </CardTitle>
                <div className="flex items-center gap-2">
                  <label className="text-sm text-muted-foreground">
                    Published
                  </label>
                  <Switch
                    checked={selectedContent.Published ?? false}
                    onCheckedChange={handlePublishToggle}
                  />
                </div>
              </CardHeader>

              <CardContent className="grid gap-3">
                {selectedSchema.Definition.map((field: any, idx: number) => (
                  <div key={idx}>
                    <label className="text-sm font-medium">{field.name}</label>
                    <Input
                      className="mt-1"
                      type="text"
                      value={
                        editingContent?.ID === selectedContent.ID
                          ? editingContent.Data?.[field.name] ?? ""
                          : selectedContent.Data?.[field.name] ?? ""
                      }
                      disabled={editingContent?.ID !== selectedContent.ID}
                      onChange={(e) => {
                        if (editingContent?.ID === selectedContent.ID) {
                          setEditingContent((prev: any) => ({
                            ...prev,
                            Data: {
                              ...prev.Data,
                              [field.name]: e.target.value,
                            },
                          }));
                        }
                      }}
                    />
                  </div>
                ))}

                <div className="flex gap-2 justify-end mt-4">
                  {editingContent?.ID === selectedContent.ID ? (
                    <Button size="sm" onClick={handleSave}>
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
                  >
                    <Trash2 className="size-4" />
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </section>
    </>
  );
}
