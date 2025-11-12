import { AuthProvider } from "./auth";
import { ContentContextProvider } from "./contents";
import { SchemaContextProvider } from "./schema";

export const ContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  return (
    <AuthProvider>
      <SchemaContextProvider>
        <ContentContextProvider>{children}</ContentContextProvider>
      </SchemaContextProvider>
    </AuthProvider>
  );
};
