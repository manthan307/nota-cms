"use client";

import { fetch } from "@/lib/instance";
import { createContext, useEffect, useState } from "react";

type SchemaResponse = {
  error?: string;
  data?: any[];
  count?: number;
  message?: string;
};

export const SchemaContext = createContext<{
  loading: boolean;
  schemas: any[];
  refreshSchemas: () => Promise<void>;
  addSchema: (schema: any) => void;
  updateSchema: (updatedSchema: any) => void;
  deleteSchema: (id: string) => void;
  setSchemas: React.Dispatch<React.SetStateAction<any[]>>; // optional utility
}>({
  loading: true,
  schemas: [],
  refreshSchemas: async () => {},
  addSchema: () => {},
  updateSchema: () => {},
  deleteSchema: () => {},
  setSchemas: () => {},
});

export const SchemaContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [schemas, setSchemas] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  const refreshSchemas = async () => {
    try {
      setLoading(true);
      const res = await fetch.get<SchemaResponse>("/api/v1/schemas/list");
      setSchemas(res.data?.data || []);
    } catch (err) {
      console.error("Failed to fetch schemas:", err);
    } finally {
      setLoading(false);
    }
  };

  const addSchema = (schema: any) => {
    setSchemas((prev) => [...prev, schema]);
  };

  const updateSchema = (updatedSchema: any) => {
    setSchemas((prev) =>
      prev.map((schema) =>
        schema.ID === updatedSchema.ID ? updatedSchema : schema
      )
    );
  };

  const deleteSchema = (id: string) => {
    setSchemas((prev) => prev.filter((schema) => schema.ID !== id));
  };

  useEffect(() => {
    refreshSchemas();
  }, []);

  return (
    <SchemaContext.Provider
      value={{
        schemas,
        loading,
        refreshSchemas,
        addSchema,
        updateSchema,
        deleteSchema,
        setSchemas,
      }}
    >
      {children}
    </SchemaContext.Provider>
  );
};
