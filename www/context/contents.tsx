"use client";

import { fetch } from "@/lib/instance";
import { createContext, useEffect, useState } from "react";

type ContentResponse = {
  [x: string]: any;
  error?: string;
  data?: any[];
  count?: number;
  message?: string;
};

export const ContentContext = createContext<{
  loading: boolean;
  contents: any[];
  schemaName: string | null;
  setSchemaName: React.Dispatch<React.SetStateAction<string | null>>;
  refreshContents: (schemaName?: string) => Promise<void>;
  addContent: (content: any) => void;
  updateContent: (updatedContent: any) => void;
  deleteContent: (id: string) => void;
  setContents: React.Dispatch<React.SetStateAction<any[]>>;
}>({
  loading: true,
  contents: [],
  schemaName: null,
  setSchemaName: () => {},
  refreshContents: async () => {},
  addContent: () => {},
  updateContent: () => {},
  deleteContent: () => {},
  setContents: () => {},
});

export const ContentContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [contents, setContents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [schemaName, setSchemaName] = useState<string | null>(null);

  const refreshContents = async (schema?: string) => {
    const name = schema || schemaName;
    if (!name) {
      setContents([]);
      return;
    }

    setLoading(true);
    setContents([]); // Clear old data

    try {
      const res = await fetch.get<ContentResponse>(
        `/api/v1/content/get_all/${name}`
      );

      // Normalize key casing from backend → frontend
      const normalized = (res?.data ?? []).map((item: any) => ({
        ID: item.ID ?? item.id,
        SchemaID: item.SchemaID ?? item.schemaID,
        Data: item.Data ?? item.data,
        CreatedAt: item.CreatedAt ?? item.createdAt,
      }));

      setContents(normalized);
    } catch (err) {
      console.error("Failed to fetch contents:", err);
      setContents([]);
    } finally {
      setLoading(false);
    }
  };

  const addContent = (content: any) => {
    setContents((prev) => [...prev, content]);
  };

  const updateContent = (updatedContent: any) => {
    setContents((prev) =>
      prev.map((c) =>
        (c.ID ?? c.id) === (updatedContent.ID ?? updatedContent.id)
          ? updatedContent
          : c
      )
    );
  };

  const deleteContent = (id: string) => {
    setContents((prev) => prev.filter((c) => (c.ID ?? c.id) !== id));
  };

  useEffect(() => {
    if (schemaName) refreshContents(schemaName);
  }, [schemaName]);

  return (
    <ContentContext.Provider
      value={{
        contents,
        loading,
        schemaName,
        setSchemaName,
        refreshContents,
        addContent,
        updateContent,
        deleteContent,
        setContents,
      }}
    >
      {children}
    </ContentContext.Provider>
  );
};
